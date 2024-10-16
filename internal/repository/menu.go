package repository

import (
	"bubu_admin/internal/model"
	"context"
)

type MenuRepository interface {
	FindByID(ctx context.Context, id string) (*model.Menu, error)
	FindByCodeAndParentID(ctx context.Context, code, parentID string) (*model.Menu, error)
	Save(ctx context.Context, menu *model.Menu) (*model.Menu, error)
	List(ctx context.Context, menus *model.Menu) ([]*model.Menu, error)
	FindByIDS(ctx context.Context, ids []string) ([]*model.Menu, error)
	Update(ctx context.Context, menu *model.Menu) (*model.Menu, error)
	Delete(ctx context.Context, id string) error
}

func NewMenuRepository(
	repository *Repository,
) MenuRepository {
	return &menuRepository{
		Repository: repository,
	}
}

type menuRepository struct {
	*Repository
}

func (r *menuRepository) FindByID(ctx context.Context, id string) (*model.Menu, error) {
	var menu model.Menu
	if err := r.DB(ctx).Where("id = ?", id).Find(&menu).Error; err != nil {
		return nil, err
	}
	if menu.ID == "" {
		return nil, nil
	}
	return &menu, nil
}

func (r *menuRepository) FindByCodeAndParentID(ctx context.Context, code, parentID string) (*model.Menu, error) {
	var menu model.Menu
	if err := r.DB(ctx).Where("code = ? AND parent_id = ?", code, parentID).Find(&menu).Error; err != nil {
		return nil, err
	}
	if menu.ID == "" {
		return nil, nil
	}
	return &menu, nil
}

func (r *menuRepository) Save(ctx context.Context, menu *model.Menu) (*model.Menu, error) {
	if err := r.DB(ctx).Create(menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}

func (r *menuRepository) List(ctx context.Context, menus *model.Menu) ([]*model.Menu, error) {
	var menuList []*model.Menu
	if err := r.DB(ctx).Where(menus).Find(&menuList).Error; err != nil {
		return nil, err
	}
	return menuList, nil
}

func (r *menuRepository) FindByIDS(ctx context.Context, ids []string) ([]*model.Menu, error) {
	var menuList []*model.Menu
	if err := r.DB(ctx).Where("id IN (?)", ids).Find(&menuList).Error; err != nil {
		return nil, err
	}
	return menuList, nil
}

func (r *menuRepository) Update(ctx context.Context, menu *model.Menu) (*model.Menu, error) {
	if err := r.DB(ctx).Save(menu).Error; err != nil {
		return nil, err
	}
	return menu, nil
}


func (r *menuRepository) Delete(ctx context.Context, id string) error {
	return r.DB(ctx).Delete(&model.Menu{}, "id = ?", id).Error
}
