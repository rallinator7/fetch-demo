package user

import (
	"container/heap"
	"sync"

	"github.com/rallinator7/fetch-demo/internal"
	"github.com/rallinator7/fetch-demo/internal/points"
)

type Repository struct {
	lock  sync.RWMutex
	users map[string]User
}

func NewRepository() Repository {
	return Repository{
		lock:  sync.RWMutex{},
		users: map[string]User{},
	}
}

func (repo *Repository) GetUser(id string) (User, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	user, exist := repo.users[id]
	if !exist {
		return User{}, &MissingUserError{User: id}
	}

	return user, nil
}

func (repo *Repository) AddUser(user internal.User) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	_, exist := repo.users[user.Id]
	if exist {
		return &DuplicateUserError{User: user.Id}
	}

	repo.users[user.Id] = New(user)

	return nil
}

func (repo *Repository) ListPoints(id string) (map[string]int, error) {
	repo.lock.RLock()
	defer repo.lock.RUnlock()

	user, exist := repo.users[id]
	if !exist {
		return nil, &MissingUserError{
			User: id,
		}
	}

	return user.payerTotals, nil
}

// GivePoints adds PayerTransactions to the user.  If a PayerTransaction results in the payer's total to go negative,
// then this function returns an error.
func (repo *Repository) GivePoints(id string, transaction points.PayerTransaction) error {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	user, exist := repo.users[id]
	if !exist {
		return &MissingUserError{
			User: id,
		}
	}

	total := user.payerTotals[transaction.Payer]
	final := total + transaction.Amount

	if final < 0 {
		return &NegativeTransactionError{Target: transaction.Payer}
	}

	heap.Push(
		&user.payerTransactions,
		&PQTransaction{
			PayerTransaction: transaction,
		},
	)

	user.payerTotals[transaction.Payer] = final
	user.totalPoints += transaction.Amount

	repo.users[user.Id] = user

	return nil
}

// SpendPoints removes points from a user based on the oldest transaction.  If the requested amount of points to be removed are greater than
// the amount of points a user has, the this returns an error.  Otherwise this returns a list of PayerPoints.
func (repo *Repository) SpendPoints(id string, total int) ([]points.PayerPoints, error) {
	repo.lock.Lock()
	defer repo.lock.Unlock()

	user, exist := repo.users[id]
	if !exist {
		return nil, &MissingUserError{
			User: id,
		}
	}

	if user.totalPoints < total {
		return nil, &NegativeTransactionError{
			Target: user.FirstName,
		}
	}

	user, pts := popUserTransactions(user, total)

	repo.users[user.Id] = user

	return payerMapToSlice(pts), nil
}

// popTransactions loops through a users transactions and removes oldest transactions until total is zero.
func popUserTransactions(user User, total int) (User, map[string]int) {
	pts := map[string]int{}

	for {
		transaction := heap.Pop(&user.payerTransactions).(*PQTransaction)

		if total > transaction.Amount && total-transaction.Amount >= 0 {
			total -= transaction.Amount
			user.totalPoints -= transaction.Amount
			pts[transaction.Payer] -= transaction.Amount
			user.payerTotals[transaction.Payer] -= transaction.Amount

			continue
		}

		// correct user totals
		user.totalPoints -= total
		user.payerTotals[transaction.Payer] -= total
		pts[transaction.Payer] -= total

		// correct transaction and place back on heap
		transaction.Amount -= total
		heap.Push(&user.payerTransactions, transaction)

		break
	}

	return user, pts
}

// payerMapToSlice converts a map of payer to amounts into a slice of PayerPoints.
func payerMapToSlice(pts map[string]int) []points.PayerPoints {
	ppSlice := []points.PayerPoints{}

	for k, v := range pts {
		ppSlice = append(ppSlice, points.PayerPoints{
			Payer:  k,
			Amount: v,
		})
	}

	return ppSlice
}
