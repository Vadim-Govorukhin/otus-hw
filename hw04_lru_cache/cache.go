package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш.
}

type lruCache struct {
	capacity int               // Количество сохраняемых в кэше элементов
	queue    List              // очередь [последних используемых элементов] на основе двусвязного списка
	items    map[Key]*ListItem // словарь, отображающий ключ (строка) на элемент очереди
	mu       sync.Mutex        // Для безопасности горутин
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, inCache := c.items[key]

	if inCache {
		c.items[key].Value = cacheItem{key, value}
		c.queue.MoveToFront(c.items[key])

		return inCache
	}

	if c.queue.Len() == c.capacity {
		deleteitem := c.queue.Back()
		c.queue.Remove(deleteitem)
		delete(c.items, deleteitem.Value.(cacheItem).key)
	}

	item := c.queue.PushFront(cacheItem{key, value})
	c.items[key] = item
	return inCache
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value.(cacheItem).value, ok
	}
	return nil, ok
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) GetQueueValues() []interface{} {
	var elems []interface{}
	for i := c.queue.Front(); i.Next != nil; i = i.Next {
		elems = append(elems, i.Value.(cacheItem).value)
	}
	return elems
}
