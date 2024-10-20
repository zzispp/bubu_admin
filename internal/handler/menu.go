package handler

import (
	v1 "bubu_admin/api/v1"
	"bubu_admin/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MenuHandler struct {
	*Handler
	menuService service.MenuService
}

func NewMenuHandler(
    handler *Handler,
    menuService service.MenuService,
) *MenuHandler {
	return &MenuHandler{
		Handler:      handler,
		menuService: menuService,
	}
}

func (h *MenuHandler) CreateMenu(ctx *gin.Context) {
	req := new(v1.CreateMenuRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrBadRequest, nil)
		return
	}

	if err := h.menuService.CreateMenu(ctx, req); err != nil {
		h.logger.WithContext(ctx).Error("menuService.Create error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}


func (h *MenuHandler) ListMenu(ctx *gin.Context) {
	req := new(v1.ListMenuRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrBadRequest, nil)
		return
	}

	menus, err := h.menuService.ListMenu(ctx, req)
	if  err != nil {
		h.logger.WithContext(ctx).Error("menuService.ListMenu error", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, &v1.ListMenuResponse{Menus: menus})
}

func (h *MenuHandler) UpdateMenu(ctx *gin.Context) {
	id := ctx.Param("id")
	req := new(v1.UpdateMenuRequest)
	if err := ctx.ShouldBindJSON(req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	if err := h.menuService.UpdateMenu(ctx, id, req); err != nil {
		h.logger.WithContext(ctx).Error("menuService.UpdateMenu 错误", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}

func (h *MenuHandler) GetMenuByID(ctx *gin.Context) {
	id := ctx.Param("id")

	menu, err := h.menuService.GetMenuByID(ctx, id)
	if err != nil {
		h.logger.WithContext(ctx).Error("menuService.GetMenuByID 错误", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, menu)
}

func (h *MenuHandler) DeleteMenu(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := h.menuService.DeleteMenu(ctx, id); err != nil {
		h.logger.WithContext(ctx).Error("menuService.DeleteMenu 错误", zap.Error(err))
		v1.HandleError(ctx, http.StatusInternalServerError, err, nil)
		return
	}

	v1.HandleSuccess(ctx, nil)
}


