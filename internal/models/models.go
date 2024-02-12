package model

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey; index:idx_user_id"`
	Email     string    `gorm:"unique;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	Projects  []Project `gorm:"foreignKey:UserID;references:ID"` // One-to-Many relationship: One User can have multiple Projects
}

type Project struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Title       string    `gorm:"not null"`
	Description string    `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	UserID      string    `gorm:"not null;index:idx_project_user_id"`
	Invites     []Invite  `gorm:"foreignKey:ProjectID;references:ID"`
	Memberships []Membership `gorm:"foreignKey:ProjectID;references:ID"`
	Tasks       []Task `gorm:"foreignKey:ProjectID;references:ID"`
}

type Invite struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"not null"`
	ProjectID uint      `gorm:"unique;not null;index:idx_invite_project_id"`
	InvitedBy string    `gorm:"not null"`
	Status    string    `gorm:"not null;default:'pending'"`
	Token     string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type Membership struct {
	ID uint `gorm:"primaryKey;autoIncrement"`
	ProjectID uint `gorm:"unique;not null;index:idx_member_project_id"`
	UserID    string `gorm:"not null;index:idx_member_user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	ProjectRoles ProjectRole 
}

type ProjectRole struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	ProjectID uint `gorm:"not null;index:idx_project_role_project_id"`
	MembershipID uint `gorm:"not null;index:idx_project_role_member_id"`
	RoleID uint `gorm:"not null;index:idx_project_role_role_id"`
	Permission Permission `gorm:"foreignKey:ProjectRoleID;references:ID"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Role struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"not null"`
	ProjectID uint `gorm:"not null;index:idx_role_project_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}


type Permission struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	ProjectRoleID uint `gorm:"not null;index:idx_permission_role_id"`
	TaskR bool `gorm:"not null; default:true"`
	TaskW bool `gorm:"not null; default:false"`
	TaskD bool `gorm:"not null; default:false"`
	CreateAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

}

type Task struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	ProjectID uint `gorm:"unique;not null;index:idx_task_project_id"`
	CreatedBy string `gorm:"not null"`
	Title string `gorm:"not null"`
	AssignedTo string 
	Status     string `gorm:"not null;default:'open'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
