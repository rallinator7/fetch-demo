package payer

import (
	"encoding/json"
	"net/http"

	internal "github.com/rallinator7/fetch-demo/internal"
	"github.com/rallinator7/fetch-demo/internal/payer/api"
)

type PayerController interface {
	AddPayer(payer internal.Payer) (internal.Payer, error)
}

type StatusCoder interface {
	StatusCode() int64
}

type Handler struct {
	controller PayerController
}

func NewHandler(controller PayerController) Handler {
	return Handler{
		controller: controller,
	}
}

// AddPayer parse the open api /payer/add route and attempts to add the payer.
func (router *Handler) AddPayer(w http.ResponseWriter, r *http.Request) {
	var payerName api.PayerName

	err := json.NewDecoder(r.Body).Decode(&payerName)
	if err != nil {
		ErrorResponse(w, internal.NewStatusCodeError(http.StatusBadRequest, err))
		return
	}

	payer := internal.NewPayer(payerName.Name)

	_, err = router.controller.AddPayer(payer)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		&api.Payer{
			Name: payer.Name,
			Id:   payer.Id,
		},
	)
}

// ErrorResponse attempts to read the status code from the error and returns the code and the message to the caller.
func ErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var code int64

	coder, ok := err.(StatusCoder)
	if !ok {
		code = http.StatusInternalServerError
	} else {
		code = coder.StatusCode()
	}

	w.WriteHeader(int(code))

	resp := api.Error{
		Code:    code,
		Message: err.Error(),
	}

	json.NewEncoder(w).Encode(resp)
}
