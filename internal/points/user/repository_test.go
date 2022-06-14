package user

import (
	"testing"
	"time"

	"github.com/rallinator7/fetch-demo/internal"
	"github.com/rallinator7/fetch-demo/internal/points"
	"github.com/stretchr/testify/assert"
)

func TestRepository_ListPoints(t *testing.T) {
	tests := map[string]struct {
		PayerTotals map[string]int
	}{
		"returns correct map": {
			PayerTotals: map[string]int{
				"DANNON":       500,
				"UNILEVER":     0,
				"MILLER COORS": 200,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			user := internal.NewUser("Tester", "Last")
			repo := NewRepository()

			err := repo.AddUser(user)
			assert.NoError(err)

			pts, err := repo.ListPoints(user.Id)
			assert.NoError(err)

			for k, v := range pts {
				assert.Equal(test.PayerTotals[k], v)
			}
		})
	}
}

func TestRepository_GivePoints(t *testing.T) {
	tests := map[string]struct {
		Transactions      []points.PayerTransaction
		Err               error
		ExpectedTotals    map[string]int
		ExpectedUserTotal int
	}{
		"completes a transaction": {
			Transactions: []points.PayerTransaction{
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: 500,
					},
					TimeStamp: time.Now().UTC(),
				},
			},
			ExpectedTotals: map[string]int{
				"DANNON": 500,
			},
			ExpectedUserTotal: 500,
		},
		"returns err for negative total": {
			Transactions: []points.PayerTransaction{
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: -500,
					},
					TimeStamp: time.Now().UTC(),
				},
			},
			Err: &NegativeTransactionError{Target: "DANNON"},
		},
		"completes all transactions and allows for multiple payers": {
			Transactions: []points.PayerTransaction{
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: 500,
					},
					TimeStamp: time.Now().UTC(),
				},
				{
					PayerPoints: points.PayerPoints{
						Payer:  "UNILEVER",
						Amount: 300,
					},
					TimeStamp: time.Now().UTC(),
				},
			},
			ExpectedTotals: map[string]int{
				"DANNON":   500,
				"UNILEVER": 300,
			},
			ExpectedUserTotal: 800,
		},
		"completes all transactions and keeps correct total with negatives": {
			Transactions: []points.PayerTransaction{
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: 500,
					},
					TimeStamp: time.Now().UTC(),
				},
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: -200,
					},
					TimeStamp: time.Now().UTC(),
				},
			},
			ExpectedTotals: map[string]int{
				"DANNON": 300,
			},
			ExpectedUserTotal: 300,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			user := internal.NewUser("Tester", "Last")
			repo := NewRepository()
			err := repo.AddUser(user)
			assert.NoError(err)

			var expectedTotal int = 0

			for _, t := range test.Transactions {
				expectedTotal += t.Amount

				err = repo.GivePoints(user.Id, t)
				if test.Err != nil {
					assert.ErrorContains(err, test.Err.Error(), "err check")
					return
				}

				assert.NoError(err)
			}

			for k, v := range test.ExpectedTotals {
				assert.Equal(v, repo.users[user.Id].payerTotals[k])
			}

			assert.Equal(test.ExpectedUserTotal, expectedTotal)
			assert.Len(repo.users[user.Id].payerTransactions, len(test.Transactions))
		})
	}
}

func TestRepository_SpendPoints(t *testing.T) {
	tests := map[string]struct {
		Transactions                 []points.PayerTransaction
		Points                       int
		Err                          error
		ExpectedPayerPoints          map[string]int
		ExpectedUserTotal            int
		ExpectedLeftOverPayerPoints  map[string]int
		ExpectedLeftOverTransactions int
	}{
		"example from sheet": {
			Transactions: []points.PayerTransaction{
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: 1000,
					},
					TimeStamp: time.Now().UTC().Add(time.Minute * time.Duration(5)),
				},
				{
					PayerPoints: points.PayerPoints{
						Payer:  "UNILEVER",
						Amount: 200,
					},
					TimeStamp: time.Now().UTC().Add(time.Minute * time.Duration(2)),
				},
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: -200,
					},
					TimeStamp: time.Now().UTC().Add(time.Minute * time.Duration(3)),
				},
				{
					PayerPoints: points.PayerPoints{
						Payer:  "MILLER COORS",
						Amount: 10000,
					},
					TimeStamp: time.Now().UTC().Add(time.Minute * time.Duration(4)),
				},
				{
					PayerPoints: points.PayerPoints{
						Payer:  "DANNON",
						Amount: 300,
					},
					TimeStamp: time.Now().UTC().Add(time.Minute),
				},
			},
			Points: 5000,
			ExpectedPayerPoints: map[string]int{
				"DANNON":       -100,
				"UNILEVER":     -200,
				"MILLER COORS": -4700,
			},
			ExpectedLeftOverPayerPoints: map[string]int{
				"DANNON":       1000,
				"UNILEVER":     0,
				"MILLER COORS": 5300,
			},
			ExpectedLeftOverTransactions: 2,
		},
		"returns error for negative request": {
			Points: 5000,
			Err: &NegativeTransactionError{
				Target: "Tester",
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			user := internal.NewUser("Tester", "Last")
			repo := NewRepository()

			err := repo.AddUser(user)
			assert.NoError(err)

			for _, t := range test.Transactions {
				err = repo.GivePoints(user.Id, t)
				assert.NoError(err)
			}

			userTotal := repo.users[user.Id].totalPoints

			payerPoints, err := repo.SpendPoints(user.Id, test.Points)
			if test.Err != nil {
				assert.ErrorContains(err, test.Err.Error(), "err check")
				return
			}

			var totalSpent int = 0

			for _, pp := range payerPoints {
				assert.Equal(test.ExpectedPayerPoints[pp.Payer], pp.Amount)
				totalSpent += test.ExpectedPayerPoints[pp.Payer]
			}

			addedUser, err := repo.GetUser(user.Id)
			assert.NoError(err)

			for k, v := range addedUser.payerTotals {
				assert.Equal(test.ExpectedLeftOverPayerPoints[k], v)
			}

			assert.Equal(userTotal+totalSpent, addedUser.totalPoints)
			assert.Equal(test.ExpectedLeftOverTransactions, len(addedUser.payerTransactions))
		})
	}
}
