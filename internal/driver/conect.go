package driver

import (
	"fmt"
	"log"
	"os"

	"github.com/rashidalam9678/project-management-software-server/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*DB,error) {
	var err error // define error here to prevent overshadowing the global DB

	env := os.Getenv("DATABASE_URL")

	fmt.Println("connecting to Database...")
	db, err := gorm.Open(postgres.Open(env), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--------connected ---------")

	fmt.Println("running migrations...")


	dbCon.SQL=db

	err = db.AutoMigrate(&model.User{},&model.Project{} )
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("--------migration completed---------")

	return dbCon,nil


}