package linkedlist

import (
	"testing"
	"testing/quick"
)

func istOk(t *testing.T) {}

type data struct {
	key   uint
	value int
}

func TestInsert(t *testing.T) {

	validTests := []int{
		0,
		1,
		2,
		3,
		4,
	}

	invalidTests := []data{
		{
			6,
			17,
		},
		{
			10,
			18,
		},
		{
			12,
			18,
		},
		{
			9,
			12,
		},
	}

	l := New()

	for k, v := range validTests {
		ok := l.Insert(uint(k), v)
		if !ok {
			t.Fatalf("Error inserting item at: %d with value: %d", k, v)
		}
	}

	for k, v := range validTests {
		item, ok := l.Get(uint(k))
		if !ok {
			t.Fatalf("Item at index %d should exists", k)
		}
		if item != v {
			t.Fatalf("Item at index %d should be %d but is, %d", k, v, item)
		}
	}

	for _, v := range invalidTests {
		ok := l.Insert(v.key, v.value)
		if ok {
			t.Fatalf("Shouldn't be able to insert at %d", v.key)
		}
	}

	l.Insert(2, 100)

	validTests = append(validTests, 0)
	copy(validTests[3:], validTests[2:])
	validTests[2] = 100
	for k, v := range validTests {
		item, ok := l.Get(uint(k))
		if !ok {
			t.Fatalf("Item at index %d should exists", uint(k))
		}
		if item != v {
			t.Fatalf("Item at index %d should be %d but is, %d", uint(k), v, item)
		}
	}
	item, ok := l.Get(2)
	if !ok {
		t.Fatalf("Item at index %d should exists with value", 2)
	}
	if item != 100 {
		t.Fatalf("Item at index %d should be %d but is, %d", 2, 100, item)
	}

}

func TestRemove(t *testing.T) {
	validTests := []data{
		{
			0,
			0,
		},
		{
			1,
			1,
		},
		{
			2,
			2,
		},
		{
			3,
			3,
		},
		{
			4,
			4,
		},
	}

	removedKeys := []uint{
		1,
		3,
	}

	l := New()

	for _, v := range validTests {
		ok := l.Insert(v.key, v.value)
		if !ok {
			t.Fatalf("Error inserting item at: %d with value: %d", v.key, v.value)
		}
	}

	for _, k := range removedKeys {
		ok := l.Remove(uint(k))
		if !ok {
			t.Fatalf("Error removing item at: %d", k)
		}
	}

	validTests = append(validTests[:1], validTests[1+1:]...)
	validTests = append(validTests[:3], validTests[3+1:]...)

	for k, v := range validTests {
		item, ok := l.Get(uint(k))
		if !ok {
			t.Fatalf("Item at index %d should exists with value", 2)
		}
		if item != v.value {
			t.Fatalf("Item at index %d should be %d but is, %d", v.key, v.value, item)
		}
	}

}

func TestGet(t *testing.T) {
	validTests := []int{
		0,
		1,
		2,
		3,
		4,
	}

	l := New()

	for k, v := range validTests {
		ok := l.Insert(uint(k), v)
		if !ok {
			t.Fatalf("Error inserting item at: %d with value: %d", k, v)
		}
	}

	for k, v := range validTests {
		item, ok := l.Get(uint(k))
		if !ok {
			t.Fatalf("Can't Get item at index %d", k)
		}
		if item != v {
			t.Fatalf("Item at index %d should be %d but is %d", k, v, item)
		}
	}
}

func TestFind(t *testing.T) {
	insertTests := []int{
		0,
		1,
		2,
		3,
		4,
	}

	l := New()

	for k, v := range insertTests {
		ok := l.Insert(uint(k), v)
		if !ok {
			t.Fatalf("Error inserting item at: %d with value: %d", k, v)
		}
	}

	validTests := []data{
		{
			key:   1,
			value: 1,
		},
		{
			key:   3,
			value: 3,
		},
	}

	invalidTests := []int{5, 7}

	for _, v := range validTests {
		index, ok := l.Find(v.value)
		if !ok {
			t.Fatalf("Can't find item with value %d", v.value)
		}

		if index != v.key {
			t.Fatalf("Item with value %d should be in index %d but is %d", v.value, v.key, index)
		}
	}

	for _, v := range invalidTests {
		_, ok := l.Find(v)
		if ok {
			t.Fatalf("Item should not be found with value %d", v)
		}
	}

}

func TestPropertyBasedTest(t *testing.T) {

	err := quick.Check(func(inputs []int) bool {
		l := New()

		for k, v := range inputs {
			ok := l.Insert(uint(k), v)
			if !ok {
				return false
			}
		}

		for k, v := range inputs {
			out, ok := l.Get(uint(k))
			if !ok {
				return false
			}
			if out != v {
				return false
			}
		}

		for k, v := range inputs {
			index, found := l.Find(v)
			if !found {
				return false
			}
			if index != uint(k) {
				return false
			}
		}

		for v := uint(len(inputs)); v > 0; v-- {
			ok := l.Remove(v - 1)
			if !ok {
				t.Fatal(v)
			}
		}

		for k := range inputs {
			_, ok := l.Get(uint(k))
			if ok {
				return false
			}
		}

		return true
	}, nil)

	if err != nil {
		t.Fatal(err)
	}
}
