package service

import (
	v1 "bubu_admin/api/v1"
	"bubu_admin/internal/model"
	"bubu_admin/internal/repository"
	"bubu_admin/internal/utils"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
	AddRoleToUser(ctx context.Context, userId string, roleIds []string) error
	List(ctx context.Context, req *v1.ListUserRequest) ([]*model.User, error)
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, id string, req *v1.UpdateUserRequest) error
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
	menuRepo repository.MenuRepository,
	roleRepo repository.RoleRepository,
) UserService {
	return &userService{
		userRepo: userRepo,
		menuRepo: menuRepo,
		roleRepo: roleRepo,
		Service:  service,
	}
}

type userService struct {
	userRepo repository.UserRepository
	menuRepo repository.MenuRepository
	roleRepo repository.RoleRepository
	*Service
}

func (s *userService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	// check username
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return v1.ErrInternalServerError
	}
	if err == nil && user != nil {
		return v1.ErrEmailAlreadyUse
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user = &model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	// Transaction demo
	err = s.tm.Transaction(ctx, func(ctx context.Context) error {
		// Create a user
		if err = s.userRepo.Create(ctx, user); err != nil {
			return err
		}
		// TODO: other repo
		return nil
	})
	return err
}

func (s *userService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {

	// login by root
	if req.Email == s.conf.GetString("security.root.email") {

		err := bcrypt.CompareHashAndPassword([]byte(s.conf.GetString("security.root.password")), []byte(req.Password))
		if err != nil {
			return "", v1.ErrEmailOrPasswordIncorrect
		}

		userID := s.conf.GetString("security.root.id")
		s.logger.WithContext(ctx).Info("Login by root")
		return s.jwt.GenToken(userID, time.Now().Add(time.Hour*24*90))
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil || user == nil {
		return "", v1.ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", v1.ErrEmailOrPasswordIncorrect
	}
	token, err := s.jwt.GenToken(user.ID, time.Now().Add(time.Hour*24*90))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {
	if userId == s.conf.GetString("security.root.id") {
		menus, err := s.menuRepo.List(ctx, &model.Menu{})
		if err != nil {
			return nil, fmt.Errorf("获取菜单列表失败")
		}

		// 将菜单列表转换为树结构
		menuTree := utils.GetMenuTreeBuilderInstance().BuildMenuTree(menus)

		return &v1.GetProfileResponseData{
			UserId: userId,
			Name:   s.conf.GetString("security.root.name"),
			Email:  s.conf.GetString("security.root.email"),
			Menus:  menuTree,
		}, nil
	}
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return nil, err
	}

	userRoles, err := s.casbin.GetRolesForUser(userId)
	if err != nil {
		return nil, err
	}
	menus := make([][]string, 0)
	//根据用户的角色拿到对应的菜单
	for _, role := range userRoles {
		roleMenus, err := s.casbin.GetMenusForRole(role)
		if err != nil {
			return nil, err
		}
		menus = append(menus, roleMenus...)
	}

	//提取菜单ID
	menuIds := utils.GetMenuTreeBuilderInstance().ExtractMenuIds(menus)

	menusEntities, err := s.menuRepo.FindByIDS(ctx, menuIds)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	menuTree := utils.GetMenuTreeBuilderInstance().BuildMenuTree(menusEntities)

	return &v1.GetProfileResponseData{
		UserId: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Menus:  menuTree,
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}

	user.Email = req.Email
	user.Name = req.Name

	if err = s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// 为用户添加角色
func (s *userService) AddRoleToUser(ctx context.Context, userId string, roleIds []string) error {
	user, err := s.userRepo.GetByID(ctx, userId)
	if err != nil {
		return err
	}
	if user == nil {
		return v1.ErrUserNotFound
	}
	//去重
	roleIds = utils.GetRoleBuilderInstance().RemoveDuplicateRoleIds(roleIds)

	roles, err := s.roleRepo.FindByIDS(ctx, roleIds)
	if err != nil {
		return err
	}
	if len(roles) == 0 {
		return v1.ErrRoleNotFound
	}

	// 为用户批量添加角色
	_, err = s.casbin.AddRolesForUser(userId, roleIds)
	if err != nil {
		return fmt.Errorf("为用户添加角色失败: %w", err)
	}

	return nil
}

func (s *userService) List(ctx context.Context, req *v1.ListUserRequest) ([]*model.User, error) {
	users, err := s.userRepo.List(ctx, &model.User{})
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		roles, err := s.casbin.GetRolesForUser(user.ID)
		if err != nil {
			return nil, err
		}
		roleEntities, err := s.roleRepo.FindByIDS(ctx, roles)
		if err != nil {
			return nil, err
		}
		user.Roles = roleEntities
	}
	return users, nil
}

func (s *userService) Delete(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}

func (s *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	roles, err := s.casbin.GetRolesForUser(user.ID)
	if err != nil {
		return nil, err
	}
	roleEntities, err := s.roleRepo.FindByIDS(ctx, roles)
	if err != nil {
		return nil, err
	}
	user.Roles = roleEntities
	return user, nil
}

func (s *userService) Update(ctx context.Context, id string, req *v1.UpdateUserRequest) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.Email = req.Email
	user.Name = req.Name

	_, err = s.roleRepo.FindByIDS(ctx, req.Roles)
	if err != nil {
		return fmt.Errorf("角色不存在")
	}

	// 更新角色
	// 先删除用户的所有角色
	oldRoles, err := s.casbin.GetRolesForUser(user.ID)
	if err != nil {
		return fmt.Errorf("获取用户角色失败: %w", err)
	}
	for _, role := range oldRoles {
		_, err = s.casbin.RemoveRoleForUser(user.ID, role)
		if err != nil {
			return fmt.Errorf("删除用户角色失败: %w", err)
		}
	}

	// 添加新的角色
	_, err = s.casbin.AddRolesForUser(user.ID, req.Roles)
	if err != nil {
		return fmt.Errorf("添加用户角色失败: %w", err)
	}

	return s.userRepo.Update(ctx, user)
}
