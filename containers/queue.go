package containers

type Queue[T any] []T

func (q *Queue[T]) IsEmpty() bool {
    return len(*q) == 0
}

func (q *Queue[T]) Size() int {
    return len(*q)
}

func (q *Queue[T]) Enqueue(item T) {
    *q = append(*q, item)
}

func (q *Queue[T]) Dequeue() (T, bool) {
    if q.IsEmpty() {
        var empty T
        return empty, false
    } else {
        item := (*q)[0]
        *q = (*q)[1:]
        return item, true
    }
}

func (q *Queue[T]) Peek() (T, bool) {
    if q.IsEmpty() {
        var empty T
        return empty, false
    } else {
        return (*q)[0], true
    }
}

func (q *Queue[T]) Front() T {
    return (*q)[0]
}