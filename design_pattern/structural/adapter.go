package structural

import "fmt"

// 适配器模式 - 给原本不匹配使用的接口增加一个中间层，用于适配兼容
// 如下例子，AWS和阿里云各自有创建主机的方法，两者并不相同，为了统一使用
// 增加一个适配层接口ICreateServer，各自衍生一个适配类均实现CreateServer适配方法

// ICreateServer 创建云主机
type ICreateServer interface {
	CreateServer(cpu, mem float64) error
}

// AWSClient aws sdk
type AWSClient struct{}

// RunInstance 启动实例
func (c *AWSClient) RunInstance(cpu, mem float64) error {
	fmt.Printf("aws client run success, cpu： %f, mem: %f", cpu, mem)
	return nil
}

// AwsClientAdapter 适配器
type AwsClientAdapter struct {
	Client AWSClient
}

// CreateServer 启动实例
func (a *AwsClientAdapter) CreateServer(cpu, mem float64) error {
	_ = a.Client.RunInstance(cpu, mem)
	return nil
}

// AliyunClient aliyun sdk
type AliyunClient struct{}

// CreateServer 启动实例
func (c *AliyunClient) CreateServer(cpu, mem int) error {
	fmt.Printf("aws client run success, cpu： %d, mem: %d", cpu, mem)
	return nil
}

// AliyunClientAdapter 适配器
type AliyunClientAdapter struct {
	Client AliyunClient
}

// CreateServer 启动实例
func (a *AliyunClientAdapter) CreateServer(cpu, mem float64) error {
	_ = a.Client.CreateServer(int(cpu), int(mem))
	return nil
}
