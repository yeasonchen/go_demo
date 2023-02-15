package leetcode

type DLinkedNode[K, V any] struct {
	Key        K
	Value      V
	Prev, Next *DLinkedNode[K, V]
}

func InitLinkedNode[K, V any](key K, value V) *DLinkedNode[K, V] {
	return &DLinkedNode[K, V]{
		Key:   key,
		Value: value,
	}
}
