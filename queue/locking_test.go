package queue

import (
	"sync"
	"testing"
)

func TestLockingMultiProducerConsumer(t *testing.T) {
	var (
		q             = NewLockingQueue()
		lock          sync.Mutex
		enqueueWg     sync.WaitGroup
		dequeueWg     sync.WaitGroup
		expectedQueue = make(chan int, 10)
	)

	if startLength := q.Length(); startLength != 0 {
		t.Log("queue lengths starts non 0")
		t.Log("expected : 0")
		t.Log("actual : ", startLength)
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		enqueueWg.Add(1)
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

	for i := 0; i < 10; i++ {
		dequeueWg.Add(1)
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

func BenchmarkLocking100Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 100, b)
}

func BenchmarkLocking200Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 200, b)
}

func BenchmarkLocking300Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 300, b)
}

func BenchmarkLocking400Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 400, b)
}

func BenchmarkLocking500Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 500, b)
}

func BenchmarkLocking600Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 600, b)
}

func BenchmarkLocking700Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 700, b)
}

func BenchmarkLocking800Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 800, b)
}

func BenchmarkLocking900Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 900, b)
}

func BenchmarkLocking1000Enqueues2Workers(b *testing.B) {
	benchLockingEnqueue(2, 1000, b)
}

func BenchmarkLocking100Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 100, b)
}

func BenchmarkLocking200Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 200, b)
}

func BenchmarkLocking300Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 300, b)
}

func BenchmarkLocking400Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 400, b)
}

func BenchmarkLocking500Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 500, b)
}

func BenchmarkLocking600Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 600, b)
}

func BenchmarkLocking700Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 700, b)
}

func BenchmarkLocking800Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 800, b)
}

func BenchmarkLocking900Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 900, b)
}

func BenchmarkLocking1000Dequeues2Workers(b *testing.B) {
	benchLockingDequeue(2, 1000, b)
}

func benchLockingEnqueue(numOfTasks int, numOfOps int, b *testing.B) {
	var (
		q         = NewLockingQueue()
		enqueueWg sync.WaitGroup
	)

	for run := 0; run < b.N; run++ {
		enqueueWg.Add(numOfTasks)
		for i := 0; i < numOfTasks; i++ {
			go func(value int, wg *sync.WaitGroup) {
				var (
					val interface{} = value
				)
				defer wg.Done()
				for j := 0; j < numOfOps; j++ {
					_ = q.Enqueue(&val)
				}
			}(i, &enqueueWg)
		}
		enqueueWg.Wait()
	}
}

func benchLockingDequeue(numOfTasks int, numOfOps int, b *testing.B) {
	for run := 0; run < b.N; run++ {
		var (
			q         = NewLockingQueue()
			enqueueWg sync.WaitGroup
			dequeueWg sync.WaitGroup
		)

		enqueueWg.Add(numOfTasks * numOfOps)
		dequeueWg.Add(numOfTasks)

		for i := 0; i < numOfTasks*numOfOps; i++ {
			go func(value int, wg *sync.WaitGroup) {
				var (
					val interface{} = value
				)
				defer wg.Done()
				_ = q.Enqueue(&val)
			}(i, &enqueueWg)
		}

		enqueueWg.Wait()

		for i := 0; i < numOfTasks; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				for j := 0; j < numOfOps; j++ {
					q.Dequeue()
				}
			}(&dequeueWg)
		}

		dequeueWg.Wait()

	}

}
