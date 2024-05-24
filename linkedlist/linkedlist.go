package linkedlist

type node struct {
	//prev *node
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
	node := &node{
		Data: data,
	}
	if index == 0 {
		node.next = l.head
		l.head = node
		return true
	}

	prev := l.head
	for i := 0; i < int(index)-1 && prev != nil; i++ {
		// list length is smaller than index
		/*
			if head.next == nil {
				return false
			}
			prev = head
			head = head.next
		*/
		prev = prev.next
	}

	if prev == nil {
		return false
	}

	node.next = prev.next
	prev.next = node
	// head.next = node

	return true
}

func (l *linkedList) Remove(index uint) bool {
	if l.head == nil {
		return false
	}
	if index == 0 {
		l.head = l.head.next
		return true
	}
	prev := l.head
	for i := uint(0); i < index-1 && prev != nil; i++ {
		prev = prev.next
	}

	if prev == nil || prev.next == nil {
		return false
	}
	prev.next = prev.next.next

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
	for i := uint(0); i < index && head != nil; i++ {
		head = head.next
	}
	if head == nil {
		return 0, false
	}
	return head.Data, true
}
