package user

import (
	"container/heap"

	"github.com/rallinator7/fetch-demo/internal/points"
)

type PriorityQueuer interface {
	Len() int
	Less(int, int) bool
	Swap(int, int)
	Push(interface{})
	Pop() interface{}
	Empty()
}

type PQTransaction struct {
	points.PayerTransaction
	index int
}

type TransactionPriorityQueue []*PQTransaction

func NewTransactionPriorityQueue() *TransactionPriorityQueue {
	pq := &TransactionPriorityQueue{}
	heap.Init(pq)

	return pq
}

func (pq TransactionPriorityQueue) Len() int { return len(pq) }

func (pq TransactionPriorityQueue) Less(i, j int) bool {
	return pq[i].TimeStamp.Before(pq[j].TimeStamp)
}

// needed to implemented heap interface and is used for swapping nodes in queue
func (pq TransactionPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

// retrieves a lock and adds a node into the queue
func (pq *TransactionPriorityQueue) Push(new interface{}) {
	n := len(*pq)
	transaction := new.(*PQTransaction)
	transaction.index = n
	*pq = append(*pq, transaction)
}

// retrieves a lock and pops highest priority node from queue
func (pq *TransactionPriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	transaction := old[n-1]
	old[n-1] = nil         // avoid memory leak
	transaction.index = -1 // for safety
	*pq = old[0 : n-1]
	return transaction
}
