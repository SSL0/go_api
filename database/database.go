package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "am"
	password = "postgres"
	dbname   = "godb"
)

var db_params string = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
	host, port, user, password, dbname)

var DB *gorm.DB

func Connect() {
	db, err := gorm.Open(postgres.Open(db_params), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
	fmt.Print(db)

	db.AutoMigrate(&User{})
}
