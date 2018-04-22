package queue

import (
	"sync/atomic"
	"unsafe"
)

type LinkedListNode struct {
	value *interface{}
	Next  *LinkedListNode
}

type LockFreeQueue struct {
	Front  *LinkedListNode
	Back   *LinkedListNode
	length int
}

func (q *LockFreeQueue) Enqueue(payload *interface{}) error {
	var (
		tailSwapped = false
		currTail    *LinkedListNode
	)

	newEnd := LinkedListNode{
		value: payload,
		Next:  nil,
	}

	for !tailSwapped {
		currTail = q.Back
		tailSwapped = atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&currTail.Next)),
			(unsafe.Pointer(nil)),
			(unsafe.Pointer(&newEnd)))
		if !tailSwapped {
			atomic.CompareAndSwapPointer(
				(*unsafe.Pointer)(unsafe.Pointer(&q.Back)),
				(unsafe.Pointer(currTail)),
				(unsafe.Pointer(currTail.Next)))
		}
	}

	_ = atomic.CompareAndSwapPointer(
		(*unsafe.Pointer)(unsafe.Pointer(&q.Back)),
		(unsafe.Pointer(currTail)),
		(unsafe.Pointer(&newEnd)))

	q.length++

	return nil
}

func (q *LockFreeQueue) Dequeue() (*interface{}, error) {
	var (
		elementRetrieved = false
		result           *LinkedListNode
	)

	for !elementRetrieved {
		if result = q.Front; result.Next == nil {
			return nil, EmptyQueueError{}
		}

		elementRetrieved = atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&q.Front)),
			unsafe.Pointer(result),
			unsafe.Pointer(result.Next))
	}
	q.length--
	return result.Next.value, nil
}

func (q *LockFreeQueue) Length() int {
	return q.length
}

func (n *LinkedListNode) SetNext(next *LinkedListNode) {
	n.Next = next
}

func (n *LinkedListNode) GetValue() *interface{} {
	return n.value
}

func (n *LinkedListNode) SetValue(val *interface{}) {
	n.value = val
}

func NewLockFreeQueue() *LockFreeQueue {
	dummy := LinkedListNode{
		value: nil,
		Next:  nil,
	}

	return &LockFreeQueue{
		Front:  &dummy,
		Back:   &dummy,
		length: 0,
	}
}
