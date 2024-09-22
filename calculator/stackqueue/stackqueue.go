package stackqueue

type Item interface{}

type Stack struct {
	Data []Item
}

func (s *Stack) Push(value Item) {
	s.Data = append(s.Data, value)
}

func (s *Stack) Pop() (Item, bool) {
	if len(s.Data) == 0 {
		return nil, false
	}
	index := len(s.Data) - 1
	value := (s.Data)[index]
	s.Data = (s.Data)[:index]
	return value, true
}

func (s *Stack) Peek() (Item, bool) {
	if len(s.Data) == 0 {
		return nil, false
	}
	return (s.Data)[len(s.Data)-1], true
}

type Queue struct {
	Data []Item
}

func (q *Queue) Enqueue(value Item) {
	q.Data = append(q.Data, value)
}

func (q *Queue) Dequeue() (Item, bool) {
	if len(q.Data) == 0 {
		return nil, false
	}
	value := (q.Data)[0]
	q.Data = (q.Data)[1:]
	return value, true
}

func (q *Queue) First() (Item, bool) {
	if len(q.Data) == 0 {
		return nil, false
	}
	return (q.Data)[0], true
}
