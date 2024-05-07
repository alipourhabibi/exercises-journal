package linkedlist

type node struct {
	prev *node
	next *node
	Data int
}

type linkedList struct {
	head *node
}

func New() *linkedList {
	return &linkedList{}
}

func (l *linkedList) Insert(index uint, data int) bool {
	head := l.head
	node := &node{
		Data: data,
	}
	if head == nil {
		l.head = node
		return true
	}

	for i := 0; i < int(index)-1; i++ {
		if head.next == nil {
			return false
		}
		head = head.next
	}
	head.next = node
	return true
}

func (l *linkedList) Remove(index uint) bool {
	head := l.head
	prev := l.head
	if head == nil {
		return false
	}
	if index == 0 {
		l.head = head.next
		return true
	}
	for i := 0; i < int(index); i++ {
		if head.next == nil {
			return false
		}
		prev = head
		head = head.next
	}
	prev.next = head.next

	return true
}

func (l *linkedList) Find(n int) (index uint, found bool) {
	head := l.head
	index = 0
	for head != nil {
		if head.Data == n {
			return index, true
		}
		index++
		head = head.next
	}
	return 0, false
}

func (l *linkedList) Get(index uint) (int, bool) {
	head := l.head
	for i := 0; i < int(index); i++ {
		head = head.next
		if head == nil {
			return 0, false
		}
	}
	return head.Data, true
}
