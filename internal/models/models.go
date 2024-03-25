package model

import (
	"time"
)

type User struct {
	ID        string    `gorm:"primaryKey; index:idx_user_id"`
	Email     string    `gorm:"unique;not null"`
	FirstName string    `gorm:"not null"`
	LastName  string    `gorm:"not null"`
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
	Invites     []Invite
	Members     []Membership
	Tasks       []Task
	Roles       []Role
	Messages    []Message
}

type Invite struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Email       string `gorm:"uniqueIndex:idx_email_project_id"`
	ProjectID   uint   `gorm:"uniqueIndex:idx_email_project_id"`
	Description string
	InvitedBy   string    `gorm:"not null"`
	Status      string    `gorm:"not null;default:'pending'"`
	Token       string    `gorm:"not null"`
	ExpiresAt   time.Time `gorm:"not null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}

type Membership struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	ProjectID uint      `gorm:"not null;uniqueInndex:idx_member_project_id"`
	UserID    string    `gorm:"not null;uniqueIndex:idx_member_project_id"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	RoleId    uint
}

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"not null"`
	ProjectID   uint      `gorm:"not null;index:idx_role_project_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	Permissions Permission
}

type Permission struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	RoleID    uint      `gorm:"not null;index:idx_permission_role_id"`
	TaskR     bool      `gorm:"not null; default:true"`
	TaskW     bool      `gorm:"not null; default:false"`
	TaskD     bool      `gorm:"not null; default:false"`
	MemberR   bool      `gorm:"not null; default:false"`
	MemberW   bool      `gorm:"not null; default:false"`
	MemberD   bool      `gorm:"not null; default:false"`
	CreateAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Task struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	ProjectID   uint   `gorm:"not null;index:idx_task_project_id"`
	CreatedBy   string `gorm:"not null"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	AssignedTo  string
	Status      string    `gorm:"not null;default:'todo'"`
	Priority    string    `gorm:"not null;default:'low'"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type Message struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProjectID uint      `json:"project_id"`
	SenderID  uint      `json:"sender_id"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
