package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"size:64;not null;uniqueIndex;comment:登录用户名" json:"username"`
	Password  string    `gorm:"size:256;not null;comment:密码" json:"-"`
	RealName  string    `gorm:"size:64;not null;comment:真实姓名" json:"real_name"`
	Role      string    `gorm:"size:32;not null;default:student;index;comment:角色" json:"role"`
	Email     string    `gorm:"size:128;default:null;comment:邮箱" json:"email,omitempty"`
	Phone     string    `gorm:"size:20;default:null;comment:手机号" json:"phone,omitempty"`
	Avatar    string    `gorm:"size:256;default:null;comment:头像URL" json:"avatar,omitempty"`
	Status    int       `gorm:"not null;default:1;index;comment:状态 1-正常 0-禁用" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
