package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("pushfront+pushback", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		require.Equal(t, 2, l.Len())

		item := l.Front() // 10
		require.Equal(t, 10, item.Value)
		require.Nil(t, item.Prev.Prev)
		require.Equal(t, 20, item.Next.Value)

		item = l.Back()
		require.Equal(t, 20, item.Value)
		require.Nil(t, item.Next.Next)
		require.Equal(t, 10, item.Prev.Value)

		l.PushBack(30) // [10, 20, 30]
		elems := make([]int, 0, l.Len())
		for i := l.Front(); i.Next != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{10, 20, 30}, elems)
	})

	t.Run("different types of elements", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		item := l.Front()
		require.Equal(t, 10, item.Value)

		l.PushFront(nil) // [nil, 10]
		item = l.Front()
		require.Nil(t, item.Value)

		l.PushFront(1.1) // [1.1, nil, 10]
		item = l.Front()
		require.Equal(t, 1.1, item.Value)

		l.PushFront("aaa") // ["aaa", 1.1, nil, 10]
		item = l.Front()
		require.Equal(t, "aaa", item.Value)
	})

	t.Run("remove", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		item := l.Front() // 10
		require.Equal(t, 10, item.Value)
		require.Nil(t, item.Prev.Prev)
		require.Equal(t, 30, item.Next.Value)
		require.Nil(t, item.Next.Next.Next)

		l.Remove(item)   // [30]
		item = l.Front() // 30
		l.Remove(item)   // []

		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("move to front", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]

		l.MoveToFront(l.Front()) // [10, 20, 30]
		item := l.Front()
		require.Equal(t, 10, item.Value)
		require.Nil(t, item.Prev.Prev)
		require.Equal(t, 20, item.Next.Value)
		require.Equal(t, 30, item.Next.Next.Value)

		l.MoveToFront(l.Back()) // [30, 10, 20]
		item = l.Front()
		require.Equal(t, 30, item.Value)
		require.Nil(t, item.Prev.Prev)
		require.Equal(t, 10, item.Next.Value)
		require.Equal(t, 20, item.Next.Next.Value)
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i.Next != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
