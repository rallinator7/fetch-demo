package user

import (
	"encoding/json"
	"net/http"

	internal "github.com/rallinator7/fetch-demo/internal"
	"github.com/rallinator7/fetch-demo/internal/user/api"
)

type UserController interface {
	AddUser(user internal.User) (internal.User, error)
}

type StatusCoder interface {
	StatusCode() int64
}

type Handler struct {
	controller UserController
}

func NewHandler(controller UserController) Handler {
	return Handler{
		controller: controller,
	}
}

// AddUser parse the open api /user/add route and attempts to add the user.
func (router *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	var userName api.UserName

	err := json.NewDecoder(r.Body).Decode(&userName)
	if err != nil {
		ErrorResponse(w, internal.NewStatusCodeError(http.StatusBadRequest, err))
		return
	}

	user := internal.NewUser(userName.FirstName, userName.LastName)

	_, err = router.controller.AddUser(user)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		&api.User{
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Id:        user.Id,
		},
	)
}

// ErrorResponse attempts to read the status code from the error and returns the code and the message to the caller.
func ErrorResponse(w http.ResponseWriter, err error) {
	var code int64

	coder, ok := err.(StatusCoder)
	if !ok {
		code = http.StatusInternalServerError
	} else {
		code = coder.StatusCode()
	}

	w.WriteHeader(int(code))
	w.Header().Set("Content-Type", "application/json")

	resp := api.Error{
		Code:    code,
		Message: err.Error(),
	}

	json.NewEncoder(w).Encode(resp)
}
