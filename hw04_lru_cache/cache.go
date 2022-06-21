package hw04lrucache

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
	_, ok := c.items[key]
	if ok {
		item := c.queue.PushFront(value)
		c.items[key] = item
	} else {
		item := c.queue.PushFront(value)
		c.items[key] = item
		if c.queue.Len() > c.capacity {
			deleteitem := c.queue.Back()
			c.queue.Remove(deleteitem)
			for deletekey, val := range c.items {
				if val.Value == deleteitem.Value {
					delete(c.items, deletekey)
				}
			}
		}
	}
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
		return item.Value, ok
	}
	return nil, ok
}

func (c *lruCache) Clear() {
	// TO DO
}
