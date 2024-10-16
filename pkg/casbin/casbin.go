package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

type Casbin struct {
	enforcer *casbin.Enforcer
}

func NewCasbin(db *gorm.DB) *Casbin {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}

	m, err := model.NewModelFromString(`
		[request_definition]
		r = sub, obj, type

		[policy_definition]
		p = sub, obj, type

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.type == p.type
	`)
	if err != nil {
		panic(err)
	}

	enforcer, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		panic(err)
	}

	return &Casbin{
		enforcer: enforcer,
	}
}

// Enforce 判断是否具有权限
func (c *Casbin) Enforce(sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}

// AddPolicy 添加权限
func (c *Casbin) AddPolicy(sub, obj, act string) (bool, error) {
	return c.enforcer.AddPolicy(sub, obj, act)
}

// AddPolicies 批量添加权限
func (c *Casbin) AddPolicies(policies [][]string) (bool, error) {
	return c.enforcer.AddPoliciesEx(policies)
}

//批量修改
func (c *Casbin) UpdatePolicies(oldPolicies [][]string, newPolicies [][]string) (bool, error) {
	// 先删除旧策略
	_, err := c.enforcer.RemovePolicies(oldPolicies)
	if err != nil {
		return false, err
	}
	
	// 再添加新策略
	return c.enforcer.AddPolicies(newPolicies)
}

// RemovePolicy 删除权限
func (c *Casbin) RemovePolicy(sub, obj, act string) (bool, error) {
	return c.enforcer.RemovePolicy(sub, obj, act)
}

// AddRoleForUser 为用户添加角色
func (c *Casbin) AddRoleForUser(user, role string) (bool, error) {
	return c.enforcer.AddGroupingPolicy(user, role)
}

// 批量添加角色
func (c *Casbin) AddRolesForUser(user string, roles []string) (bool, error) {
	policies := make([][]string, len(roles))
	for i, role := range roles {
		policies[i] = []string{user, role}
	}
	return c.enforcer.AddGroupingPoliciesEx(policies)
}

// RemoveRoleForUser 为用户删除角色
func (c *Casbin) RemoveRoleForUser(user, role string) (bool, error) {
	return c.enforcer.RemoveGroupingPolicy(user, role)
}

//获取用户的角色
func (c *Casbin) GetRolesForUser(user string) ([]string, error) {
	return c.enforcer.GetRolesForUser(user)
}

// 获取角色的菜单
func (c *Casbin) GetMenusForRole(role string) ([][]string, error) {
	return c.enforcer.GetPermissionsForUser(role)
}
