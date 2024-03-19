package dbrepo

import (
	"context"
	"errors"

	"time"
	model "github.com/rashidalam9678/project-management-software-server/internal/models"
)

//GetUserByEmail gets a user by email and return error if user not found or any other error
func (p *postgresDBRepo) GetUserByEmail(email string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user= model.User{}
	result:=p.DB.WithContext(ctx).Where("email = ?",email).First(&user)
	if result.Error != nil{
		return nil,result.Error
	
	}

	return &user, nil
}

//GetUserByExternalID gets a user by external id and return error if user not found or any other error
func (p *postgresDBRepo) GetUserByID(id string) (*model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user= model.User{}
	result:=p.DB.WithContext(ctx).Where("id = ?",id).First(&user)
	if result.Error != nil{
		return nil,result.Error
	
	}
	return &user, nil
}


//InsertUser inserts a new user into the database and return the id of the user and error if any
func (p *postgresDBRepo) InsertUser(email,id, firstName, lastName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	user:=model.User{
		Email: email,
		ID: id,
		FirstName: firstName,
		LastName: lastName,
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
	if result.Error != nil{
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

// Create Role
func (p *postgresDBRepo) InsertRole(projectID uint,roleName string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	role:=model.Role{
		ProjectID: projectID,
		Name: roleName,
	}

	result:= p.DB.WithContext(ctx).Create(&role)
	if result.Error != nil{
		return 0, result.Error
	}

	return role.ID, nil
}

//InsertInvite inserts a new invite into the database and return the id of the invite and error if any
func (p *postgresDBRepo) InsertInvite(email string,projectID uint,invitedBy string, description string, token string) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	
	invite:=model.Invite{
		Email: email,
		ProjectID: projectID,
		InvitedBy: invitedBy,
		Token: token,
		Description: description,
		ExpiresAt: time.Now().Add(time.Hour*24),
	}

	result:= p.DB.WithContext(ctx).Create(&invite)
	if result.Error != nil{
		return 0, result.Error
	}

	return invite.ID, nil
}

//GetInviteByID gets a invite by id and return error if invite not found or any other error
func (p *postgresDBRepo) GetInviteByID(id uint) (*model.Invite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var invite  model.Invite
	result:=p.DB.WithContext(ctx).Where("id = ?",id).First(&invite)
	if result.Error != nil{
		return nil,result.Error
	}
	return &invite, nil
}

//GetInvitesByProjectID gets all the invites by project id and return error if any
func (p *postgresDBRepo) GetInvitesByProjectID(projectID uint) ([]model.Invite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var invites= []model.Invite{}
	result:=p.DB.WithContext(ctx).Where("project_id = ?",projectID).Find(&invites)
	if result.Error != nil{
		return nil,result.Error
	}

	return invites, nil
}

//GetInviteByEmailAndProjectID gets a invite by email and project id and return error if invite not found or any other error
func (p *postgresDBRepo) GetInviteByEmailAndProjectID(email string,projectID uint) (*model.Invite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var invite  model.Invite
	result:=p.DB.WithContext(ctx).Where("email = ? AND project_id = ?",email,projectID).First(&invite)
	if result.Error != nil{
		return nil,result.Error
	}
	return &invite, nil
}

// GetInvite by token 
func (p *postgresDBRepo) GetInviteByToken(token string) (*model.Invite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var invite  model.Invite
	result:=p.DB.WithContext(ctx).Where("token = ?",token).First(&invite)
	if result.Error != nil{
		return nil,result.Error
	}
	return &invite, nil

}

// DeleteInviteByID deletes a invite by id and return error if any
func (p *postgresDBRepo) DeleteInviteByID(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	result:= p.DB.WithContext(ctx).Where("id = ?",id).Delete(&model.Invite{})
	if result.Error != nil{
		return result.Error
	}

	return nil
}

//UpdateInviteByID updates a invite by id and return error if any
func (p *postgresDBRepo) UpdateInviteByID(id uint,status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	invite:=model.Invite{
		Status: status,
	}

	result:= p.DB.WithContext(ctx).Model(&invite).Where("id = ?",id).Updates(&invite)
	if result.Error != nil{
		return result.Error
	}

	return nil
}

//Membership related

//InsertMembership inserts a new membership into the database and return the id of the membership and error if any
func (p *postgresDBRepo) InsertMembership(projectID uint,userID string,roleId uint) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	membership:=model.Membership{
		ProjectID: projectID,
		UserID: userID,
	}

	result:= p.DB.WithContext(ctx).Create(&membership)
	if result.Error != nil{
		return 0, result.Error
	}

	return membership.ID, nil
}

