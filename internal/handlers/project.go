package handlers

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"github.com/rashidalam9678/project-management-software-server/internal/helpers"
)

// CreateNewProject Handler creates a new project
func (m *Repository) CreateNewProject(w http.ResponseWriter, r *http.Request) {

	// Retrieve the attached data from the context
	userID, ok := r.Context().Value("externalId").(string)
	if !ok {
		payLoad := jsonResponse{}
		payLoad.Error = true
		payLoad.Message = "external id context not found"
		helpers.WriteJSON(w, http.StatusInternalServerError, payLoad)
		return
	}

	// extract the request body
	type credentials struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	var cred credentials
	var payload jsonResponse

	err := helpers.ReadJSON(w, r, &cred)
	if err != nil {
		if err != nil {
			m.App.ErrorLog.Println(err)
			payload.Error = true
			payload.Message = "invalid json supplied, or json missing entirely"
			_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
		}
	}

	// create the project
	projectID, err := m.DB.InsertProject(cred.Title, cred.Description, userID)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to create project"
		_ = helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	//Create the guest role for the project
	err = m.DB.InsertDefaultRoleAndPermissions(projectID)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to create default role and permissions"
		_ = helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	// create the response
	payload.Data = projectID
	payload.Error = false
	payload.Message = "project created successfully"

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
}

// GetProjects Handler gets all the projects
func (m *Repository) GetUserProjects(w http.ResponseWriter, r *http.Request) {
	
	// Retrieve the attached data from the context
	userID, ok := r.Context().Value("externalId").(string)
	if !ok {
		payLoad := jsonResponse{}
		payLoad.Error = true
		payLoad.Message = "external id context not found"
		helpers.WriteJSON(w, http.StatusInternalServerError, payLoad)
		return
	}

	// get the projects
	projects, err := m.DB.GetAllProjectsByUserID(userID)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload := jsonResponse{}
		payload.Error = true
		payload.Message = "unable to get projects"
		helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	// create the response
	payload := jsonResponse{}
	payload.Data = projects
	payload.Error = false
	payload.Message = "success"

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
}

// GetProjectByID Handler gets a project by id
func (m *Repository) GetProject(w http.ResponseWriter, r *http.Request) {
	//retrive project id from the url
	vars:=mux.Vars(r)
	id:=vars["id"]

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	projectId := uint(uint64Val)

	var payload jsonResponse

	// get the project
	project, err := m.DB.GetProjectByID(projectId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "project not found"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// create the response
	payload.Data = project
	payload.Error = false
	payload.Message = "success"

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
}

// DeleteProject Handler deletes a project by id
func (m *Repository) DeleteProject(w http.ResponseWriter, r *http.Request) {
	//retrive project id from the url
	vars:=mux.Vars(r)
	id:=vars["id"]

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	projectId := uint(uint64Val)

	var payload jsonResponse

	//check if the project exists or not
	_, err = m.DB.GetProjectByID(projectId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "project not found"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// delete the project
	err = m.DB.DeleteProjectByID(projectId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to delete project"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// create the response
	payload.Error = false
	payload.Message = "project deleted successfully"
	helpers.WriteJSON(w, http.StatusOK, payload)
}

// UpdateProject Handler updates a project by id
func (m *Repository) UpdateProject(w http.ResponseWriter, r *http.Request) {
	//retrive project id from the url
	vars:=mux.Vars(r)
	id:=vars["id"]

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	projectId := uint(uint64Val)

	// extract the request body
	type credentials struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	var cred credentials
	var payload jsonResponse

	err = helpers.ReadJSON(w, r, &cred)
	if err != nil {
		if err != nil {
			m.App.ErrorLog.Println(err)
			payload.Error = true
			payload.Message = "invalid json supplied, or json missing entirely"
			_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
		}
	}

	//check if the project exists or not
	_, err = m.DB.GetProjectByID(projectId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "project not found"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// update the project
	err = m.DB.UpdateProjectByID(cred.Title, cred.Description, projectId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to update project"
		_ = helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	// create the response
	payload.Error = false
	payload.Message = "project updated successfully"
	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}
