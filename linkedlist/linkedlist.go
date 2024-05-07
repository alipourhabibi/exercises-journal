// A bad copy of https://github.com/cloudflare/pingora/blob/main/pingora-lru/src/linked_list.rs
// this linkedlist is not safe to use concurrently
// if you want to use it concurrently you should hanlde the locking yourself in you package

package linkedlist

import (
	"fmt"
	"log/slog"
	"math"
)

type Index uint64

const VEC_EXP_GROWTH_CAP uint64 = 65536

var (
	NIL    = Index(math.MaxUint64)
	HEAD   = &Node{Data: 0}
	TAIL   = &Node{Data: 1}
	OFFSET = Index(0)
)

type Node struct {
	prev *Node
	next *Node
	Data int
}

type Nodes struct {
	head      *Node
	tail      *Node
	dataNodes []*Node
}

func (n Nodes) getTail() *Node {
	return n.tail
}

func NewNodesWithCapacity(capacity uint64) *Nodes {
	return &Nodes{
		dataNodes: make([]*Node, 0, capacity),
	}
}

func (n Nodes) newNode(data int) Index {
	node := &Node{
		prev: nil,
		next: nil,
		Data: data,
	}

	/*
		if cap(n.dataNodes) > int(VEC_EXP_GROWTH_CAP) && cap(n.dataNodes)-len(n.dataNodes) < 2 {
			n.dataNodes = append(n.dataNodes[:0], n.dataNodes[:cap(n.dataNodes)/10]...)
		}
	*/

	n.dataNodes = append(n.dataNodes, node)
	return Index(len(n.dataNodes)-1) + OFFSET
}

func (n Nodes) byIndex(index uint) *Node {
	if Index(len(n.dataNodes)) < Index(index)-OFFSET {
		return nil
	}
	return n.dataNodes[Index(index)-OFFSET]
}

// Linked list
type LinkedList struct {
	nodes *Nodes
	free  []*Node // to keep track of freed node to be used again
}

func NewLinkedListCap(capacity uint64) *LinkedList {
	return &LinkedList{
		nodes: NewNodesWithCapacity(capacity),
		free:  []*Node{},
	}
}

func (l *LinkedList) newNode(data int) *Node {
	if len(l.free) > 0 {
		slog.Debug(
			"Getting from the free node pool",
			"data", data,
		)
		// get the first free node
		node := l.free[0]

		// remove it from free
		l.free = l.free[1:]

		// update the payload
		node.Data = data
		return node
	} else {
		return &Node{
			Data: data,
			prev: nil,
			next: nil,
		}
	}
}

func (l *LinkedList) Insert(index uint, data int) bool {
	if index > uint(len(l.nodes.dataNodes)) {
		slog.Error(
			"can't insert with index more than the len of datas",
			"index", index,
			"data", data,
			"len", len(l.nodes.dataNodes),
		)
		return false
	}

	node := l.newNode(data)
	// is not empty
	if len(l.nodes.dataNodes) != 0 {
		prev := l.nodes.byIndex(index - 1)
		prevNext := prev.next

		prev.next = node
		node.next = prevNext
		node.prev = prev
		l.nodes.dataNodes = append(l.nodes.dataNodes[:index], l.nodes.dataNodes[index-1:]...)
		l.nodes.dataNodes[index] = node
	} else {
		l.nodes.dataNodes = append(l.nodes.dataNodes, node)
	}

	return true
}

func (l *LinkedList) Remove(index uint) bool {
	node := l.nodes.byIndex(index)
	if node == nil {
		slog.Debug(
			"index does not exists",
			"index", index,
		)
		// does not exists
		return false
	}
	l.free = append(l.free, node)
	prev := node.prev
	next := node.next
	node.prev = nil
	node.next = nil
	l.nodes.dataNodes = append(l.nodes.dataNodes[:index], l.nodes.dataNodes[index+1:]...)

	// fmt.Println(prev, next)
	prev.next = next

	return true
}

func (l *LinkedList) Find(n int) (index uint, found bool) {
	for k, v := range l.nodes.dataNodes {
		if v != nil && v.Data == n {
			return uint(k), true
		}
	}
	return 0, false
}

func (l *LinkedList) Get(index uint) (int, bool) {
	if Index(len(l.nodes.dataNodes)) <= Index(index)-OFFSET {
		return 0, false
	}
	node := l.nodes.dataNodes[Index(index)-OFFSET]
	if node != nil {
		return node.Data, true
	}
	return 0, false
}

func (l *LinkedList) Show() {
	h := l.nodes.dataNodes[0]
	for h.next != nil {
		fmt.Printf("data: %d\n", h.Data)
		h = h.next
	}
	fmt.Println("data >", h.Data)
}

func (n *Node) String() string {
	return fmt.Sprintf("SELF: %p, PREV: %p, NEXT: %p", n, n.prev, n.next)
}
