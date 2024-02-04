package driver

import "gorm.io/gorm"

// DB gorm connector

type DB struct{
	SQL *gorm.DB
}

var dbCon= &DB{}