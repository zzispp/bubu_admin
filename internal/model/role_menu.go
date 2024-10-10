package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type RoleMenu struct {
	ID        string    `json:"id" gorm:"size:20;primarykey"` // Unique ID
	RoleID    string    `json:"role_id" gorm:"size:20;index"` // From Role.ID
	MenuID    string    `json:"menu_id" gorm:"size:20;index"` // From Menu.ID
	CreatedAt time.Time `json:"created_at" gorm:"index;"`     // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`     // Update time
}

func (r *RoleMenu) TableName() string {
	return "role_menus"
}

func (r *RoleMenu) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = xid.New().String()
	}
	// 设置创建时间和更新时间
	now := time.Now()
	if r.CreatedAt.IsZero() {
		r.CreatedAt = now
	}
	r.UpdatedAt = now
	return nil
}
