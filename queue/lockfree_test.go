package queue

import (
	"sync"
	"testing"
)

func TestLockFreeMultiProducerConsumer(t *testing.T) {
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

func BenchmarkLockFreeBench10Enqueues(b *testing.B) {
	benchLockFreeEnqueue(10, b)
}

func BenchmarkLockFreeBench20Enqueues(b *testing.B) {
	benchLockFreeEnqueue(20, b)
}

func BenchmarkLockFreeBench30Enqueues(b *testing.B) {
	benchLockFreeEnqueue(30, b)
}

func BenchmarkLockFreeBench40Enqueues(b *testing.B) {
	benchLockFreeEnqueue(40, b)
}

func BenchmarkLockFreeBench50Enqueues(b *testing.B) {
	benchLockFreeEnqueue(50, b)
}

func BenchmarkLockFreeBench60Enqueues(b *testing.B) {
	benchLockFreeEnqueue(60, b)
}

func BenchmarkLockFreeBench70Enqueues(b *testing.B) {
	benchLockFreeEnqueue(70, b)
}

func BenchmarkLockFreeBench80Enqueues(b *testing.B) {
	benchLockFreeEnqueue(80, b)
}

func BenchmarkLockFreeBench90Enqueues(b *testing.B) {
	benchLockFreeEnqueue(90, b)
}

func BenchmarkLockFreeBench100Enqueues(b *testing.B) {
	benchLockFreeEnqueue(100, b)
}

func BenchmarkLockFreeBench10Dequeues(b *testing.B) {
	benchLockFreeDequeue(10, b)
}

func BenchmarkLockFreeBench20Dequeues(b *testing.B) {
	benchLockFreeDequeue(20, b)
}

func BenchmarkLockFreeBench30Dequeues(b *testing.B) {
	benchLockFreeDequeue(30, b)
}

func BenchmarkLockFreeBench40Dequeues(b *testing.B) {
	benchLockFreeDequeue(40, b)
}

func BenchmarkLockFreeBench50Dequeues(b *testing.B) {
	benchLockFreeDequeue(50, b)
}

func BenchmarkLockFreeBench60Dequeues(b *testing.B) {
	benchLockFreeDequeue(60, b)
}

func BenchmarkLockFreeBench70Dequeues(b *testing.B) {
	benchLockFreeDequeue(70, b)
}

func BenchmarkLockFreeBench80Dequeues(b *testing.B) {
	benchLockFreeDequeue(80, b)
}

func BenchmarkLockFreeBench90Dequeues(b *testing.B) {
	benchLockFreeDequeue(90, b)
}

func BenchmarkLockFreeBench100Dequeues(b *testing.B) {
	benchLockFreeDequeue(100, b)
}

func benchLockFreeEnqueue(numOfOps int, b *testing.B) {
	var (
		q         = NewLockFreeQueue()
		enqueueWg sync.WaitGroup
	)

	for run := 0; run < b.N; run++ {
		enqueueWg.Add(numOfOps)
		for i := 0; i < numOfOps; i++ {
			go func(value int, wg *sync.WaitGroup) {
				var (
					val interface{} = value
				)
				defer wg.Done()
				_ = q.Enqueue(&val)
			}(i, &enqueueWg)
		}
		enqueueWg.Wait()
	}
}

func benchLockFreeDequeue(numOfOps int, b *testing.B) {
	for run := 0; run < b.N; run++ {
		var (
			q         = NewLockFreeQueue()
			enqueueWg sync.WaitGroup
			dequeueWg sync.WaitGroup
		)

		enqueueWg.Add(numOfOps)
		dequeueWg.Add(numOfOps)

		for i := 0; i < numOfOps; i++ {
			go func(value int, wg *sync.WaitGroup) {
				var (
					val interface{} = value
				)
				defer wg.Done()
				_ = q.Enqueue(&val)
			}(i, &enqueueWg)
		}

		enqueueWg.Wait()

		for i := 0; i < numOfOps; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				q.Dequeue()
			}(&dequeueWg)
		}

		dequeueWg.Wait()

	}

}
