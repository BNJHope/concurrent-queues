package queue

import (
	"sync"
	"testing"
)

func TestLinkedLockingMultiProducerConsumer(t *testing.T) {
	var (
		q             = NewLockingLinkedQueue()
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

func BenchmarkLockingLinked100Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 100, b)
}

func BenchmarkLockingLinked200Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 200, b)
}

func BenchmarkLockingLinked300Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 300, b)
}

func BenchmarkLockingLinked400Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 400, b)
}

func BenchmarkLockingLinked500Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 500, b)
}

func BenchmarkLockingLinked600Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 600, b)
}

func BenchmarkLockingLinked700Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 700, b)
}

func BenchmarkLockingLinked800Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 800, b)
}

func BenchmarkLockingLinked900Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 900, b)
}

func BenchmarkLockingLinked1000Enqueues2Workers(b *testing.B) {
	benchLockingLinkedEnqueue(2, 1000, b)
}

func BenchmarkLockingLinked100Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 100, b)
}

func BenchmarkLockingLinked200Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 200, b)
}

func BenchmarkLockingLinked300Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 300, b)
}

func BenchmarkLockingLinked400Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 400, b)
}

func BenchmarkLockingLinked500Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 500, b)
}

func BenchmarkLockingLinked600Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 600, b)
}

func BenchmarkLockingLinked700Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 700, b)
}

func BenchmarkLockingLinked800Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 800, b)
}

func BenchmarkLockingLinked900Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 900, b)
}

func BenchmarkLockingLinked1000Dequeues2Workers(b *testing.B) {
	benchLockingLinkedDequeue(2, 1000, b)
}

func BenchmarkLockingLinked1Enqueues100Workers(b *testing.B) {
	benchLockingLinkedEnqueue(100, 1, b)
}

func BenchmarkLockingLinked1Enqueues200Workers(b *testing.B) {
	benchLockingLinkedEnqueue(200, 1, b)
}

func BenchmarkLockingLinked1Enqueues300Workers(b *testing.B) {
	benchLockingLinkedEnqueue(300, 1, b)
}

func BenchmarkLockingLinked1Enqueues400Workers(b *testing.B) {
	benchLockingLinkedEnqueue(400, 1, b)
}

func BenchmarkLockingLinked1Enqueues500Workers(b *testing.B) {
	benchLockingLinkedEnqueue(500, 1, b)
}

func BenchmarkLockingLinked1Enqueues600Workers(b *testing.B) {
	benchLockingLinkedEnqueue(600, 1, b)
}

func BenchmarkLockingLinked1Enqueues700Workers(b *testing.B) {
	benchLockingLinkedEnqueue(700, 1, b)
}

func BenchmarkLockingLinked1Enqueues800Workers(b *testing.B) {
	benchLockingLinkedEnqueue(800, 1, b)
}

func BenchmarkLockingLinked1Enqueues900Workers(b *testing.B) {
	benchLockingLinkedEnqueue(900, 1, b)
}

func BenchmarkLockingLinked1Enqueues1000Workers(b *testing.B) {
	benchLockingLinkedEnqueue(1000, 1, b)
}

func benchLockingLinkedEnqueue(numOfTasks int, numOfOps int, b *testing.B) {
	var (
		q         = NewLockingLinkedQueue()
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

func benchLockingLinkedDequeue(numOfTasks int, numOfOps int, b *testing.B) {
	for run := 0; run < b.N; run++ {
		var (
			q         = NewLockingLinkedQueue()
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
