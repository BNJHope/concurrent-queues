package queue

import (
	"testing"
	"time"
)

func TestLockingSingleProducerConsumer(t *testing.T) {
	var (
		val interface{}
		q   = NewLockingQueue()
	)

	for i := 0; i < 10; i++ {
		val = i
		go q.Enqueue(&val)
	}

	time.Sleep(time.Second)

	expected := 10
	actual := q.Length()
	if expected != actual {
		t.Log("enqueueing : length failed")
		t.Log("expected : ", expected)
		t.Log("actual : ", actual)
		t.Fail()
	}

	for i := 0; i < 10; i++ {
		go q.Dequeue()
	}

	time.Sleep(time.Second)

	expected = 0
	actual = q.Length()
	if q.Length() != 0 {
		t.Log("dequeueing : length failed")
		t.Log("expected : ", expected)
		t.Log("actual : ", actual)
		t.Fail()
	}

}
