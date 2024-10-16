package service

import (
	v1 "bubu_admin/api/v1"
	"bubu_admin/internal/consts"
	"bubu_admin/internal/model"
	"bubu_admin/internal/repository"
	"bubu_admin/internal/utils"
	"context"
	"fmt"

	"github.com/rs/xid"
)

type RoleService interface {
	CreateRole(ctx context.Context, req *v1.CreateRoleRequest) error
	UpdateRole(ctx context.Context, id string, req *v1.UpdateRoleRequest) error
	ListRole(ctx context.Context, req *v1.ListRoleRequest) ([]*model.Role, error)
	GetRoleByID(ctx context.Context, id string) (*model.Role, error)
	DeleteRole(ctx context.Context, id string) error
}

func NewRoleService(
	service *Service,
	roleRepository repository.RoleRepository,
	menuRepository repository.MenuRepository,
	userRepository repository.UserRepository,
) RoleService {
	return &roleService{
		Service:        service,
		roleRepository: roleRepository,
		menuRepository: menuRepository,
		userRepository: userRepository,
	}
}

type roleService struct {
	*Service
	roleRepository repository.RoleRepository
	menuRepository repository.MenuRepository
	userRepository repository.UserRepository
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
	}

	err = r.tm.Transaction(ctx, func(ctx context.Context) error {
		// 保存角色
		_, err = r.roleRepository.Save(ctx, role)
		if err != nil {
			return err
		}

		if len(req.Menus) > 0 {
			// 获取所有菜单
			allMenus, err := r.menuRepository.List(ctx, &model.Menu{})
			if err != nil {
				return fmt.Errorf("获取菜单列表失败: %w", err)
			}

			// 收集需要添加的菜单ID，包括父级菜单
			menusToAdd := utils.GetMenuTreeBuilderInstance().CollectMenusAndParents(req.Menus, allMenus)

			// 为角色添加菜单权限
			policies := make([][]string, len(menusToAdd))
			for i, menuID := range menusToAdd {
				policies[i] = []string{roleId, menuID, consts.StrategyTypeMenu}
			}

			_, err = r.casbin.AddPolicies(policies)
			if err != nil {
				return fmt.Errorf("添加角色菜单权限失败: %w", err)
			}
		}

		return nil
	})
	return err
}

func (r *roleService) UpdateRole(ctx context.Context, id string, req *v1.UpdateRoleRequest) error {
	// 根据ID查询角色
	existingRole, err := r.roleRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingRole == nil {
		return fmt.Errorf("角色不存在")
	}

	// 检查是否存在相同Code但ID不同的角色
	roleWithSameCode, err := r.roleRepository.FindByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if roleWithSameCode != nil && roleWithSameCode.ID != existingRole.ID {
		return fmt.Errorf("已存在相同代码的角色")
	}

	existingRole.Code = req.Code
	existingRole.Name = req.Name
	existingRole.Description = req.Description
	existingRole.Sequence = req.Sequence
	existingRole.Status = req.Status

	err = r.tm.Transaction(ctx, func(ctx context.Context) error {
		// 修改角色
		_, err = r.roleRepository.Update(ctx, existingRole)
		if err != nil {
			return err
		}
		if len(req.Menus) > 0 {
			// 获取所有菜单
			allMenus, err := r.menuRepository.List(ctx, &model.Menu{})
			if err != nil {
				return fmt.Errorf("获取菜单列表失败: %w", err)
			}

			// 收集需要添加的菜单ID，包括父级菜单
			menusToUpdate := utils.GetMenuTreeBuilderInstance().CollectMenusAndParents(req.Menus, allMenus)

			// 为角色添加菜单权限
			policies := make([][]string, len(menusToUpdate))
			for i, menuID := range menusToUpdate {
				policies[i] = []string{existingRole.ID, menuID, consts.StrategyTypeMenu}
			}

			//获取原有菜单
			oldMenus, err := r.casbin.GetMenusForRole(existingRole.ID)
			if err != nil {
				return fmt.Errorf("获取原有菜单失败: %w", err)
			}

			_, err = r.casbin.UpdatePolicies(oldMenus, policies)
			if err != nil {
				return fmt.Errorf("添加角色菜单权限失败: %w", err)
			}
		}

		return nil
	})
	return err
}

func (r *roleService) ListRole(ctx context.Context, req *v1.ListRoleRequest) ([]*model.Role, error) {
	return r.roleRepository.List(ctx, &model.Role{
		Code:   req.Code,
		Name:   req.Name,
		Status: req.Status,
	})
}

func (r *roleService) GetRoleByID(ctx context.Context, id string) (*model.Role, error) {
	role, err := r.roleRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("角色不存在")
	}

	//查询角色的菜单
	menus, err := r.casbin.GetMenusForRole(id)
	if err != nil {
		return nil, fmt.Errorf("获取角色菜单失败: %w", err)
	}

	if len(menus) > 0 {
		//提取菜单ID
		menuIds := utils.GetMenuTreeBuilderInstance().ExtractMenuIds(menus)

		menusEntities, err := r.menuRepository.FindByIDS(ctx, menuIds)
		if err != nil {
			return nil, err
		}

		// 构建菜单树
		menuTree := utils.GetMenuTreeBuilderInstance().BuildMenuTree(menusEntities)
		role.Menus = menuTree
	}
	return role, nil
}

func (r *roleService) DeleteRole(ctx context.Context, id string) error {
	// 首先检查角色是否存在
	existingRole, err := r.roleRepository.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取角色失败: %w", err)
	}
	if existingRole == nil {
		return fmt.Errorf("角色不存在")
	}

	//获取所有用户
	users, err := r.userRepository.List(ctx, &model.User{})
	if err != nil {
		return fmt.Errorf("获取用户失败: %w", err)
	}
	if len(users) > 0 {
		// 检查用户是否拥有该角色
		for _, user := range users {
			roles, err := r.casbin.GetRolesForUser(user.ID)
			if err != nil {
				return fmt.Errorf("获取用户角色失败: %w", err)
			}
			for _, role := range roles {
				if role == id {
					return fmt.Errorf("无法删除包含用户的角色")
				}
			}
		}
	}

	//查询角色的菜单
	menus, err := r.casbin.GetMenusForRole(id)
	if err != nil {
		return fmt.Errorf("获取角色菜单失败: %w", err)
	}
	if len(menus) > 0 {
		return fmt.Errorf("无法删除包含菜单的角色")
	}

	// 删除角色
	return r.roleRepository.Delete(ctx, id)
}
