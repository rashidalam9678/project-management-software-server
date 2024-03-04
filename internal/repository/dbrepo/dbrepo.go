package dbrepo

import (

	"github.com/rashidalam9678/project-management-software-server/internal/config"
	"github.com/rashidalam9678/project-management-software-server/internal/repository"
	"gorm.io/gorm"
	model "github.com/rashidalam9678/project-management-software-server/internal/models"
)

type postgresDBRepo struct{
	App *config.AppConfig
	DB *gorm.DB
}

func NewPostgresRepo(conn *gorm.DB, a *config.AppConfig ) repository.Database{
	return &postgresDBRepo{
		App:a,
		DB:conn,
	}
}

// Insert Default Role and Permissions
func (p *postgresDBRepo) InsertDefaultRoleAndPermissions(projectID uint) error {
	// Insert Default Role
	role := model.Role{
		Name:      "guest",
		ProjectID: projectID,
	}
	result := p.DB.Create(&role)
	if result.Error != nil {
		return result.Error
	}
	// Insert Default Permissions
	err := p.InsertDefaultPermissions(role.ID)
	if err != nil {
		return err
	}
	return nil
}

// Insert Default Permissions
func (p *postgresDBRepo) InsertDefaultPermissions(roleID uint) error {
	permission := model.Permission{
		RoleID: roleID,
	}
	result := p.DB.Create(&permission)
	if result.Error != nil {
		return result.Error
	}
	return nil
}