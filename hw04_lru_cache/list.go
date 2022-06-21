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
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

func NewListItem(v interface{}, next, prev *ListItem) *ListItem {
	return &ListItem{
		Value: v,
		Next:  next,
		Prev:  prev,
	}
}

func NewList() List {
	newlist := new(list)
	//newback, newfront := NewListItem(nil, nil, nil), NewListItem(nil, nil, nil)
	newlist.front, newlist.back = NewListItem(nil, nil, nil), NewListItem(nil, nil, nil)
	newlist.back.Prev, newlist.front.Next = newlist.front, newlist.back
	//fmt.Println("front addr", newlist.front)
	//fmt.Println("front Next", newlist.front.Next, "front Prev", newlist.front.Prev)
	//fmt.Println("back  addr", newlist.back)
	//fmt.Println("back Next", newlist.back.Next, "back Prev", newlist.back.Prev)
	return newlist
}

type list struct {
	size  int
	front *ListItem
	back  *ListItem
}

func (l *list) Len() int {
	return l.size
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
	//fmt.Println("PushFront", newItem)
	l.front.Next.Prev = &newItem
	l.front.Next = &newItem
	//fmt.Println("front ", l.front)
	//fmt.Println("front Next", l.front.Next, "front Prev", l.front.Prev)
	//fmt.Println("back ", l.back)
	//fmt.Println("back Next", l.back.Next, "back Prev", l.back.Prev)
	l.size++
	return &newItem
}

func (l *list) PushBack(v interface{}) *ListItem {
	newItem := *NewListItem(v, l.back, l.back.Prev)
	//fmt.Println("PushBack", newItem)
	l.back.Prev.Next = &newItem
	l.back.Prev = &newItem
	//fmt.Println("front ", l.front)
	//fmt.Println("front Next", l.front.Next, "front Prev", l.front.Prev)
	//fmt.Println("back ", l.back)
	//fmt.Println("back Next", l.back.Next, "back Prev", l.back.Prev)
	l.size++
	return &newItem
}

func (l *list) Remove(i *ListItem) {
	// refactoring
	prevItem := i.Prev
	nextItem := i.Next

	prevItem.Next = nextItem
	nextItem.Prev = prevItem
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	_ = l.PushFront(i.Value)
	l.Remove(i)
}
