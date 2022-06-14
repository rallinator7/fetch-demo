package user

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	internal "github.com/rallinator7/fetch-demo/internal"
)

type Repository struct {
	lock  sync.RWMutex
	users map[string]internal.User
}

func NewRepository() Repository {
	return Repository{
		lock:  sync.RWMutex{},
		users: map[string]internal.User{},
	}
}

// Add attempts to add a user to memory.  If a user with the same name already exists, it returns an error.
func (repo *Repository) Add(user internal.User) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	mapper := fmt.Sprintf("%s-%s", strings.ToUpper(user.FirstName), strings.ToUpper(user.LastName))

	_, exist := repo.users[mapper]
	if exist {
		return internal.NewStatusCodeError(http.StatusBadRequest, &DuplicateUserError{First: user.FirstName, Last: user.LastName})
	}

	repo.users[mapper] = user

	return nil
}
