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
	UserID      string    `gorm:"not null;index:idx_project_user_id"` // Foreign key referencing the ID of the User
}
