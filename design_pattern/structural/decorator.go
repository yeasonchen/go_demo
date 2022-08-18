package structural

// 装饰模式 - 组合替代继承
// 当要给一个类扩展功能时，可以使用组合，在其原有功能上添加新功能
// 和代理的区别，代理添加的是无关的功能，装饰的是相关的功能

// IDraw IDraw
type IDraw interface {
	Draw() string
}

// Square 正方形
type Square struct{}

func (s Square) Draw() string {
	return "this is a square"
}

// ColorSquare 有颜色的正方形
type ColorSquare struct {
	square IDraw
	color  string
}

func NewColorSquare(square IDraw, color string) ColorSquare {
	return ColorSquare{color: color, square: square}
}

func (c ColorSquare) Draw() string {
	return c.square.Draw() + ", color is " + c.color
}
