package utils

import (
	"bubu_admin/internal/model"
	"strings"
	"sync"
)

// MenuTreeBuilder 是一个菜单树构建器的单例
type MenuTreeBuilder struct{}

var instance *MenuTreeBuilder
var once sync.Once

// GetMenuTreeBuilderInstance 返回MenuTreeBuilder的单例实例
func GetMenuTreeBuilderInstance() *MenuTreeBuilder {
	once.Do(func() {
		instance = &MenuTreeBuilder{}
	})
	return instance
}

// BuildMenuTree 构建菜单树结构
func (mtb *MenuTreeBuilder) BuildMenuTree(menus []*model.Menu) []*model.Menu {
	menuMap := make(map[string]*model.Menu)
	var rootMenus []*model.Menu

	// 第一次遍历，将所有菜单放入map中
	for _, menu := range menus {
		menuMap[menu.ID] = menu
		menu.Children = make([]*model.Menu, 0)
	}

	// 第二次遍历，构建树结构和当前路径
	for _, menu := range menus {
		if menu.ParentID == "root" {
			rootMenus = append(rootMenus, menu)
			if menu.Type != "button" {
				menu.CurrentPath = mtb.normalizePath(menu.Path)
			}
		} else if parent, exists := menuMap[menu.ParentID]; exists {
			parent.Children = append(parent.Children, menu)
			if menu.Type != "button" {
				menu.CurrentPath = mtb.buildCurrentPath(parent, menu)
			}
		}
	}

	// 递归排序所有层级的菜单
	mtb.sortMenus(rootMenus)

	return rootMenus
}

// buildCurrentPath 构建当前路径
func (mtb *MenuTreeBuilder) buildCurrentPath(parent *model.Menu, menu *model.Menu) string {
	parentPath := parent.CurrentPath
	menuPath := mtb.normalizePath(menu.Path)
	return strings.TrimSuffix(parentPath, "/") + "/" + strings.TrimPrefix(menuPath, "/")
}

// normalizePath 规范化路径
func (mtb *MenuTreeBuilder) normalizePath(path string) string {
	if path == "" {
		return "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// sortMenus 递归排序菜单
func (mtb *MenuTreeBuilder) sortMenus(menus []*model.Menu) {
	// 使用冒泡排序按Sequence升序排列
	for i := 0; i < len(menus)-1; i++ {
		for j := 0; j < len(menus)-1-i; j++ {
			if menus[j].Sequence > menus[j+1].Sequence {
				menus[j], menus[j+1] = menus[j+1], menus[j]
			}
		}
	}

	// 递归排序子菜单
	for _, menu := range menus {
		if len(menu.Children) > 0 {
			mtb.sortMenus(menu.Children)
		}
	}
}



// CollectMenusAndParents 给定等待添加，已有的菜单项和所有的菜单项，返回去除已有的菜单ID后，需要添加的子菜单ID，包括所有父级菜单
func (mtb *MenuTreeBuilder) CollectMenusAndParents(waitAddMenus []string, existingMenus []*model.RoleMenu, allMenus []*model.Menu) []string {
	// 创建用于存储需要添加的菜单ID的集合
	menusToAdd := make(map[string]bool)
	
	// 创建已存在菜单ID的集合，用于快速查找
	existingMenuSet := make(map[string]bool)
	for _, roleMenu := range existingMenus {
		existingMenuSet[roleMenu.MenuID] = true
	}

	// 创建菜单ID到菜单对象的映射，用于快速查找菜单信息
	menuMap := make(map[string]*model.Menu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}

	// 定义递归函数，用于收集菜单及其所有父级菜单
	var collectMenu func(string)
	collectMenu = func(menuId string) {
		// 如果菜单ID不在需要添加的集合中且不在已存在的集合中
		if !menusToAdd[menuId] && !existingMenuSet[menuId] {
			// 将菜单ID添加到需要添加的集合中
			menusToAdd[menuId] = true
			// 获取菜单对象
			menu, exists := menuMap[menuId]
			// 如果菜单存在且有父级菜单（不是根菜单）
			if exists && menu.ParentID != "" && menu.ParentID != "root" {
				// 递归收集父级菜单
				collectMenu(menu.ParentID)
			}
		}
	}

	// 对每个待添加的菜单ID执行收集操作
	for _, menuId := range waitAddMenus {
		collectMenu(menuId)
	}

	// 将需要添加的菜单ID集合转换为切片
	result := make([]string, 0, len(menusToAdd))
	for menuId := range menusToAdd {
		result = append(result, menuId)
	}

	// 返回需要添加的菜单ID切片
	return result
}