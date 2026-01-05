package queue

type Queue[T any] struct {
	ch chan T
}

func NewQueue[T any](bufferSize int) *Queue[T] {
	return &Queue[T]{
		ch: make(chan T, bufferSize),
	}
}
func (q *Queue[T]) Enqueue(item T) {
	q.ch <- item
}

// Dequeue retrieves from queue
func (q *Queue[T]) Dequeue() <-chan T {
	return q.ch
}

// close closes the queue channel
func (q *Queue[T]) Close() {
	close(q.ch)
}

func (q *Queue[T]) TryEnqueue(item T) bool {
	select {
	case q.ch <- item:
		return true
	default:
		return false
	}
}
