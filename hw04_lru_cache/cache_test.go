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
				t.Helper()
				c := NewCache(10)

				val, ok := c.Get("aaa")
				require.False(t, ok)
				require.Nil(t, val)

				val, ok = c.Get("bbb")
				require.False(t, ok)
				require.Nil(t, val)
			},
		},
		{
			name: "add to cache",
			runfunc: func(t *testing.T) {
				t.Helper()
				c := NewCache(2)

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
				t.Helper()
				c := NewCache(2)

				c.Set("aaa", 100)

				wasInCache := c.Set("aaa", 1)
				require.True(t, wasInCache)

				val, wasInCache := c.Get("aaa")
				require.True(t, wasInCache)
				require.Equal(t, 1, val)
			},
		},
		{
			name: "overflow cache",
			runfunc: func(t *testing.T) {
				t.Helper()
				c := NewCache(2)

				c.Set("aaa", 1) // [1]
				c.Set("bbb", 2) // [2, 1]
				c.Set("ccc", 3) // [3, 2]

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
				t.Helper()
				c := NewCache(2)

				c.Set("aaa", 1)
				c.Set("bbb", 2)
				c.Clear()

				val, wasInCache := c.Get("aaa")
				require.False(t, wasInCache)
				require.Nil(t, val)

				val, wasInCache = c.Get("bbb")
				require.False(t, wasInCache)
				require.Nil(t, val)
			},
		},
		{
			name: "different types of elements",
			runfunc: func(t *testing.T) {
				t.Helper()
				c := NewCache(2)

				c.Set("aaa", 1)
				c.Set("это тест", "и это тест")

				val, wasInCache := c.Get("aaa")
				require.True(t, wasInCache)
				require.Equal(t, 1, val)

				val, wasInCache = c.Get("это тест")
				require.True(t, wasInCache)
				require.Equal(t, "и это тест", val)
			},
		},
		{
			name: "simple",
			runfunc: func(t *testing.T) {
				t.Helper()
				c := NewCache(5)

				wasInCache := c.Set("aaa", 100)
				require.False(t, wasInCache)

				wasInCache = c.Set("bbb", 200)
				require.False(t, wasInCache)

				val, wasInCache := c.Get("aaa") // [100, 200]
				require.True(t, wasInCache)
				require.Equal(t, 100, val)

				val, wasInCache = c.Get("bbb") // [200, 100]
				require.True(t, wasInCache)
				require.Equal(t, 200, val)

				wasInCache = c.Set("aaa", 300) // [300, 200]
				require.True(t, wasInCache)

				val, wasInCache = c.Get("aaa")
				require.True(t, wasInCache)
				require.Equal(t, 300, val)

				val, wasInCache = c.Get("ccc")
				require.False(t, wasInCache)
				require.Nil(t, val)
			},
		},
		{
			name: "purge logic",
			runfunc: func(t *testing.T) {
				t.Helper()
				c := NewCache(3)

				c.Set("aaa", 100)
				c.Set("bbb", 200)
				c.Set("ccc", 300)
				c.Set("bbb", 400) // [400, 300, 100]

				val, ok := c.Get("aaa") // [100, 400, 300]
				require.True(t, ok)
				require.Equal(t, 100, val)

				val, ok = c.Get("bbb") // [400, 100, 300]
				require.True(t, ok)
				require.Equal(t, 400, val)

				val, ok = c.Get("ccc") // [300, 400, 100]
				require.True(t, ok)
				require.Equal(t, 300, val)

				val, ok = c.Get("ddd") // [300, 400, 100]
				require.False(t, ok)
				require.Nil(t, val)

				c.Set("ddd", 500)
				c.Set("aaa", 600)

				val, ok = c.Get("aaa")
				require.True(t, ok)
				require.Equal(t, 600, val)

				val, ok = c.Get("bbb")
				require.False(t, ok)
				require.Nil(t, val)

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
	testcases := []struct {
		name    string
		runfunc func(t *testing.T)
	}{
		{
			name: "simple",
			runfunc: func(t *testing.T) {
				t.Helper()
				c := NewCache(1_0000)
				wg := &sync.WaitGroup{}
				var mu sync.Mutex

				for i := 0; i < 1_000; i++ {
					wg.Add(1)
					go func(i int) {
						defer wg.Done()
						mu.Lock()
						c.Set(Key("a"), i)
						val, _ := c.Get(Key("a"))
						require.Equal(t, i, val)
						mu.Unlock()
					}(i)
				}

				wg.Wait()
			},
		},
		{
			name: "complex",
			runfunc: func(t *testing.T) {
				t.Helper()
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
					var num int
					for i := 0; i < 1_000_000; i++ {
						num = rand.Intn(1_000_000)
						val, wasIsCache := c.Get(Key(strconv.Itoa(num)))
						if wasIsCache {
							require.Equal(t, num, val)
						}
					}
				}()

				wg.Wait()
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, tc.runfunc)
	}
}
