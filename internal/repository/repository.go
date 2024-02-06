package repository

import model "github.com/rashidalam9678/project-management-software-server/internal/models"

type Database interface{
	GetUserByEmail(email string) (model.User,error)
	InsertUser(email,externalId string)(string,error)
	GetUserByID(externalId string) (model.User, error)
	UpdateUserByID(email,id string) error
	DeleteUserByID(id string) error

	//project related
	InsertProject(title,description,userId string)(uint,error)
	GetProjectByID(id uint)(*model.Project,error)
	GetAllProjectsByUserID(userId string)([]model.Project,error)
	UpdateProjectByID(title,description string,id uint) error
	DeleteProjectByID(id uint) error

}