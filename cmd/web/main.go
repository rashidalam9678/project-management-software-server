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

	//Connect to Databas and intitalise the configs
	db,err:= driver.ConnectDB()

	if err!=nil{
		log.Fatal("Not Connected to DB")
	}

	infoLog:= log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	errorLog:= log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime|log.Lshortfile)

	app.ErrorLog=errorLog
	app.InfoLog=infoLog
	app.Port=8080

	repo:=handlers.NewRepo(&app,db)
	handlers.NewHandlers(repo)

	//create the server
	srv:= &http.Server{
		Addr: ":8080",
		Handler: routes(&app),
	}

	//start the http server
	fmt.Println("Starting the server at http://localhost:8080")
	err=srv.ListenAndServe()
	if err !=nil{
		log.Fatal("Unable to start the server")
	}

}
