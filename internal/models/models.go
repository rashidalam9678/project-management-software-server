package model

import (

	"gorm.io/gorm"
)

type User struct {
    gorm.Model // This includes fields like ID, CreatedAt, UpdatedAt, DeletedAt
    Email            string `gorm:"unique;not null"`
    ExternalID      string `gorm:"not null;unique;index:idx_user_external_id"` // Unique identifier for the user in the external system
    Projects []Project `gorm:"foreignKey:UserID;references:ExternalID"`// One-to-Many relationship: One User can have multiple Projects
}

type Project struct {
    gorm.Model 

    Title       string `gorm:"not null"`
    Description string  `gorm:"not null"`
    UserID      string    `gorm:"not null;index:idx_project_user_id"`// Foreign key referencing the ID of the User
}

// type UserVerification struct {
//     gorm.Model 

//     UserID           uint `gorm:"unique;not null"` // Foreign key referencing the ID of the User
//     TokenType        string `gorm:"not null"`      // Type of the token (e.g., "email_verification", "invitation")
//     Token            string `gorm:"not null"`
//     ExpiresAt        time.Time
//     IsTokenConsumed  bool   // Flag indicating whether the token has been used or consumed
//     InvitationReason string // Additional information for invitation tokens
// }