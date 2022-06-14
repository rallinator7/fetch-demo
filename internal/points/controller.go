package points

import (
	"net/http"

	"github.com/rallinator7/fetch-demo/internal"
)

type PayerRepoer interface {
	GetPayer(id string) (internal.Payer, error)
}

type UserRepoer interface {
	ListPoints(id string) (map[string]int, error)
	GivePoints(id string, transaction PayerTransaction) error
	SpendPoints(id string, total int) ([]PayerPoints, error)
}

type Controller struct {
	userRepo  UserRepoer
	payerRepo PayerRepoer
}

func NewController(user UserRepoer, payer PayerRepoer) Controller {
	return Controller{
		userRepo:  user,
		payerRepo: payer,
	}
}

func (controller *Controller) DescribeBalance(id string) (map[string]int, error) {
	balance, err := controller.userRepo.ListPoints(id)
	if err != nil {
		return nil, internal.NewStatusCodeError(http.StatusBadRequest, err)
	}

	return balance, nil
}

func (controller *Controller) GiveUserPoints(id string, transaction PayerTransaction) error {
	_, err := controller.payerRepo.GetPayer(transaction.Payer)
	if err != nil {
		return internal.NewStatusCodeError(http.StatusBadRequest, err)
	}

	err = controller.userRepo.GivePoints(id, transaction)
	if err != nil {
		return internal.NewStatusCodeError(http.StatusBadRequest, err)
	}

	return nil
}

func (controller *Controller) SpendUserPoints(id string, points int) ([]PayerPoints, error) {
	payerPoints, err := controller.userRepo.SpendPoints(id, points)
	if err != nil {
		return nil, internal.NewStatusCodeError(http.StatusBadRequest, err)
	}

	return payerPoints, nil
}
