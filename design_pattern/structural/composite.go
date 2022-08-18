package structural

// 组合模式【用的很少】 - 当一组对象能够形成树形结构时，使用组合模式
// 类似于文件系统或者公司组织架构，将文件和文件夹都看作是一个结点

// IOrganization 组织接口，都实现统计人数的功能
type IOrganization interface {
	Count() int
}

// Employee 员工
type Employee struct {
	Name string
}

// Count 人数统计
func (Employee) Count() int {
	return 1
}

// Department 部门
type Department struct {
	Name string

	SubOrganizations []IOrganization
}

// Count 人数统计
func (d Department) Count() int {
	c := 0
	for _, org := range d.SubOrganizations {
		c += org.Count()
	}
	return c
}

// AddSub 添加子节点
func (d *Department) AddSub(org IOrganization) {
	d.SubOrganizations = append(d.SubOrganizations, org)
}

// NewOrganization 构建组织架构 demo
func NewOrganization() IOrganization {
	root := &Department{Name: "root"}
	for i := 0; i < 10; i++ {
		root.AddSub(&Employee{})
		root.AddSub(&Department{Name: "sub", SubOrganizations: []IOrganization{&Employee{}}})
	}
	return root
}
