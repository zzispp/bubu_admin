package repository

import (
	"bubu_admin/internal/model"
	"context"
)

type RoleRepository interface {
	FindByCode(ctx context.Context, code string) (*model.Role, error)
	Save(ctx context.Context, role *model.Role) (*model.Role, error)
	FindByIDS(ctx context.Context, ids []string) ([]*model.Role, error)
	FindByID(ctx context.Context, id string) (*model.Role, error)
	FindByUserID(ctx context.Context, userId string) ([]*model.Role, error)
}

func NewRoleRepository(
	repository *Repository,
) RoleRepository {
	return &roleRepository{
		Repository: repository,
	}
}

type roleRepository struct {
	*Repository
}

func (r *roleRepository) FindByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	if err := r.DB(ctx).Where("code = ?", code).Find(&role).Error; err != nil {
		return nil, err
	}
	if role.ID == "" {
		return nil, nil
	}
	return &role, nil
}
func (r *roleRepository) Save(ctx context.Context, role *model.Role) (*model.Role, error) {
	if err := r.DB(ctx).Create(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) FindByIDS(ctx context.Context, ids []string) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.DB(ctx).Where("id IN (?)", ids).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) FindByID(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	if err := r.DB(ctx).Where("id = ?", id).Find(&role).Error; err != nil {
		return nil, err
	}
	if role.ID == "" {
		return nil, nil
	}
	return &role, nil
}

func (r *roleRepository) FindByUserID(ctx context.Context, userId string) ([]*model.Role, error) {
	var roles []*model.Role
	if err := r.DB(ctx).Where("user_id = ?", userId).Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}
