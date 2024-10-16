package model

import (
	"bubu_admin/internal/consts"
	"fmt"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"` // Unique ID
	Code        string    `json:"code" gorm:"size:32;index;"`    // Code of role (unique)
	Name        string    `json:"name" gorm:"size:128;index"`    // Display name of role
	Description string    `json:"description" gorm:"size:1024"`  // Details about role
	Sequence    int32       `json:"sequence" gorm:"index"`         // Sequence for sorting
	Status      string    `json:"status" gorm:"size:20;index"`   // Status of role (disabled, enabled)
	CreatedAt   time.Time `json:"created_at" gorm:"index;"`      // Create time
	UpdatedAt   time.Time `json:"updated_at" gorm:"index;"`      // Update time
	Menus       []*Menu   `json:"menus" gorm:"many2many:role_menus;"` // 角色菜单
}

func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	// 如果ID为空，则生成一个新的ID
	if r.ID == "" {
		r.ID = xid.New().String()
	}

	// 设置创建时间和更新时间
	now := time.Now()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	r.UpdatedAt = now

	if !consts.IsValidStatus(r.Status) {
		return fmt.Errorf("无效的状态: %s", r.Status)
	}
	return nil
}
