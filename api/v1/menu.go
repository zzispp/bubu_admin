package v1

import "bubu_admin/internal/model"

type CreateMenuRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sequence    int32    `json:"sequence"`
	Type        string `json:"type"`
	Icon        string `json:"icon"`
	Path        string `json:"path"`
	Component   string `json:"component"`
	Redirect    string `json:"redirect,omitempty"`
	Status      string    `json:"status"`
	ParentID    string `json:"parent_id"`
}
type ListMenuRequest struct {
	Code string `json:"code"`	
	Name string `json:"name"`
	Status string `json:"status"`
}
type ListMenuResponse struct {
	Menus []*model.Menu `json:"menus"`
}

type UpdateMenuRequest struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sequence    int32  `json:"sequence"`
	Type        string `json:"type"`
	Icon        string `json:"icon"`
	Path        string `json:"path"`
	Component   string `json:"component"`
	Redirect    string `json:"redirect,omitempty"`
	Status      string `json:"status"`
	ParentID    string `json:"parent_id"`
}
