package queue

import (
	"sync"
	"testing"
)

func TestLockFreeSingleProducerConsumer(t *testing.T) {
	var (
		q             = NewLockFreeQueue()
		lock          sync.Mutex
		expectedQueue = make(chan int, 10)
		enqueueWg     sync.WaitGroup
		dequeueWg     sync.WaitGroup
	)

	if startLength := q.Length(); startLength != 0 {
		t.Log("queue lengths starts non 0")
		t.Log("expected : 0")
		t.Log("actual : ", startLength)
		t.Fail()
	}

	enqueueWg.Add(10)
	for i := 0; i < 10; i++ {
		go func(value int, expectedQueue chan int, lock *sync.Mutex, wg *sync.WaitGroup) {
			var (
				val interface{} = value
				err error
			)
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			if err = q.Enqueue(&val); err != nil {
				t.Error(err)
			}
			expectedQueue <- value
		}(i, expectedQueue, &lock, &enqueueWg)
	}

	enqueueWg.Wait()

	if midwayLength := q.Length(); midwayLength != 10 {
		t.Log("midway length error")
		t.Log("expected : 10")
		t.Log("actual : ", midwayLength)
		t.Fail()
	}

	dequeueWg.Add(10)
	for i := 0; i < 10; i++ {
		go func(expectedQueue chan int, lock *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			var (
				val *interface{}
				err error
			)
			lock.Lock()
			defer lock.Unlock()
			if val, err = q.Dequeue(); err != nil {
				t.Error(err)
			}
			if expected := <-expectedQueue; expected != *val {
				t.Log("out of order item dequeuing")
				t.Log("expected : ", expected)
				t.Log("actual : ", *val)
				t.Fail()
			}
		}(expectedQueue, &lock, &dequeueWg)
	}

	dequeueWg.Wait()

	if _, err := q.Dequeue(); err == nil || err.Error() != "empty queue error" {
		t.Log("no empty queue error")
		t.Fail()
	}

	if endLength := q.Length(); endLength != 0 {
		t.Log("queue lengths ends non 0")
		t.Log("expected : 0")
		t.Log("actual : ", endLength)
		t.Fail()
	}

}
