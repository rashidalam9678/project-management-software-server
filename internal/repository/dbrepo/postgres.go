package dbrepo

import (
	"context"
	"errors"

	"time"

	model "github.com/rashidalam9678/project-management-software-server/internal/models"
	"gorm.io/gorm"
)

//GetUserByEmail gets a user by email and return error if user not found or any other error
func (p *postgresDBRepo) GetUserByEmail(email string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user= model.User{}
	result:=p.DB.WithContext(ctx).Where("email = ?",email).First(&user)
	if result.Error == gorm.ErrRecordNotFound{
		return model.User{}, errors.New("record not found")
	}else if result.Error != nil{
		return model.User{},result.Error
	}

	return user, nil
}

//GetUserByExternalID gets a user by external id and return error if user not found or any other error
func (p *postgresDBRepo) GetUserByID(id string) (model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user= model.User{}
	result:=p.DB.WithContext(ctx).Where("id = ?",id).First(&user)
	if result.Error == gorm.ErrRecordNotFound{
		return model.User{}, errors.New("record not found")
	}else if result.Error != nil{
		return model.User{},result.Error
	}

	return user, nil
}


//InsertUser inserts a new user into the database and return the id of the user and error if any
func (p *postgresDBRepo) InsertUser(email,id string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	user:=model.User{
		Email: email,
		ID: id,
	}

	result:= p.DB.WithContext(ctx).Create(&user)
	if result.Error != nil{
		return "", errors.New(result.Error.Error())
	}

	return user.ID, nil
}

//UpdateUserByID updates a user by id and return error if any
func (p *postgresDBRepo) UpdateUserByID(email,id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	user:=model.User{
		Email: email,
		ID: id,
	}

	result:= p.DB.WithContext(ctx).Model(&user).Where("id = ?",id).Updates(&user)
	if result.Error != nil{
		return result.Error
	}

	return nil
}

//DeleteUserByID deletes a user by id and return error if any
func (p *postgresDBRepo) DeleteUserByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result:= p.DB.WithContext(ctx).Where("id = ?",id).Delete(&model.User{})
	if result.Error != nil{
		return result.Error
	}

	return nil
}

//project related

// CreateProject creates a new project and return the id of the project and error if any
func (p *postgresDBRepo) InsertProject(title,description,userID string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	project:=model.Project{
		Title: title,
		Description: description,
		UserID: userID,
	}

	result:= p.DB.WithContext(ctx).Create(&project)
	if result.Error != nil{
		return 0, result.Error
	}

	return project.ID, nil
}

//GetProjectByID gets a project by id and return error if project not found or any other error
func (p *postgresDBRepo) GetProjectByID(id uint) (*model.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var project  model.Project
	result:=p.DB.WithContext(ctx).Where("id = ?",id).First(&project)
	if result.Error == gorm.ErrRecordNotFound{
		return nil, errors.New("record not found")
	}else if result.Error != nil{
		return nil,result.Error
	}

	return &project, nil
}

//UpdateProjectByID updates a project by id and return error if any
func (p *postgresDBRepo) UpdateProjectByID(title,description string,id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	project:=model.Project{
		Title: title,
		Description: description,
	}

	result:= p.DB.WithContext(ctx).Model(&project).Where("id = ?",id).Updates(&project)
	if result.Error != nil{
		return result.Error
	}

	return nil
}

//DeleteProjectByID deletes a project by id and return error if any
func (p *postgresDBRepo) DeleteProjectByID(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result:= p.DB.WithContext(ctx).Where("id = ?",id).Delete(&model.Project{})
	if result.Error != nil{
		return result.Error
	}

	return nil
}

//GetProjectsByUserID gets all the projects by user id and return error if any
func (p *postgresDBRepo) GetAllProjectsByUserID(userID string) ([]model.Project, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var projects= []model.Project{}
	result:=p.DB.WithContext(ctx).Where("user_id = ?",userID).Find(&projects)
	if result.Error != nil{
		return nil,result.Error
	}

	return projects, nil
}
