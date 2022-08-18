package behavioral

// 迭代器模式 - 将集合对象的遍历操作从集合对象中分离出来，用一个单独的类来实现，实现单一原则

// Iterator 迭代器接口
type Iterator interface {
	HasNext() bool
	Next()
	// 获取当前元素，由于 Go 1.15 中还没有泛型，所以我们直接返回 interface{}
	CurrentItem() interface{}
}

// ArrayInt 数组
type ArrayInt []int

// Iterator 返回迭代器
func (a ArrayInt) Iterator() Iterator {
	return &ArrayIntIterator{
		arrayInt: a,
		index:    0,
	}
}

// ArrayIntIterator 数组迭代
type ArrayIntIterator struct {
	arrayInt ArrayInt
	index    int
}

// HasNext 是否有下一个
func (iter *ArrayIntIterator) HasNext() bool {
	return iter.index < len(iter.arrayInt)-1
}

// Next 游标加一
func (iter *ArrayIntIterator) Next() {
	iter.index++
}

// CurrentItem 获取当前元素
func (iter *ArrayIntIterator) CurrentItem() interface{} {
	return iter.arrayInt[iter.index]
}
