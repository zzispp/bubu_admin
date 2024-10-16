package model

import (
	"bubu_admin/internal/consts"
	"fmt"
	"time"

	"github.com/rs/xid"
	"gorm.io/gorm"
)

type Menu struct {
	ID          string    `json:"id" gorm:"size:20;primarykey;"`   // 唯一标识
	Code        string    `json:"code" gorm:"size:32;index;"`      // 菜单代码（每个级别唯一）
	Name        string    `json:"name" gorm:"size:128;index"`      // 菜单显示名称
	Description string    `json:"description" gorm:"size:1024"`    // 菜单详细描述
	Sequence    int32     `json:"sequence" gorm:"index;"`          // 排序序列（升序排列）
	Type        string    `json:"type" gorm:"size:20;index"`       // 菜单类型（nav,page）
	Path        string    `json:"path" gorm:"size:255;"`           // 菜单访问路径
	Redirect    string    `json:"redirect,omitempty" gorm:"size:255;"` // 重定向路径
	Component   string    `json:"component,omitempty" gorm:"size:255;"` // 组件路径
	Icon        string    `json:"icon,omitempty" gorm:"size:255;"`      // 图标
	Status      string    `json:"status" gorm:"size:20;index"`     // 菜单状态（enable、disable）
	ParentID    string    `json:"parent_id" gorm:"size:20;index;"` // 父级ID（来自Menu.ID）
	Children    []*Menu   `json:"children,omitempty" gorm:"-"`     // 子菜单
	CreatedAt   time.Time `json:"created_at" gorm:"index;"`        // 创建时间
	UpdatedAt   time.Time `json:"updated_at" gorm:"index;"`        // 更新时间
}

func (m *Menu) TableName() string {
	return "menus"
}

// BeforeSave 在保存菜单记录之前执行的钩子函数
func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	// 如果ID为空，则生成一个新的ID
	if m.ID == "" {
		m.ID = xid.New().String()
	}

	// 设置创建时间和更新时间
	now := time.Now()
	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}
	m.UpdatedAt = now

	// 验证菜单类型
	if !consts.IsValidMenuType(m.Type) {
		return fmt.Errorf("无效的菜单类型: %s", m.Type)
	}


	// 验证状态
	if !consts.IsValidStatus(m.Status) {
		return fmt.Errorf("无效的状态: %s", m.Status)
	}

	return nil
}
