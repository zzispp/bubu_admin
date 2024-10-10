package repository

import (
	"bubu_admin/internal/model"
	"context"
)

type RoleMenuRepository interface {
	Save(ctx context.Context, roleMenu *model.RoleMenu) (*model.RoleMenu, error)
	FindByRoleID(ctx context.Context, roleID string) ([]*model.RoleMenu, error)
	FindByRoleIDS(ctx context.Context, roleIDS []string) ([]*model.RoleMenu, error)
}

func NewRoleMenuRepository(
	repository *Repository,
) RoleMenuRepository {
	return &roleMenuRepository{
		Repository: repository,
	}
}

type roleMenuRepository struct {
	*Repository
}

// Save implements RoleMenuRepository.
func (r *roleMenuRepository) Save(ctx context.Context, roleMenu *model.RoleMenu) (*model.RoleMenu, error) {
	if err := r.DB(ctx).Create(roleMenu).Error; err != nil {
		return nil, err
	}
	return roleMenu, nil
}

// FindByRoleID implements RoleMenuRepository.
func (r *roleMenuRepository) FindByRoleID(ctx context.Context, roleID string) ([]*model.RoleMenu, error) {
	var roleMenus []*model.RoleMenu
	if err := r.DB(ctx).Where("role_id = ?", roleID).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}

func (r *roleMenuRepository) FindByRoleIDS(ctx context.Context, roleIDS []string) ([]*model.RoleMenu, error) {
	var roleMenus []*model.RoleMenu
	if err := r.DB(ctx).Where("role_id IN (?)", roleIDS).Find(&roleMenus).Error; err != nil {
		return nil, err
	}
	return roleMenus, nil
}