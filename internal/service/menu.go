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
}

func NewMenuService(
	service *Service,
	menuRepository repository.MenuRepository,
) MenuService {
	return &menuService{
		Service:        service,
		menuRepository: menuRepository,
	}
}

type menuService struct {
	*Service
	menuRepository repository.MenuRepository
}

func (s *menuService) CreateMenu(ctx context.Context, req *v1.CreateMenuRequest) error {

	// 检查菜单类型和父级ID的合法性
	if req.Type == consts.MenuTypeNav && req.ParentID != "root" {
		return fmt.Errorf("导航菜单只能添加到最顶层")
	}
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
		Path:        req.Path,
		PathType:    req.PathType,
		Status:      req.Status,
		ParentID:    req.ParentID,
	}); err != nil {
		return fmt.Errorf("创建菜单失败")
	}

	return nil
}

func (s *menuService) ListMenu(ctx context.Context, req *v1.ListMenuRequest) ([]*model.Menu, error) {
	menus, err := s.menuRepository.List(ctx, &model.Menu{
		Code: req.Code,
		Name: req.Name,
		Status: req.Status,
	})
	if err != nil {
		return nil, fmt.Errorf("获取菜单列表失败")
	}

	// 将菜单列表转换为树结构
	menuTree := utils.GetMenuTreeBuilderInstance().BuildMenuTree(menus)
	return menuTree, nil
}