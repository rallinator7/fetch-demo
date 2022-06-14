package payer

import (
	"sync"

	"github.com/rallinator7/fetch-demo/internal"
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

func (repo *Repository) GetPayer(id string) (internal.Payer, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	payer, exist := repo.payers[id]
	if !exist {
		return internal.Payer{}, &MissingUserError{Payer: id}
	}

	return payer, nil
}

func (repo *Repository) AddPayer(payer internal.Payer) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	_, exist := repo.payers[payer.Name]
	if exist {
		return &DuplicateUserError{Payer: payer.Name}
	}

	repo.payers[payer.Name] = payer

	return nil
}
