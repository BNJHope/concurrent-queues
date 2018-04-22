package queue

import (
	"sync"
)

type LockingQueue struct {
	elements []*interface{}
	lock     *sync.RWMutex
}

func (q *LockingQueue) Enqueue(payload *interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.elements = append(q.elements, payload)
	return nil
}

func (q *LockingQueue) Dequeue() (*interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if len(q.elements) == 0 {
		return nil, EmptyQueueError{}
	}
	res := q.elements[0]
	q.elements = q.elements[1:]
	return res, nil
}

func (q *LockingQueue) Length() int {
	return len(q.elements)
}

func NewLockingQueue() *LockingQueue {
	return &LockingQueue{
		elements: make([]*interface{}, 0),
		lock:     &sync.RWMutex{},
	}
}
