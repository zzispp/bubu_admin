package model

import (
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type UserRole struct {
	ID        string    `json:"id" gorm:"size:50;primarykey"`           // Unique ID
	UserID    string    `json:"user_id" gorm:"size:50;index"`           // From User.ID
	RoleID    string    `json:"role_id" gorm:"size:50;index"`           // From Role.ID
	CreatedAt time.Time `json:"created_at" gorm:"index;"`               // Create time
	UpdatedAt time.Time `json:"updated_at" gorm:"index;"`               // Update time
	RoleName  string    `json:"role_name" gorm:"<-:false;-:migration;"` // From Role.Name
}

func (a *UserRole) TableName() string {
	return "user_role"
}

func (a *UserRole) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == "" {
		a.ID = xid.New().String()
	}
	now := time.Now()
	if a.CreatedAt.IsZero() {
		a.CreatedAt = now
	}
	a.UpdatedAt = now
	return
}