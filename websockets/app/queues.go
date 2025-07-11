package app

type Stack []UserInfo

func (s *Stack) Push(value UserInfo) *Stack {
	*s = append(*s, value)
	return s
}

func (s *Stack) Pop() UserInfo {
	result := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]

	return result
}

type Queue struct {
	in  Stack
	out Stack
}

func (q *Queue) Enqueue(value UserInfo) {
	q.in.Push(value)
}

func (q *Queue) Dequeue() UserInfo {
	if len(q.out) == 0 {
		for len(q.in) > 0 {
			q.out.Push(q.in.Pop())
		}
	}

	return q.out.Pop()
}
