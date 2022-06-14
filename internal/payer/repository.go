package payer

import (
	"net/http"
	"sync"

	internal "github.com/rallinator7/fetch-demo/internal"
)

type Repository struct {
	lock   sync.RWMutex
	payers map[string]internal.Payer
}

func NewRepository() Repository {
	return Repository{
		lock:   sync.RWMutex{},
		payers: map[string]internal.Payer{},
	}
}

// Add attempts to add a payer to memory.  If a payer with the same name already exists, it returns an error.
func (repo *Repository) Add(payer internal.Payer) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	_, exist := repo.payers[payer.Name]
	if exist {
		return internal.NewStatusCodeError(http.StatusBadRequest, &DuplicatePayerError{Name: payer.Name})
	}

	repo.payers[payer.Name] = payer

	return nil
}
