package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	testcases := []struct {
		name    string
		runfunc func(t *testing.T)
	}{
		{
			name: "empty cache",
			runfunc: func(t *testing.T) {
				c := NewCache(10)

				_, ok := c.Get("aaa")
				require.False(t, ok)

				_, ok = c.Get("bbb")
				require.False(t, ok)
			},
		},
		{
			name: "add to cache",
			runfunc: func(t *testing.T) {
				c := NewCache(2)

				_, ok := c.Get("aaa")
				require.False(t, ok)

				wasInCache := c.Set("aaa", 100)
				require.False(t, wasInCache)

				val, ok := c.Get("aaa")
				require.True(t, ok)
				require.Equal(t, 100, val)
			},
		},
		{
			name: "rewrite item",
			runfunc: func(t *testing.T) {
				c := NewCache(2)

				wasInCache := c.Set("aaa", 100)
				require.False(t, wasInCache)

				wasInCache = c.Set("aaa", 1)
				require.True(t, wasInCache)

				val, ok := c.Get("aaa")
				require.True(t, ok)
				require.Equal(t, 1, val)
			},
		},
		{
			name: "cacheoverflow",
			runfunc: func(t *testing.T) {
				c := NewCache(2)

				c.Set("aaa", 1)
				c.Set("bbb", 2)
				c.Set("ccc", 3)

				val, ok := c.Get("aaa")
				require.False(t, ok)
				require.Nil(t, val)

				val, ok = c.Get("bbb")
				require.True(t, ok)
				require.Equal(t, 2, val)

				val, ok = c.Get("ccc")
				require.True(t, ok)
				require.Equal(t, 3, val)
			},
		},
		{
			name: "clear cache",
			runfunc: func(t *testing.T) {
				c := NewCache(2)

				c.Set("aaa", 1)
				c.Set("bbb", 2)
				c.Clear()

				val, ok := c.Get("aaa")
				require.False(t, ok)
				require.Nil(t, val)

				val, ok = c.Get("bbb")
				require.False(t, ok)
				require.Nil(t, val)
			},
		},
		{
			name: "simple",
			runfunc: func(t *testing.T) {
				c := NewCache(5)

				wasInCache := c.Set("aaa", 100)
				require.False(t, wasInCache)

				wasInCache = c.Set("bbb", 200)
				require.False(t, wasInCache)

				val, ok := c.Get("aaa")
				require.True(t, ok)
				require.Equal(t, 100, val)

				val, ok = c.Get("bbb")
				require.True(t, ok)
				require.Equal(t, 200, val)

				wasInCache = c.Set("aaa", 300)
				require.True(t, wasInCache)

				val, ok = c.Get("aaa")
				require.True(t, ok)
				require.Equal(t, 300, val)

				val, ok = c.Get("ccc")
				require.False(t, ok)
				require.Nil(t, val)
			},
		},
		{
			name: "purge logic",
			runfunc: func(t *testing.T) {
				c := NewCache(3)

				c.Set("aaa", 100)
				c.Set("bbb", 200)
				c.Set("ccc", 300)
				c.Set("bbb", 400)

				elems := c.GetQueueValues()
				require.Equal(t, 3, len(elems))
				require.Equal(t, []interface{}{400, 300, 100}, elems)
				fmt.Println("[test] ", c.GetItemsKeys(), elems)

				val, ok := c.Get("aaa") // [100, 400, 300]
				require.True(t, ok)
				require.Equal(t, 100, val)

				elems = c.GetQueueValues()
				fmt.Println("[test] ", c.GetItemsKeys(), elems)
				require.Equal(t, 3, len(elems))

				val, ok = c.Get("bbb") // [400, 100, 300]
				require.True(t, ok)
				require.Equal(t, 400, val)

				elems = c.GetQueueValues()
				fmt.Println("[test] ", c.GetItemsKeys(), elems)
				require.Equal(t, 3, len(elems))

				val, ok = c.Get("ccc") // [300, 400, 100]
				require.True(t, ok)
				require.Equal(t, 300, val)

				elems = c.GetQueueValues()
				fmt.Println("[test] ", c.GetItemsKeys(), elems)
				require.Equal(t, 3, len(elems))

				c.Set("ddd", 500)

				elems = c.GetQueueValues()
				fmt.Println("[test] ", c.GetItemsKeys(), elems)

				require.Equal(t, 3, len(elems))
				require.Equal(t, []interface{}{500, 300, 400}, elems)
				fmt.Println(elems)

				val, ok = c.Get("aaa")
				require.False(t, ok)
				require.Nil(t, val)

				val, ok = c.Get("bbb")
				require.True(t, ok)
				require.Equal(t, 400, val)

				val, ok = c.Get("ccc")
				require.True(t, ok)
				require.Equal(t, 300, val)

				val, ok = c.Get("ddd")
				require.True(t, ok)
				require.Equal(t, 500, val)

			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.runfunc)
	}
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
