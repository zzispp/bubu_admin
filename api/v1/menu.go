package v1

import "bubu_admin/internal/model"

type CreateMenuRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Sequence    int32    `json:"sequence"`
	Type        string `json:"type"`
	Path        string `json:"path"`
	PathType    string `json:"path_type"`
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