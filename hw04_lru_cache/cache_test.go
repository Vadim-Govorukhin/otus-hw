package hw04lrucache

import (
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
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.runfunc)
	}

	t.Run("purge logic", func(t *testing.T) {
		// Write me
	})
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
