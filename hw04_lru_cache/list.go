package hw04lrucache

type List interface {
	Len() int                          // Длина списка
	Front() *ListItem                  // первый элемент списка
	Back() *ListItem                   // последний элемент списка
	PushFront(v interface{}) *ListItem // добавить значение в начало
	PushBack(v interface{}) *ListItem  // добавить значение в конец
	Remove(i *ListItem)                // удалить элемент
	MoveToFront(i *ListItem)           // переместить элемент в начало
}

type ListItem struct {
	Value interface{} // Значение
	Next  *ListItem   // Следующий элемент
	Prev  *ListItem   // Предыдущий элемент
}

func NewListItem(v interface{}, next, prev *ListItem) *ListItem {
	return &ListItem{
		Value: v,
		Next:  next,
		Prev:  prev,
	}
}

type list struct {
	length int       // Размер списка
	front  *ListItem // Первый элемент
	back   *ListItem // Последний элемент
}

func NewList() List {
	newlist := new(list)
	newlist.front, newlist.back = NewListItem(nil, nil, nil), NewListItem(nil, nil, nil)
	newlist.back.Prev, newlist.front.Next = newlist.front, newlist.back
	return newlist
}

func (l *list) Len() int {
	return l.length
}

func (l *list) Front() *ListItem {
	if l.front.Next == l.back {
		return nil
	}
	return l.front.Next
}

func (l *list) Back() *ListItem {
	if l.back.Prev == l.front {
		return nil
	}
	return l.back.Prev
}

func (l *list) PushFront(v interface{}) *ListItem {
	newItem := *NewListItem(v, l.front.Next, l.front)
	l.front.Next.Prev = &newItem
	l.front.Next = &newItem
	l.length++
	return &newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := *NewListItem(v, l.back, l.back.Prev)
	l.back.Prev.Next = &newItem
	l.back.Prev = &newItem
	l.length++
	return &newItem
}

func (l *list) Remove(i *ListItem) {
	i.Prev.Next = i.Next
	i.Next.Prev = i.Prev
	l.length--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)
	l.PushFront(i.Value)
}
