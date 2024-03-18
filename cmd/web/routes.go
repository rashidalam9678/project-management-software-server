package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	middlewareHandler "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rashidalam9678/project-management-software-server/internal/config"
	"github.com/rashidalam9678/project-management-software-server/internal/handlers"
	"github.com/rashidalam9678/project-management-software-server/internal/helpers"
)

// jsonResponse is the type used for generic JSON responses
type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func routes(app *config.AppConfig, clerkKey string) http.Handler {
	mux := mux.NewRouter()
	//get the private key from the environment


	//create a new clerk client
	client, err := clerk.NewClient(clerkKey)
	if err != nil {
		fmt.Println("Error in clerk")
	}

	// member invite routes


	//public routes
	mux.HandleFunc("/", handlers.Repo.Ping)
	mux.HandleFunc("/user", handlers.Repo.CreateNewUser).Methods("POST")
	mux.HandleFunc("/user", handlers.Repo.DeleteUser).Methods("DELETE")


	// mux.HandleFunc("/invite", handlers.Repo.SendInvite).Methods("POST")

	// protected routes
	subrouter := mux.PathPrefix("/api/v1").Subrouter()
	subrouter.Use(AuthenticateToken(client))

	//project routes
	subrouter.HandleFunc("/projects", handlers.Repo.CreateNewProject).Methods("POST")
	subrouter.HandleFunc("/projects", handlers.Repo.GetUserProjects).Methods("GET")
	subrouter.HandleFunc("/projects/{id}", handlers.Repo.GetProject).Methods("GET")
	subrouter.HandleFunc("/projects/{id}/invites", handlers.Repo.GetAllInvites).Methods("GET")
	subrouter.HandleFunc("/projects/{id}", handlers.Repo.DeleteProject).Methods("DELETE")
	subrouter.HandleFunc("/projects/{id}", handlers.Repo.UpdateProject).Methods("PUT")

	//project invite routes
	subrouter.HandleFunc("/invite/send", handlers.Repo.SendInvite).Methods("POST")
	subrouter.HandleFunc("/invite/confirm", handlers.Repo.AcceptInvite).Methods("POST")
	subrouter.HandleFunc("/invite/{id}", handlers.Repo.DeleteInvite).Methods("DELETE")
	// subrouter.HandleFunc("/reject-invite", handlers.Repo.RejectInvite).Methods("PUT")
	// subrouter.HandleFunc("/invite", handlers.Repo.DeleteInvite).Methods("DELETE")


	// Issues routes
	subrouter.HandleFunc("/issue", handlers.Repo.CreateIssue).Methods("POST")
	subrouter.HandleFunc("/issue", handlers.Repo.GetIssuesByProjectsID).Methods("GET")
	subrouter.HandleFunc("/issue", handlers.Repo.UpdateIssueByID).Methods("PUT")
	subrouter.HandleFunc("/issue", handlers.Repo.DeleteIssueByID).Methods("DELETE")



	// Create a new CORS middleware with a few options
	corsHandler := middlewareHandler.CORS(
		middlewareHandler.AllowedOrigins([]string{"*"}), // Adjust as needed for your frontend's origin
		middlewareHandler.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		middlewareHandler.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	// Create a new http.Handler that includes the CORS middleware
	handler := corsHandler(mux)

	return handler
}

// AuthenticateToken middleware authenticate the token and add the external id from clerk to the context
func AuthenticateToken(client clerk.Client) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//extract the token
			token, err := ExtractToken(r)
			if err != nil {
				payload := jsonResponse{
					Error:   true,
					Message: "Invalid token supplied",
				}
				_ = helpers.WriteJSON(w, http.StatusUnauthorized, payload)
				return
			}
			// verify the session
			sessClaims, err := client.VerifyToken(token)
			if err != nil {
				payload := jsonResponse{
					Error:   true,
					Message: "Unauthorised access",
				}
				_ = helpers.WriteJSON(w, http.StatusUnauthorized, payload)
				return
			}

			// get the user
			user, err := client.Users().Read(sessClaims.Claims.Subject)
			if err != nil {
				payload := jsonResponse{
					Error:   true,
					Message: "Unable to get the user with token",
				}
				_ = helpers.WriteJSON(w, http.StatusUnauthorized, payload)
				return
			}
			//add external id from clerk to context
			ctx := context.WithValue(r.Context(), "externalId", user.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ExtractToken(r *http.Request) (string, error) {
	// get the authorization header
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return "", errors.New("no authorization header received")
	}

	// get the plain text token from the header
	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("no valid authorization header received")
	}

	token := headerParts[1]
	return token, nil
}
