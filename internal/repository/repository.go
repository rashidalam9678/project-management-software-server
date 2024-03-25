package repository

import (

	model "github.com/rashidalam9678/project-management-software-server/internal/models"
)

type Database interface{
	GetUserByEmail(email string) (*model.User,error)
	InsertUser(email,externalId, firstName, lastName string)(string,error)
	GetUserByID(externalId string) (*model.User, error)
	UpdateUserByID(email,id string) error
	DeleteUserByID(id string) error

	//project related
	InsertProject(title,description,userId string)(uint,error)
	GetProjectByID(id uint)(*model.Project,error)
	GetAllProjectsByUserID(userId string)([]model.Project,error)
	UpdateProjectByID(title,description string,id uint) error
	DeleteProjectByID(id uint) error

	//Role related
	InsertDefaultRoleAndPermissions(projectID uint) error
	InsertRole(projectID uint,roleName string) (uint, error)

	// InsertInvite(email,projectId,roleId string)(uint,error)
	GetInviteByEmailAndProjectID(email string, projectId uint)(*model.Invite,error)
	InsertInvite(email string,projectID uint,invitedBy string, description string,token string) (uint, error)
	GetInviteByID(id uint) (*model.Invite, error)
	GetInviteByToken(token string) (*model.Invite, error)
	GetInvitesByProjectID(projectID uint) ([]model.Invite, error)
	DeleteInviteByID(id uint) error
	UpdateInviteByID(id uint,status string) error

	InsertMembership(projectID uint,userID string,roleId uint) (uint, error)

	// GetInviteByID(id uint)(*model.Invite,error)
	// GetAllInvitesByProjectID(projectId uint)([]model.Invite,error)

	// Issue related functions
    CreatedIssue(ProjectID uint , CreatedBy string , Title string , Description string , AssignedTo string , Status string , Priority string)(uint , error)
	GetAllIssues(ProjectId uint)([](model.Task),error)
	UpdateIssueById(Title string , Description string , AssignedTo string , Status string,issuesId uint)(error)
	DeleteIssueById(IssueId uint)(error)
	GetIssueByID(IssueId uint)(error)


	// chat related functions
	SaveMessage(ProjectID uint , SenderID uint , Message string)(error)
	GetAllMembersByProjectID(ProjectID uint) ([]model.Membership, error)
	GetMessagesByProjectId(ProjectID uint) ([]model.Message, error)

}
