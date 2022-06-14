package points

import (
	"encoding/json"
	"net/http"

	"github.com/rallinator7/fetch-demo/internal"
	"github.com/rallinator7/fetch-demo/internal/points/api"
)

type PointsController interface {
	DescribeBalance(id string) (map[string]int, error)
	GiveUserPoints(id string, transaction PayerTransaction) error
	SpendUserPoints(id string, points int) ([]PayerPoints, error)
}

type StatusCoder interface {
	StatusCode() int
}

type Handler struct {
	controller PointsController
}

func NewHandler(controller PointsController) Handler {
	return Handler{
		controller: controller,
	}
}

func (handler *Handler) GivePoints(w http.ResponseWriter, r *http.Request, user string) {
	var give api.GivePoints

	err := json.NewDecoder(r.Body).Decode(&give)
	if err != nil {
		ErrorResponse(w, internal.NewStatusCodeError(http.StatusBadRequest, err))
		return
	}

	transaction, err := NewPayerTransaction(give.Payer, give.Points, give.Timestamp)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	err = handler.controller.GiveUserPoints(user, transaction)
	if err != nil {
		ErrorResponse(w, err)
		return
	}
}

func (handler *Handler) SpendPoints(w http.ResponseWriter, r *http.Request, user string) {
	var spend api.SpendPoints

	err := json.NewDecoder(r.Body).Decode(&spend)
	if err != nil {
		ErrorResponse(w, internal.NewStatusCodeError(http.StatusBadRequest, err))
		return
	}

	payerPoints, err := handler.controller.SpendUserPoints(user, spend.Points)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	slice := convertPayerPointsSlice(payerPoints)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		slice,
	)
}

func (handler *Handler) DescribeBalance(w http.ResponseWriter, r *http.Request, user string) {
	balance, err := handler.controller.DescribeBalance(user)
	if err != nil {
		ErrorResponse(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		&api.BalanceList{
			AdditionalProperties: balance,
		},
	)
}

// ErrorResponse attempts to read the status code from the error and returns the code and the message to the caller.
func ErrorResponse(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var code int

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

func convertPayerPointsSlice(payerPoints []PayerPoints) []api.PayerPoints {
	slice := []api.PayerPoints{}

	for _, pp := range payerPoints {
		slice = append(slice, api.PayerPoints{
			Payer:  pp.Payer,
			Points: pp.Amount,
		})
	}

	return slice
}
