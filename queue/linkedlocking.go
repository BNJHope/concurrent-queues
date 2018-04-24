package queue

import (
	"sync"
)

type LockingLinkedQueue struct {
	Front  *LinkedListNode
	Back   *LinkedListNode
	length int
	lock   *sync.RWMutex
}

func (q *LockingLinkedQueue) Enqueue(payload *interface{}) error {
	q.lock.Lock()
	defer q.lock.Unlock()
	newEnd := LinkedListNode{
		value: payload,
		Next:  nil,
	}
	if q.Back != nil {
		q.Back.SetNext(&newEnd)
	}
	q.Back = &newEnd
	if q.Front == nil {
		q.Front = &newEnd
	}
	q.length++
	return nil
}

func (q *LockingLinkedQueue) Dequeue() (*interface{}, error) {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.Front == nil {
		return nil, EmptyQueueError{}
	}
	res := q.Front
	q.Front = q.Front.Next
	q.length--
	return res.GetValue(), nil
}

func (q *LockingLinkedQueue) Length() int {
	return q.length
}

func NewLockingLinkedQueue() *LockingLinkedQueue {
	return &LockingLinkedQueue{
		Front:  nil,
		Back:   nil,
		length: 0,
		lock:   &sync.RWMutex{},
	}
}
