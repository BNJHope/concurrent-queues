package queue

type EmptyQueueError struct {
}

type Queue interface {
	Enqueue(payload *interface{}) error
	Dequeue() (*interface{}, error)
	Length() int
}

func (e EmptyQueueError) Error() string {
	return "empty queue error"
}
