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
}

func NewUserService(
	service *Service,
	userRepo repository.UserRepository,
	menuRepo repository.MenuRepository,
	roleRepo repository.RoleRepository,
	userRoleRepo repository.UserRoleRepository,
	roleMenuRepo repository.RoleMenuRepository,
) UserService {
	return &userService{
		userRepo:     userRepo,
		menuRepo:     menuRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		roleMenuRepo: roleMenuRepo,
		Service:      service,
	}
}

type userService struct {
	userRepo     repository.UserRepository
	menuRepo     repository.MenuRepository
	roleRepo     repository.RoleRepository
	userRoleRepo repository.UserRoleRepository
	roleMenuRepo repository.RoleMenuRepository
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
		menus, err := s.menuRepo.List(ctx,&model.Menu{})
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

	// 获取用户角色
	userRoles, err := s.userRoleRepo.FindByUserID(ctx, userId)
	if err != nil {
		return nil, err
	}

	roleIds := make([]string, 0)
	for _, role := range userRoles {
		roleIds = append(roleIds, role.RoleID)
	}

	//根据用户的角色拿到对应的菜单
	roleMenus, err := s.roleMenuRepo.FindByRoleIDS(ctx, roleIds)
	if err != nil {
		return nil, err
	}

	// 使用map来存储菜单ID，确保菜单ID的唯一性
	menuIdMap := make(map[string]struct{})
	for _, roleMenu := range roleMenus {
		menuIdMap[roleMenu.MenuID] = struct{}{}
	}

	// 将map中的键转换为切片，得到去重后的菜单ID列表
	menuIds := make([]string, 0, len(menuIdMap))
	for menuId := range menuIdMap {
		menuIds = append(menuIds, menuId)
	}

	// 获取去重后的菜单列表
	menus, err := s.menuRepo.FindByIDS(ctx, menuIds)
	if err != nil {
		return nil, err
	}

	// 构建菜单树
	menuTree := utils.GetMenuTreeBuilderInstance().BuildMenuTree(menus)

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

	roles, err := s.roleRepo.FindByIDS(ctx, roleIds)
	if err != nil {
		return err
	}

	if len(roleIds) != len(roles) {
		return v1.ErrRoleNotFound
	}

	// 获取用户已有的角色
	existingUserRoles, err := s.userRoleRepo.FindByUserID(ctx, userId)
	if err != nil {
		return err
	}

	// 创建一个map来存储已有的角色ID
	existingRoleIDs := make(map[string]bool)
	for _, ur := range existingUserRoles {
		existingRoleIDs[ur.RoleID] = true
	}

	// 只添加新的角色
	userRoles := make([]*model.UserRole, 0)
	for _, role := range roles {
		if !existingRoleIDs[role.ID] {
			userRoles = append(userRoles, &model.UserRole{
				UserID: userId,
				RoleID: role.ID,
			})
		}
	}

	// 如果有新角色需要添加，则保存
	if len(userRoles) > 0 {
		err = s.userRoleRepo.SaveBatch(ctx, userRoles)
		if err != nil {
			return err
		}
	}

	return nil
}
