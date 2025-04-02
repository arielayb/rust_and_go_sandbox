package app

type Stack []any

func (s *Stack) Push(value any) {
	*s = append(*s, value)
}

func (s *Stack) Pop() any {
	result := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return result
}

type Queue struct {
	in  Stack
	out Stack
}

func (q *Queue) Enqueue(value any) {
	q.in.Push(value)
}

func (q *Queue) Dequeue() any {
	if len(q.out) == 0 {
		for len(q.in) > 0 {
			q.out.Push(q.in.Pop())
		}
	}

	return q.out.Pop()
}
