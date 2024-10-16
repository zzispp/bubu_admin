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
	Update(ctx context.Context, role *model.Role) (*model.Role, error)
	List(ctx context.Context, role *model.Role) ([]*model.Role, error)
	Delete(ctx context.Context, id string) error
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

func (r *roleRepository) Update(ctx context.Context, role *model.Role) (*model.Role, error) {
	if err := r.DB(ctx).Model(role).Updates(role).Error; err != nil {
		return nil, err
	}
	return role, nil
}

func (r *roleRepository) List(ctx context.Context, query *model.Role) ([]*model.Role, error) {
	var roles []*model.Role
	db := r.DB(ctx)

	if query != nil {
		if query.Name != "" {
			db = db.Where("name LIKE ?", "%"+query.Name+"%")
		}
		if query.Code != "" {
			db = db.Where("code LIKE ?", "%"+query.Code+"%")
		}
		if query.Status != "" {
			db = db.Where("status = ?", query.Status)
		}
	}

	if err := db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *roleRepository) Delete(ctx context.Context, id string) error {
	return r.DB(ctx).Where("id = ?", id).Delete(&model.Role{}).Error
}
