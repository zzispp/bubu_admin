package utils

// RoleBuilder 是一个菜单树构建器的单例
type RoleBuilder struct{}

var roleBuilder *RoleBuilder

// GetRoleBuilderInstance 返回RoleBuilder的单例实例
func GetRoleBuilderInstance() *RoleBuilder {
	if roleBuilder == nil {
		roleBuilder = &RoleBuilder{}
	}
	return roleBuilder
}

func (rb *RoleBuilder) RemoveDuplicateRoleIds(roleIds []string) []string {
	roleIdMap := make(map[string]bool)
	for _, roleId := range roleIds {
		roleIdMap[roleId] = true
	}
	uniqueRoleIds := make([]string, 0, len(roleIdMap))
	for roleId := range roleIdMap {
		uniqueRoleIds = append(uniqueRoleIds, roleId)
	}
	return uniqueRoleIds
}