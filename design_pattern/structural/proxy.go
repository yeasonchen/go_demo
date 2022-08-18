package structural

import (
	"log"
	"time"
)

// 代理模式 - 建立一个代理类，内部包含被代理类的相同方法，该方法内会调用被代理类的相同方法
// 并且在此之上可以添加自定义逻辑，比如可以添加开始和结束时间，记录调用的时长

// IUser IUser
type IUser interface {
	Login(username, password string) error
}

// User 用户
type User struct{}

// Login 用户登录
func (u *User) Login(username, password string) error {
	// 不实现细节
	return nil
}

// UserProxy 代理类
type UserProxy struct {
	user *User
}

func NewUserProxy(user *User) *UserProxy {
	return &UserProxy{
		user: user,
	}
}

// Login 登录，和 user 实现相同的接口
func (p *UserProxy) Login(username, password string) error {
	// before 这里可能会有一些统计的逻辑
	start := time.Now()

	// 这里是原有的业务逻辑
	if err := p.user.Login(username, password); err != nil {
		return err
	}

	// after 这里可能也有一些监控统计的逻辑
	log.Printf("user login cost time: %s", time.Now().Sub(start))

	return nil
}
