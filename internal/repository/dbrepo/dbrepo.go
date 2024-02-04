package dbrepo

import (

	"github.com/rashidalam9678/project-management-software-server/internal/config"
	"github.com/rashidalam9678/project-management-software-server/internal/repository"
	"gorm.io/gorm"
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