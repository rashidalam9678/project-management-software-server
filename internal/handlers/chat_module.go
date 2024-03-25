package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
	"github.com/rashidalam9678/project-management-software-server/internal/helpers"
)



type Message struct {
	ProjectID int       `json:"project_id"`
	SenderID  int       `json:"sender_id"`
	Message   string    `json:"message"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan Message)

func (m *Repository) CreateMessage(w http.ResponseWriter, r *http.Request){
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
    var payload jsonResponse
	clients[ws] = true
	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		broadcast <- msg
		resErr := m.DB.SaveMessage(uint(msg.ProjectID) , uint(msg.SenderID), msg.Message)

		if resErr != nil {
			m.App.ErrorLog.Println(err)
			payload.Error = true
			payload.Message = "unable to save the message in the backend"
			_ = helpers.WriteJSON(w, http.StatusForbidden, payload)
			return
		}

	}

}


func (m *Repository) GetMessages(w http.ResponseWriter , r *http.Request){
	id := r.URL.Query().Get("id")
	uint64Val, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		m.App.ErrorLog.Println(err)
		return
	}
	projectID := uint(uint64Val)


	messages, err := m.DB.GetMessagesByProjectId(projectID)
	payload := jsonResponse{}
	if err != nil {
		payload.Error = true
		payload.Message = "unable to get messages"
		helpers.WriteJSON(w, http.StatusAccepted, payload)
	}

	payload.Data = messages
	payload.Error = false
	payload.Message = "find all messages"
	err = helpers.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		m.App.ErrorLog.Println(err)
	}

}
