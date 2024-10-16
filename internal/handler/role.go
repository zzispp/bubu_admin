package handler

import (
	v1 "bubu_admin/api/v1"
	"bubu_admin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RoleHandler struct {
	*Handler
	roleService service.RoleService
}

func NewRoleHandler(
	handler *Handler,
	roleService service.RoleService,
) *RoleHandler {
	return &RoleHandler{
		Handler:     handler,
		roleService: roleService,
	}
}

func (h *RoleHandler) CreateRole(ctx *gin.Context) {
	req := new(v1.CreateRoleRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrBadRequest, nil)
		return
	}

	if err := h.roleService.CreateRole(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("roleService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *RoleHandler) UpdateRole(ctx *gin.Context) {
	id := ctx.Param("id")
	req := new(v1.UpdateRoleRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrBadRequest, nil)
		return
	}

	if err := h.roleService.UpdateRole(ctx, id,req); err != nil {
		h.logger.WithContext(ctx).Error("roleService.Update error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *RoleHandler) ListRole(ctx *gin.Context) {
	req := new(v1.ListRoleRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrBadRequest, nil)
		return
	}

	roles, err := h.roleService.ListRole(ctx, req)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.List error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, &v1.ListRoleResponse{
		Roles: roles,
	})
}

func (h *RoleHandler) GetRoleByID(ctx *gin.Context) {
	id := ctx.Param("id")

	role, err := h.roleService.GetRoleByID(ctx, id)
	if err != nil {
		h.logger.WithContext(ctx).Error("roleService.GetRoleByID error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}



	v1.HandleSuccess(ctx, role)
}

func (h *RoleHandler) DeleteRole(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.roleService.DeleteRole(ctx, id); err != nil {
		h.logger.WithContext(ctx).Error("roleService.DeleteRole error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}