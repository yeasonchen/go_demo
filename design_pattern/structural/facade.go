package structural

// 外观模式 - 将系统内的复杂逻辑屏蔽起来，暴露一个简单的接口让用户使用
// 如下例子，系统内部有登录和注册两个逻辑，但是只暴露给用户一个注册或登录的逻辑，
// 用户已注册就会自动调用登录，未注册就会调用注册方法方便用户使用

type FacadeIUser interface {
	Login(phone int, code int) (*User, error)
	Register(phone int, code int) (*User, error)
}

// IUserFacade 门面模式
type IUserFacade interface {
	LoginOrRegister(phone int, code int) error
}

type FacadeUser struct {
	Name string
}

type UserService struct{}

// Login 登录
func (u UserService) Login(phone int, code int) (*FacadeUser, error) {
	// 校验操作 ...
	return &FacadeUser{Name: "test login"}, nil
}

// Register 注册
func (u UserService) Register(phone int, code int) (*FacadeUser, error) {
	// 校验操作 ...
	// 创建用户
	return &FacadeUser{Name: "test register"}, nil
}

// LoginOrRegister 登录或注册
func (u UserService) LoginOrRegister(phone int, code int) (*FacadeUser, error) {
	user, err := u.Login(phone, code)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return user, nil
	}

	return u.Register(phone, code)
}
