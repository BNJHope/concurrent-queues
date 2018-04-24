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

func TestLockFreeSimultaneousWriterReader(t *testing.T) {
	var (
		q  = NewLockFreeQueue()
		wg sync.WaitGroup
	)

	wg.Add(200)

	for k := 0; k < 100; k++ {
		go func(q *LockFreeQueue, wg *sync.WaitGroup) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				var val interface{} = i
				q.Enqueue(&val)
			}
		}(q, &wg)
	}

	for k := 0; k < 100; k++ {
		go func(q *LockFreeQueue, wg *sync.WaitGroup) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				_, _ = q.Dequeue()
			}
		}(q, &wg)
	}

	wg.Wait()

}

func BenchmarkLockFree100Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 100, b)
}

func BenchmarkLockFree200Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 200, b)
}

func BenchmarkLockFree300Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 300, b)
}

func BenchmarkLockFree400Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 400, b)
}

func BenchmarkLockFree500Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 500, b)
}

func BenchmarkLockFree600Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 600, b)
}

func BenchmarkLockFree700Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 700, b)
}

func BenchmarkLockFree800Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 800, b)
}

func BenchmarkLockFree900Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 900, b)
}

func BenchmarkLockFree1000Enqueues2Workers(b *testing.B) {
	benchLockFreeEnqueue(2, 1000, b)
}

func BenchmarkLockFree100Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 100, b)
}

func BenchmarkLockFree200Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 200, b)
}

func BenchmarkLockFree300Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 300, b)
}

func BenchmarkLockFree400Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 400, b)
}

func BenchmarkLockFree500Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 500, b)
}

func BenchmarkLockFree600Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 600, b)
}

func BenchmarkLockFree700Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 700, b)
}

func BenchmarkLockFree800Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 800, b)
}

func BenchmarkLockFree900Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 900, b)
}

func BenchmarkLockFree1000Dequeues2Workers(b *testing.B) {
	benchLockFreeDequeue(2, 1000, b)
}

func benchLockFreeEnqueue(numOfTasks int, numOfOps int, b *testing.B) {
	var (
		q         = NewLockFreeQueue()
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

func benchLockFreeDequeue(numOfTasks int, numOfOps int, b *testing.B) {
	for run := 0; run < b.N; run++ {
		var (
			q         = NewLockFreeQueue()
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
