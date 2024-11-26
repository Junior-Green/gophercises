package queue

type Queue[T any] struct {
	list []T
	len  int
}

func (q *Queue[T]) Push(node T) {
	q.list = append(q.list, node)
	q.len++
}

func (q *Queue[T]) Dequeue() T {
	node := q.list[0]
	q.list = q.list[1:]
	q.len--
	return node
}

func (q *Queue[T]) Length() int {
	return q.len
}
