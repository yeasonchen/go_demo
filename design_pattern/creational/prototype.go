package creational

import (
	"encoding/json"
	"time"
)

// 原型模式 - 当新建一个类耗费很大资源时（例如某些属性需要经过复杂的计算，或者通过hash得到），可以通过原型模式复制的产生对象
// 其实就是复制

// Keyword 搜索关键字
type Keyword struct {
	word      string
	visit     int
	UpdatedAt *time.Time
}

// Clone 这里使用序列化与反序列化的方式深拷贝
func (k *Keyword) Clone() *Keyword {
	var newKeyword Keyword
	b, _ := json.Marshal(k)
	err := json.Unmarshal(b, &newKeyword)
	if err != nil {
		return nil
	}
	return &newKeyword
}

// Keywords 关键字 map
type Keywords map[string]*Keyword

// Clone 复制一个新的 keywords
// updatedWords: 需要更新的关键词列表，由于从数据库中获取数据常常是数组的方式
func (words Keywords) Clone(updatedWords []*Keyword) Keywords {
	newKeywords := Keywords{}

	for k, v := range words {
		// 这里是浅拷贝，直接拷贝了地址
		newKeywords[k] = v
	}

	// 替换掉需要更新的字段，这里用的是深拷贝
	for _, word := range updatedWords {
		newKeywords[word.word] = word.Clone()
	}

	return newKeywords
}
