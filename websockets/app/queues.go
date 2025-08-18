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
	In  Stack
	Out Stack
}

func (q *Queue) Enqueue(value UserInfo) {
	q.In.Push(value)
}

func (q *Queue) Dequeue() UserInfo {
	if len(q.Out) == 0 {
		for len(q.In) > 0 {
			q.Out.Push(q.In.Pop())
		}
	}

	return q.Out.Pop()
}
