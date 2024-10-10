package consts



const (
	MenuTypeNav    string = "nav"
	MenuTypePage   string = "page"
	MenuTypeButton string = "button"
)



const (
	PathTypeInternal string = "internal"
	PathTypeExternal string = "external"
)


// 菜单类型
var ValidMenuTypes = []string{MenuTypeNav, MenuTypePage, MenuTypeButton}


// 路径类型
var ValidPathTypes = []string{PathTypeInternal, PathTypeExternal}

// 验证菜单类型是否有效
func IsValidMenuType(t string) bool {
	for _, validType := range ValidMenuTypes {
		if t == validType {
			return true
		}
	}
	return false
}


// 验证路径类型是否有效
func IsValidPathType(t string) bool {
	for _, validType := range ValidPathTypes {
		if t == validType {
			return true
		}
	}
	return false
}
