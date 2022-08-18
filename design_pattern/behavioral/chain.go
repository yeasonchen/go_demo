package behavioral

// 职责链模式 - 类似于web框架中的filter
// 将一个事件的发生者和接收者解耦，将多个接收者链式存储，每个接收者都有机会处理事件
// 该模式能够实现扩展，如web框架里面，不修改源码也能使用户的逻辑得以执行

// SensitiveWordFilter 敏感词过滤器，判定是否是敏感词
type SensitiveWordFilter interface {
	Filter(content string) bool
}

// SensitiveWordFilterChain 职责链
type SensitiveWordFilterChain struct {
	filters []SensitiveWordFilter
}

// AddFilter 添加一个过滤器
func (c *SensitiveWordFilterChain) AddFilter(filter SensitiveWordFilter) {
	c.filters = append(c.filters, filter)
}

// Filter 执行过滤
func (c *SensitiveWordFilterChain) Filter(content string) bool {
	for _, filter := range c.filters {
		// 如果发现敏感直接返回结果
		if filter.Filter(content) {
			return true
		}
	}
	return false
}

// AdSensitiveWordFilter 广告
type AdSensitiveWordFilter struct{}

// Filter 实现过滤算法
func (f *AdSensitiveWordFilter) Filter(content string) bool {
	// TODO: 实现算法
	return false
}

// PoliticalWordFilter 政治敏感
type PoliticalWordFilter struct{}

// Filter 实现过滤算法
func (f *PoliticalWordFilter) Filter(content string) bool {
	// TODO: 实现算法
	return true
}
