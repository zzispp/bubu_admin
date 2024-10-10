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
		Handler:      handler,
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
