package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	testcases := []struct {
		name    string
		runfunc func(t *testing.T)
	}{
		{
			name: "empty list",
			runfunc: func(t *testing.T) {
				l := NewList()

				require.Equal(t, 0, l.Len())
				require.Nil(t, l.Front())
				require.Nil(t, l.Back())
			},
		},
		{
			name: "pushfront",
			runfunc: func(t *testing.T) {
				l := NewList()

				l.PushFront(10) // [10]
				require.Equal(t, 1, l.Len())
				item := l.Front()
				require.Equal(t, 10, item.Value)
				require.Nil(t, item.Prev.Prev)
				require.Nil(t, item.Next.Next)

			},
		},
		{
			name: "pushback",
			runfunc: func(t *testing.T) {
				l := NewList()

				l.PushBack(100) // [100]
				require.Equal(t, 1, l.Len())
				item := l.Front()
				require.Equal(t, 100, item.Value)
				require.Nil(t, item.Prev.Prev)
				require.Nil(t, item.Next.Next)

			},
		},
		{
			name: "pushfront+pushback",
			runfunc: func(t *testing.T) {
				l := NewList()

				l.PushFront(10) // [10]
				l.PushBack(100) // [10, 100]
				require.Equal(t, 2, l.Len())
				item := l.Front() // 10
				require.Equal(t, 10, item.Value)
				require.Nil(t, item.Prev.Prev)
				require.Equal(t, 100, item.Next.Value)
				require.Nil(t, item.Next.Next.Next)

			},
		},
		{
			name: "remove middle",
			runfunc: func(t *testing.T) {
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

			},
		},
		{
			name: "move to front",
			runfunc: func(t *testing.T) {
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
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.runfunc)
	}

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
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
