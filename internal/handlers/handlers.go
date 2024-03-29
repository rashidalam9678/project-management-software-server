package handlers

import (
	"net/http"
	"fmt"
	"github.com/rashidalam9678/project-management-software-server/internal/config"
	"github.com/rashidalam9678/project-management-software-server/internal/driver"
	"github.com/rashidalam9678/project-management-software-server/internal/helpers"
	"github.com/rashidalam9678/project-management-software-server/internal/repository"
	"github.com/rashidalam9678/project-management-software-server/internal/repository/dbrepo"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.Database
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:dbrepo.NewPostgresRepo(db.SQL,a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}


// Home is the handler for the home page
func (m *Repository) Ping(w http.ResponseWriter, r *http.Request) {

	payLoad:=jsonResponse{}
	payLoad.Error=false
	payLoad.Message="Alive"

	err:=helpers.WriteJSON(w,http.StatusOK,payLoad)
	if err != nil{
		m.App.ErrorLog.Println(err)
	}
}

// CreateNewUser insert new user in the database
func (m *Repository) CreateNewUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin","*")

	type credentials struct {
		Email string `json:"email"`
		ID string `json:"id"`
		FirstName string `json:"first_name"`
		LastName string `json:"last_name"`
	}

	var creds credentials
	var payload jsonResponse

	err := helpers.ReadJSON(w, r, &creds)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}

	_,err= m.DB.GetUserByID(creds.ID)

	if err==nil  {
		m.App.ErrorLog.Println(err)
		payload.Error=true
		payload.Message="user already exist"
		helpers.WriteJSON(w,http.StatusAccepted,payload)
		return
	}

	//now insert user in database
	userId,err:=m.DB.InsertUser(creds.Email,creds.ID, creds.FirstName, creds.LastName)

	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error=true
		payload.Message=fmt.Sprint(err)
		helpers.WriteJSON(w,http.StatusForbidden,payload)
		return
	}

	payload.Data=userId
	payload.Error=false
	payload.Message="success"

	err=helpers.WriteJSON(w,http.StatusOK,payload)
	if err != nil{
		m.App.ErrorLog.Println(err)
	}
}

// DeleteUser handler deletes the user from the database
func (m *Repository) DeleteUser(w http.ResponseWriter, r *http.Request) {

	type credentials struct {
		ID string `json:"id"`
	}

	var creds credentials
	var payload jsonResponse

	err := helpers.ReadJSON(w, r, &creds)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}

	user,err:=m.DB.GetUserByID(creds.ID)
	if err!=nil{
		m.App.ErrorLog.Println(err)
		payload.Error=true
		payload.Message="user not found"
		helpers.WriteJSON(w,http.StatusAccepted,payload)
		return
	}

	err=m.DB.DeleteUserByID(user.ID)
	if err!=nil{
		m.App.ErrorLog.Println(err)
		payload.Error=true
		payload.Message="unable to delete user"
		helpers.WriteJSON(w,http.StatusAccepted,payload)
		return
	}

	payload.Error=false
	payload.Message="user deleted successfully"
	helpers.WriteJSON(w,http.StatusOK,payload)
}
