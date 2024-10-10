package v1

import "bubu_admin/internal/model"

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
	Password string `json:"password" binding:"required" example:"123456"`
}
type LoginResponseData struct {
	AccessToken string `json:"accessToken"`
}
type LoginResponse struct {
	Response
	Data LoginResponseData
}

type UpdateProfileRequest struct {
	Name string `json:"name" example:"alan"`
	Email    string `json:"email" binding:"required,email" example:"1234@gmail.com"`
}
type GetProfileResponseData struct {
	UserId   string `json:"user_id"`
	Name string `json:"name" example:"alan"`
	Email string `json:"email" example:"1234@gmail.com"`
	Roles []*model.Role `json:"roles"`
	Menus []*model.Menu `json:"menus"`
}
type GetProfileResponse struct {
	Response
	Data GetProfileResponseData
}

type AddRoleToUserRequest struct {
	UserId string `json:"user_id" binding:"required"`
	RoleIds []string `json:"role_ids" binding:"required"`
}
