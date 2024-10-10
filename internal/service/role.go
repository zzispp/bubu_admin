package service

import (
	v1 "bubu_admin/api/v1"
	"bubu_admin/internal/model"
	"bubu_admin/internal/repository"
	"bubu_admin/internal/utils"
	"context"
	"fmt"

	"github.com/rs/xid"
)

type RoleService interface {
	CreateRole(ctx context.Context, req *v1.CreateRoleRequest) error
}

func NewRoleService(
	service *Service,
	roleRepository repository.RoleRepository,
	menuRepository repository.MenuRepository,
	roleMenuRepository repository.RoleMenuRepository,
) RoleService {
	return &roleService{
		Service:        service,
		roleRepository: roleRepository,
		menuRepository: menuRepository,
		roleMenuRepository: roleMenuRepository,
	}
}

type roleService struct {
	*Service
	roleRepository repository.RoleRepository
	menuRepository repository.MenuRepository
	roleMenuRepository repository.RoleMenuRepository
}

// CreateRole 实现 RoleService 接口。
func (r *roleService) CreateRole(ctx context.Context, req *v1.CreateRoleRequest) error {
	// 检查角色是否已存在
	existingRole, err := r.roleRepository.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if existingRole != nil {
		return fmt.Errorf("角色已存在")
	}

	roleId := xid.New().String()
	role := &model.Role{
		ID:          roleId,
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Sequence:    req.Sequence,
		Status:      req.Status,
		MenuIDs:     req.Menus,
	}

	err = r.tm.Transaction(ctx, func(ctx context.Context) error {
		// 获取所有菜单
		allMenus, err := r.menuRepository.List(ctx, &model.Menu{})
		if err != nil {
			return err
		}
		//查询当前已经拥有的菜单
		existingMenus, err := r.roleMenuRepository.FindByRoleID(ctx, roleId)
		if err != nil {
			return err
		}

		// 获取需要添加的菜单
		menusToAdd := utils.GetMenuTreeBuilderInstance().CollectMenusAndParents(role.MenuIDs, existingMenus, allMenus)

		// 添加收集到的所有菜单
		for _, menuId := range menusToAdd {
			if _, err := r.roleMenuRepository.Save(ctx, &model.RoleMenu{
				RoleID: roleId,
				MenuID: menuId,
			}); err != nil {
				return err
			}
		}

		_, err = r.roleRepository.Save(ctx, role)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

