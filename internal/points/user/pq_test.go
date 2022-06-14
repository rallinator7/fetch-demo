package user

import (
	"container/heap"
	"math/rand"
	"testing"
	"time"

	"github.com/rallinator7/fetch-demo/internal/points"
	"github.com/stretchr/testify/assert"
)

func randomTimeTransaction(pts points.PayerPoints) points.PayerTransaction {
	now := time.Now()

	return points.PayerTransaction{
		PayerPoints: pts,
		TimeStamp:   now.Add(time.Minute * time.Duration(rand.Intn(1000000))),
	}
}

func Test(t *testing.T) {
	tests := map[string]struct {
		Transactions []points.PayerTransaction
	}{
		"sorts correctly": {
			Transactions: []points.PayerTransaction{
				randomTimeTransaction(points.PayerPoints{}),
				randomTimeTransaction(points.PayerPoints{}),
				randomTimeTransaction(points.PayerPoints{}),
				randomTimeTransaction(points.PayerPoints{}),
				randomTimeTransaction(points.PayerPoints{}),
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			pq := NewTransactionPriorityQueue()
			sorted := []points.PayerTransaction{}

			for _, t := range test.Transactions {
				heap.Push(
					pq,
					&PQTransaction{
						PayerTransaction: t,
					},
				)
			}

			for i := 0; i < len(test.Transactions); i++ {
				t := heap.Pop(pq).(*PQTransaction).PayerTransaction

				sorted = append(sorted, t)
			}

			for i := 0; i < len(sorted)-1; i++ {
				before := sorted[i].TimeStamp.Before(sorted[i+1].TimeStamp)
				assert.True(before)
			}
		})
	}
}
