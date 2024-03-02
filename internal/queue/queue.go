package queue

type ListNode[T any] struct {
	Value    T
	Next     *ListNode[T]
	Priority int
}

type PriorityQueue[T any] struct {
	Start *ListNode[T]
	count int
}

func NewListNode[T any](val T, priority int) *ListNode[T] {
	return &ListNode[T]{
		Value:    val,
		Priority: priority,
	}
}

func NewPriorityQueue[T any]() *PriorityQueue[T] {
	return &PriorityQueue[T]{
		Start: nil,
		count: 0,
	}
}

func (q *PriorityQueue[T]) Insert(node *ListNode[T]) {
	q.count++
	if q.Start == nil {
		q.Start = node
		node.Next = nil
		return
	}

	if q.Start.Priority >= node.Priority {
		node.Next = q.Start
		q.Start = node
		return
	}

	cur := q.Start
	var prev *ListNode[T]
	for cur != nil && cur.Priority < node.Priority {
		prev = cur
		cur = cur.Next
	}
	node.Next = cur
	prev.Next = node
}

func (q *PriorityQueue[T]) Extract() *ListNode[T] {
	if q.Start == nil {
		return nil
	}
	q.count--
	node := q.Start
	q.Start = q.Start.Next
	return node
}

func (q *PriorityQueue[T]) GetCount() int {
	return q.count
}
