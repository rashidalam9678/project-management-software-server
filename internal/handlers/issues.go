package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rashidalam9678/project-management-software-server/internal/helpers"
)

func (m *Repository) CreateIssue(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value("externalId").(string)
	if !ok {
		payLoad := jsonResponse{}
		payLoad.Error = true
		payLoad.Message = "external id context not found"
		helpers.WriteJSON(w, http.StatusInternalServerError, payLoad)
		return
	}

	type task struct {
		ProjectID   uint   `json:"projectid"`
		CreatedBy   string `json:"createdby"`
		Title       string `json:"title"`
		Description string `json:"description"`
		AssignedTo  string `json:"assignedto"`
		Status      string `json:"status"`
		Priority    string `json:"priority"`
	}
	var cred task
	var payload jsonResponse

	err := helpers.ReadJSON(w, r, &cred)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}

	TaskID, err := m.DB.CreatedIssue(cred.ProjectID, cred.CreatedBy, cred.Title, cred.Description, cred.AssignedTo, cred.Status, cred.Priority)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to create Issue or task"
		_ = helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	//Create the guest role for the project

	// create the response
	payload.Data = TaskID
	payload.Error = false
	payload.Message = "Issue created successfully"

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

}

func (m *Repository) GetIssuesByProjectsID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}

	projectId := uint(uint64Val)

	_, ok := r.Context().Value("externalId").(string)
	if !ok {
		payLoad := jsonResponse{}
		payLoad.Error = true
		payLoad.Message = "external id context not found"
		helpers.WriteJSON(w, http.StatusInternalServerError, payLoad)
		return
	}

	issues, err := m.DB.GetAllIssues(projectId)

	if err != nil {
		m.App.ErrorLog.Println(err)
		payload := jsonResponse{}
		payload.Error = true
		payload.Message = "unable to get Issues"
		helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	payload := jsonResponse{}
	payload.Data = issues
	payload.Error = false
	payload.Message = "success"

	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
}

func (m *Repository) UpdateIssueByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	issueId := uint(uint64Val)

	// extract the request body
	type credentials struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		AssignedTo  string    `json:"assigned_to"`
		Status      string    `json:"status"`
	}
	var cred credentials
	var payload jsonResponse

	err = helpers.ReadJSON(w, r, &cred)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}

	//check if the issue exists or not
	err = m.DB.GetIssueByID(issueId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "issue not found"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// update the issue
	err = m.DB.UpdateIssueById(cred.Title,cred.Description,cred.AssignedTo, cred.Status, issueId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to update the issue"
		_ = helpers.WriteJSON(w, http.StatusForbidden, payload)
		return
	}

	// create the response
	payload.Error = false
	payload.Message = "issue updated successfully"
	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}

func (m *Repository) DeleteIssueByID(w http.ResponseWriter, r *http.Request){
	var payload jsonResponse
	vars := mux.Vars(r)
	id := vars["id"]

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	issueId := uint(uint64Val)

	//check if the issue exists or not
	err = m.DB.GetIssueByID(issueId)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "issue not found"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}


	delErr := m.DB.DeleteIssueById(issueId)
	if delErr != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to delete the issue"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	payload.Error = false
	payload.Message = "issue deleted successfully"
	_ = helpers.WriteJSON(w, http.StatusOK, payload)
}
