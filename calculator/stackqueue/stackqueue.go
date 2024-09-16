package stackqueue

type Item interface{}

type Stack []Item

func (s *Stack) Push(value Item) {
	*s = append(*s, value)
}

func (s *Stack) Pop() (Item, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	index := len(*s) - 1
	value := (*s)[index]
	*s = (*s)[:index]
	return value, true
}

func (s *Stack) Peek() (Item, bool) {
	if len(*s) == 0 {
		return nil, false
	}
	return (*s)[len(*s)-1], true
}

type Queue []Item

func (q *Queue) Enqueue(value Item) {
	*q = append(*q, value)
}

func (q *Queue) Dequeue() (Item, bool) {
	if len(*q) == 0 {
		return nil, false
	}
	value := (*q)[0]
	*q = (*q)[1:]
	return value, true
}

func (q *Queue) First() (Item, bool) {
	if len(*q) == 0 {
		return nil, false
	}
	return (*q)[0], true
}
