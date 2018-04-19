package queue

import (
	"sync"
)

type QueueElement interface {
}

type Queue interface {
	enqueue(payload *QueueElement)
	dequeue() *QueueElement
}

type LockingQueue struct {
	elements []*QueueElement
	lock     *sync.RWMutex
}

type LockFreeQueue struct {
}

func (q *LockingQueue) enqueue(payload *QueueElement) {
	q.lock.Lock()
	q.elements.append(payload)
	q.lock.Unlock()
}

func (q *LockingQueue) dequeue() (res *QueueElement) {
	q.lock.Lock()
	res = q[0]
	q = q[1:]
	q.lock.Unlock()
}
