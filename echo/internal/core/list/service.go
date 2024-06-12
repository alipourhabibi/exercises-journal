package list

import (
	"sync"

	"github.com/alipourhabibi/exercises-journal/linkedlist"
)

type ListEntity struct {
	Index uint `json:"index"`
	Value int  `json:"value" validate:"required"`
}

type ListService struct {
	sync.Mutex
	linkedlist *linkedlist.LinkedList
}

type ListConfiguration func(*ListService) error

func New(cfgs ...ListConfiguration) (*ListService, error) {
	ls := &ListService{}

	for _, cfg := range cfgs {
		err := cfg(ls)
		if err != nil {
			return nil, err
		}
	}

	return ls, nil
}

func WithList(l *linkedlist.LinkedList) ListConfiguration {
	return func(ls *ListService) error {
		ls.linkedlist = l
		return nil
	}
}

func BootList() ListConfiguration {
	return func(ls *ListService) error {
		l := linkedlist.New()
		ls.linkedlist = l
		return nil
	}
}

func (l *ListService) Insert(index uint, value int) bool {
	l.Lock()
	defer l.Unlock()
	return l.linkedlist.Insert(index, value)
}

func (l *ListService) Remove(index uint) bool {
	l.Lock()
	defer l.Unlock()
	return l.linkedlist.Remove(index)
}

func (l *ListService) Find(value int) (uint, bool) {
	l.Lock()
	defer l.Unlock()
	return l.linkedlist.Find(value)
}

func (l *ListService) Get(index uint) (int, bool) {
	l.Lock()
	defer l.Unlock()
	return l.linkedlist.Get(index)
}
