package hw04lrucache

import "fmt"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш.

	GetQueueValues() []interface{} // Для тестов
	GetItemsKeys() []string        // Для тестов
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
		fmt.Println("[set] update map: ", key, value)
		c.items[key].Value = value
		c.queue.MoveToFront(c.items[key])

		return ok
	}
	fmt.Println("[set] add to map: ", key, value)
	item := c.queue.PushFront(value)
	c.items[key] = item
	if c.queue.Len() > c.capacity {
		deleteitem := c.queue.Back()
		c.queue.Remove(deleteitem)
		for deletekey, val := range c.items {
			if val == deleteitem {
				fmt.Println("[set] delete from map: ", deletekey, deleteitem.Value)
				delete(c.items, deletekey)
				break
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
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func (c *lruCache) GetQueueValues() []interface{} {
	var elems []interface{}
	for i := c.queue.Front(); i.Next != nil; i = i.Next {
		elems = append(elems, i.Value)
	}
	return elems
}

func (c *lruCache) GetItemsKeys() []string {
	var keys []string
	for key, _ := range c.items {
		keys = append(keys, string(key))
	}
	return keys
}
