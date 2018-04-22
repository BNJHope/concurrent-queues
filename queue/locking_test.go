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

func BenchmarkLockingBench10Enqueues(b *testing.B) {
	benchLockingEnqueue(10, b)
}

func BenchmarkLockingBench20Enqueues(b *testing.B) {
	benchLockingEnqueue(20, b)
}

func BenchmarkLockingBench30Enqueues(b *testing.B) {
	benchLockingEnqueue(30, b)
}

func BenchmarkLockingBench40Enqueues(b *testing.B) {
	benchLockingEnqueue(40, b)
}

func BenchmarkLockingBench50Enqueues(b *testing.B) {
	benchLockingEnqueue(50, b)
}

func BenchmarkLockingBench60Enqueues(b *testing.B) {
	benchLockingEnqueue(60, b)
}

func BenchmarkLockingBench70Enqueues(b *testing.B) {
	benchLockingEnqueue(70, b)
}

func BenchmarkLockingBench80Enqueues(b *testing.B) {
	benchLockingEnqueue(80, b)
}

func BenchmarkLockingBench90Enqueues(b *testing.B) {
	benchLockingEnqueue(90, b)
}

func BenchmarkLockingBench100Enqueues(b *testing.B) {
	benchLockingEnqueue(100, b)
}

func BenchmarkLockingBench10Dequeues(b *testing.B) {
	benchLockingDequeue(10, b)
}

func BenchmarkLockingBench20Dequeues(b *testing.B) {
	benchLockingDequeue(20, b)
}

func BenchmarkLockingBench30Dequeues(b *testing.B) {
	benchLockingDequeue(30, b)
}

func BenchmarkLockingBench40Dequeues(b *testing.B) {
	benchLockingDequeue(40, b)
}

func BenchmarkLockingBench50Dequeues(b *testing.B) {
	benchLockingDequeue(50, b)
}

func BenchmarkLockingBench60Dequeues(b *testing.B) {
	benchLockingDequeue(60, b)
}

func BenchmarkLockingBench70Dequeues(b *testing.B) {
	benchLockingDequeue(70, b)
}

func BenchmarkLockingBench80Dequeues(b *testing.B) {
	benchLockingDequeue(80, b)
}

func BenchmarkLockingBench90Dequeues(b *testing.B) {
	benchLockingDequeue(90, b)
}

func BenchmarkLockingBench100Dequeues(b *testing.B) {
	benchLockingDequeue(100, b)
}

func benchLockingEnqueue(numOfOps int, b *testing.B) {
	var (
		q         = NewLockingQueue()
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

func benchLockingDequeue(numOfOps int, b *testing.B) {
	for run := 0; run < b.N; run++ {
		var (
			q         = NewLockingQueue()
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
