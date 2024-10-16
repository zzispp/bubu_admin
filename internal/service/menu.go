package service

import (
	v1 "bubu_admin/api/v1"
	"bubu_admin/internal/consts"
	"bubu_admin/internal/model"
	"bubu_admin/internal/repository"
	"bubu_admin/internal/utils"
	"context"
	"fmt"
)

type MenuService interface {
	CreateMenu(ctx context.Context, req *v1.CreateMenuRequest) error
	ListMenu(ctx context.Context, req *v1.ListMenuRequest) ([]*model.Menu, error)
	UpdateMenu(ctx context.Context, id string, req *v1.UpdateMenuRequest) error
	GetMenuByID(ctx context.Context, id string) (*model.Menu, error)
	DeleteMenu(ctx context.Context, id string) error
}

func NewMenuService(
	service *Service,
	menuRepository repository.MenuRepository,
	roleRepository repository.RoleRepository,
) MenuService {
	return &menuService{
		Service:        service,
		menuRepository: menuRepository,
		roleRepository: roleRepository,
	}
}

type menuService struct {
	*Service
	menuRepository repository.MenuRepository
	roleRepository repository.RoleRepository
}

func (s *menuService) CreateMenu(ctx context.Context, req *v1.CreateMenuRequest) error {

	if req.Type == consts.MenuTypeButton {
		parent, err := s.menuRepository.FindByID(ctx, req.ParentID)
		if err != nil {
			return err
		}
		if parent.Type != "page" {
			return fmt.Errorf("按钮只能添加在页面类型菜单下")
		}
	}
	// 检查父级ID的合法性
	if req.ParentID != "root" {
		parentMenu, err := s.menuRepository.FindByID(ctx, req.ParentID)
		if err != nil {
			return fmt.Errorf("无效的父级ID")
		}
		if parentMenu == nil {
			return fmt.Errorf("父级菜单不存在")
		}
	}

	// 检查指定的code是否已存在于指定的parentID下
	existingMenu, err := s.menuRepository.FindByCodeAndParentID(ctx, req.Code, req.ParentID)
	if err != nil {
		return fmt.Errorf("检查菜单代码是否存在时发生错误")
	}
	if existingMenu != nil {
		return fmt.Errorf("该父级菜单下已存在相同代码的菜单")
	}

	if _, err := s.menuRepository.Save(ctx, &model.Menu{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Sequence:    req.Sequence,
		Type:        req.Type,
		Component:   req.Component,
		Redirect:    req.Redirect,
		Path:        req.Path,
		Icon:        req.Icon,
		Status:      req.Status,
		ParentID:    req.ParentID,
	}); err != nil {
		return fmt.Errorf("创建菜单失败")
	}

	return nil
}

func (s *menuService) ListMenu(ctx context.Context, req *v1.ListMenuRequest) ([]*model.Menu, error) {
	menus, err := s.menuRepository.List(ctx, &model.Menu{
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("获取菜单列表失败")
	}

	// 将菜单列表转换为树结构
	menuTree := utils.GetMenuTreeBuilderInstance().BuildMenuTree(menus)
	return menuTree, nil
}

func (s *menuService) UpdateMenu(ctx context.Context, id string, req *v1.UpdateMenuRequest) error {
	// 首先获取现有菜单
	existingMenu, err := s.menuRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取菜单失败: %w", err)
	}
	if existingMenu == nil {
		return fmt.Errorf("菜单不存在")
	}

	// 更新菜单字段
	existingMenu.Code = req.Code
	existingMenu.Name = req.Name
	existingMenu.Description = req.Description
	existingMenu.Sequence = req.Sequence
	existingMenu.Type = req.Type
	existingMenu.Component = req.Component
	existingMenu.Redirect = req.Redirect
	existingMenu.Path = req.Path
	existingMenu.Icon = req.Icon
	existingMenu.Status = req.Status
	existingMenu.ParentID = req.ParentID

	// 保存更新后的菜单
	if _, err := s.menuRepository.Update(ctx, existingMenu); err != nil {
		return fmt.Errorf("更新菜单失败: %w", err)
	}

	return nil
}

func (s *menuService) GetMenuByID(ctx context.Context, id string) (*model.Menu, error) {
	menu, err := s.menuRepository.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取菜单失败: %w", err)
	}
	if menu == nil {
		return nil, fmt.Errorf("菜单不存在")
	}
	return menu, nil
}

func (s *menuService) DeleteMenu(ctx context.Context, id string) error {
	// 首先检查菜单是否存在
	existingMenu, err := s.menuRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取菜单失败: %w", err)
	}
	if existingMenu == nil {
		return fmt.Errorf("菜单不存在")
	}

	// 检查是否有子菜单
	childMenus, err := s.menuRepository.List(ctx, &model.Menu{ParentID: id})
	if err != nil {
		return fmt.Errorf("检查子菜单失败: %w", err)
	}
	if len(childMenus) > 0 {
		return fmt.Errorf("无法删除包含子菜单的菜单")
	}

	// 获取所有角色
	roles, err := s.roleRepository.List(ctx, &model.Role{})
	if err != nil {
		return fmt.Errorf("获取角色列表失败: %w", err)
	}

	// 检查每个角色的菜单权限
	for _, role := range roles {
		menus, err := s.casbin.GetMenusForRole(role.ID)
		if err != nil {
			return fmt.Errorf("获取角色菜单权限失败: %w", err)
		}

		for _, menu := range menus {
			if menu[1] == id {
				return fmt.Errorf("无法删除已授权给角色 %s 的菜单", role.Name)
			}
		}
	}

	// 删除菜单
	err = s.menuRepository.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("删除菜单失败: %w", err)
	}

	return nil
}
