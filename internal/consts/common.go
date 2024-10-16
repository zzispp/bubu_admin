package consts

//公共状态
const (
	StatusEnable  string = "enable"
	StatusDisable string = "disable"
)

//策略类型 menu/api
const (
	StrategyTypeMenu string = "menu"
	StrategyTypeApi  string = "api"
)

// 菜单状态
var ValidStatuses = []string{StatusEnable, StatusDisable}

// 验证状态是否有效
func IsValidStatus(s string) bool {
	for _, validStatus := range ValidStatuses {
		if s == validStatus {
			return true
		}
	}
	return false
}

