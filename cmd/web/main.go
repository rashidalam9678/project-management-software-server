package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/rashidalam9678/project-management-software-server/internal/config"
	"github.com/rashidalam9678/project-management-software-server/internal/driver"
	"github.com/rashidalam9678/project-management-software-server/internal/handlers"
)

var app config.AppConfig

func main(){
	//load envs
	godotenv.Load()
	clerk_key := os.Getenv("CLERK_PRIVATE_KEY")
	// sendgrid_key:=os.Getenv("SENDGRID_API_KEY")
	port:=os.Getenv("PORT")

	//Connect to Databas and intitalise the configs
	db,err:= driver.ConnectDB()

	if err!=nil{
		log.Fatal("Not Connected to DB")
	}

	infoLog:= log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	errorLog:= log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime|log.Lshortfile)

	app.ErrorLog=errorLog
	app.InfoLog=infoLog

	repo:=handlers.NewRepo(&app,db)
	handlers.NewHandlers(repo)

	//create the server
	srv:= &http.Server{
		Addr: ":"+port,
		Handler: routes(&app,clerk_key),
	}

	//start the http server
	fmt.Println("Starting the server at port 8080")
	err=srv.ListenAndServe()
	if err !=nil{
		log.Fatal("Unable to start the server")
	}

}
