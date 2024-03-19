package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rashidalam9678/project-management-software-server/internal/helpers"
)

func (m *Repository) SendInvite(w http.ResponseWriter, r *http.Request) {
	//Retrieve the attached data from the context
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
		Email       string `json:"email"`
		ProjectID   string `json:"project_id"`
		Description string `json:"description"`
		RoleID      string `json:"role_id"`
	}
	var cred credentials
	var payload jsonResponse
	err := helpers.ReadJSON(w, r, &cred)

	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "invalid json supplied, or json missing entirely"
		_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
	}

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(cred.ProjectID, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	projectId := uint(uint64Val)
	// check if the invitation already sent
	invit, err := m.DB.GetInviteByEmailAndProjectID(cred.Email, projectId)
	if invit != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "project invite already exist"
		_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// generate a token
	token, _ := helpers.GenerateToken()

	// make the entry in the database
	invite, err := m.DB.InsertInvite(cred.Email, projectId, userID, cred.Description, token)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to send invite"
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	//send the email for invitation
	_, err = helpers.SendInviteMail(cred.Email, cred.RoleID, token)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to send invite"
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	// create the response
	payload.Data = invite
	payload.Error = false
	payload.Message = "invite sent successfully"
	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
}

// AcceptInvite handler accepts the invite and add the user to the project
func (m *Repository) AcceptInvite(w http.ResponseWriter, r *http.Request) {
	//Retrieve the attached data from the context
	userID, ok := r.Context().Value("externalId").(string)
	if !ok {
		payLoad := jsonResponse{}
		payLoad.Error = true
		payLoad.Message = "external id context not found"
		helpers.WriteJSON(w, http.StatusInternalServerError, payLoad)
		return
	}
	// extract token from the url
	token := r.URL.Query().Get("token")
	roleId := r.URL.Query().Get("role")

	// Convert string to uint64
	uint64Val, err := strconv.ParseUint(roleId, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	roleID := uint(uint64Val)

	var payload jsonResponse
	if token == "" {
		payload.Error = true
		payload.Message = "token not found"
		_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
		return

	}

	// check if the invitation exists
	invite, err := m.DB.GetInviteByToken(token)
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "project invite not found"
		_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// check if the invite is already accepted or rejected
	if invite.Status != "pending" {
		payload.Error = true
		payload.Message = "project invite already accepted or rejected"
		_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// check if the invite is expired
	if invite.ExpiresAt.Before(time.Now()) {
		payload.Error = true
		payload.Message = "project invite expired, ask admin to send again"
		_ = helpers.WriteJSON(w, http.StatusAccepted, payload)
		return
	}

	// update invite status
	err = m.DB.UpdateInviteByID(invite.ID, "accepted")
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to accept invite"
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
		}



	// Add them to to the project as a guest member
	_, err = m.DB.InsertMembership(invite.ProjectID, userID, roleID)

	// Add their default role permissions
	if err != nil {
		m.App.ErrorLog.Println(err)
		payload.Error = true
		payload.Message = "unable to accept invite"
		_ = helpers.WriteJSON(w, http.StatusInternalServerError, payload)
		return
	}

	// create the response
	payload.Error = false
	payload.Message = "invite accepted successfully"
	payload.Data = invite.ProjectID
	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

}

func (m *Repository) GetAllInvites(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	projectID := uint(uint64Val)

	invites, err := m.DB.GetInvitesByProjectID(projectID)
	payload := jsonResponse{}
	if err != nil {
		payload.Error = true
		payload.Message = "unable to get data"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
	}

	payload.Data = invites
	payload.Error = false
	payload.Message = "find all invites"
	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

}

func (m *Repository) DeleteInvite(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	payload := jsonResponse{}

	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	ID := uint(uint64Val)

	invite, _ := m.DB.GetInviteByID(ID)

	if invite != nil {
		payload.Error = true
		payload.Message = "invite does not exists"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
	}

	err = m.DB.DeleteInviteByID(ID)
	if err != nil {
		payload.Error = true
		payload.Message = "unable to delete invite"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
	}

	payload.Error = false
	payload.Message = "invite deleted successfully"
	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}
}

// func (m *Repository) ResendInvite(w http.ResponseWriter, r *http.Request) {
// 	//Retrieve the attached data from the context
// 	userID, ok := r.Context().Value("externalId").(string)
// 	if !ok {
// 		payLoad := jsonResponse{}
// 		payLoad.Error = true
// 		payLoad.Message = "external id context not found"
// 		helpers.WriteJSON(w, http.StatusInternalServerError, payLoad)
// 		return
// 	}

// 	// extract the request body
// 	type credentials struct {
// 		Email       string `json:"email"`
// 		ProjectID uint `json:"project_id"`
// 		RoleID string `json:"role_id"`
// 	}
// 	var cred credentials
// 	var payload jsonResponse
// 	err := helpers.ReadJSON(w, r, &cred)
// 	if err != nil {
// 		if err != nil {
// 			m.App.ErrorLog.Println(err)
// 			payload.Error = true
// 			payload.Message = "invalid json supplied, or json missing entirely"
// 			_ = helpers.WriteJSON(w, http.StatusBadRequest, payload)
// 		}
// 	}

// 	// delete the previous invite

// 	//create the new invite

// 	//send the email for invitation

// 	// send the response back

// 	// create the project

// 	// create the response
// 	payload.Data = projectID
// 	payload.Error = false
// 	payload.Message = "project created successfully"
// 	err = helpers.WriteJSON(w, http.StatusOK, payload)
// 	if err != nil {
// 		m.App.ErrorLog.Println(err)
// 	}
// }
