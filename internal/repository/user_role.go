package repository

import (
	"bubu_admin/internal/model"
	"context"

)

type UserRoleRepository interface {
	Save(ctx context.Context, userRole *model.UserRole) error
	SaveBatch(ctx context.Context, userRoles []*model.UserRole) error
	FindByUserID(ctx context.Context, userID string) ([]*model.UserRole, error)
}

func NewUserRoleRepository(
	repository *Repository,
) UserRoleRepository {
	return &userRoleRepository{
		Repository: repository,
	}
}

type userRoleRepository struct {
	*Repository
}

func (r *userRoleRepository) Save(ctx context.Context, userRole *model.UserRole) error {
	return r.DB(ctx).Create(userRole).Error
}

func (r *userRoleRepository) SaveBatch(ctx context.Context, userRoles []*model.UserRole) error {
	return r.DB(ctx).Create(&userRoles).Error
}

func (r *userRoleRepository) FindByUserID(ctx context.Context, userID string) ([]*model.UserRole, error) {
	var userRoles []*model.UserRole
	err := r.DB(ctx).Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}
