package creational

// 简单工厂模式
// 由于 Go 本身是没有构造函数的，一般而言我们采用 NewName 的方式创建对象/接口，当它返回的是接口的时候，其实就是简单工厂模式

type IRuleConfigParser interface {
	Parse(data []byte)
}

type jsonRuleConfigParser struct{}

func (j jsonRuleConfigParser) Parse(data []byte) {}

type yamlRuleConfigParser struct{}

func (j yamlRuleConfigParser) Parse(data []byte) {}

func NewIRuleConfigParser(t string) IRuleConfigParser {
	if t == "json" {
		return jsonRuleConfigParser{}
	} else {
		return yamlRuleConfigParser{}
	}
}

// 工厂方法
// 当对象的创建逻辑比较复杂，不只是简单的 new 一下就可以，而是要组合其他类对象，做各种初始化操作的时候，
// 推荐使用工厂方法模式，将复杂的创建逻辑拆分到多个工厂类中，让每个工厂类都不至于过于复杂

// IRuleConfigParserFactory 工厂方法接口
type IRuleConfigParserFactory interface {
	CreateParser() IRuleConfigParser
}

// yamlRuleConfigParserFactory yamlRuleConfigParser 的工厂类
type yamlRuleConfigParserFactory struct {
}

func (y yamlRuleConfigParserFactory) CreateParser() IRuleConfigParser {
	return yamlRuleConfigParser{}
}

// jsonRuleConfigParserFactory jsonRuleConfigParser 的工厂类
type jsonRuleConfigParserFactory struct {
}

func (j jsonRuleConfigParserFactory) CreateParser() IRuleConfigParser {
	return jsonRuleConfigParser{}
}

// NewIRuleConfigParserFactory 用一个简单工厂封装产生工厂类
func NewIRuleConfigParserFactory(t string) IRuleConfigParserFactory {
	switch t {
	case "json":
		return jsonRuleConfigParserFactory{}
	case "yaml":
		return yamlRuleConfigParserFactory{}
	}
	return nil
}

// 抽象工厂方法 - 一个工厂方法可以创建相关联的多个类的时候就是抽象工厂模式，这个不太常用

type ISystemConfigParser interface {
	ParseSystem(data []byte)
}

// jsonSystemConfigParser jsonSystemConfigParser
type jsonSystemConfigParser struct{}

func (j jsonSystemConfigParser) ParseSystem(data []byte) {}

// IConfigParserFactory 工厂方法接口
type IConfigParserFactory interface {
	CreateRuleParser() IRuleConfigParser
	CreateSystemParser() ISystemConfigParser
}

type jsonConfigParserFactory struct{}

func (j jsonConfigParserFactory) CreateRuleParser() IRuleConfigParser {
	return jsonRuleConfigParser{}
}

func (j jsonConfigParserFactory) CreateSystemParser() ISystemConfigParser {
	return jsonSystemConfigParser{}
}
