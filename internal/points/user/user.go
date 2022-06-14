package user

import (
	"github.com/rallinator7/fetch-demo/internal"
)

type User struct {
	internal.User
	payerTransactions TransactionPriorityQueue
	payerTotals       map[string]int
	totalPoints       int
}

func New(user internal.User) User {
	return User{
		User:              user,
		payerTransactions: TransactionPriorityQueue{},
		payerTotals:       map[string]int{},
	}
}
