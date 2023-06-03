package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password []byte `json:"-"`
}

type Account struct {
	ID        uint
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"`
	CreatedAt time.Time `db:"created_at"`
	LastLogin time.Time `db:"last_login"`
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

var DB *sqlx.DB

func Connect() {
	db, err := sqlx.Open("pgx", db_params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	DB = db
}

func Disconnect() {
	DB.Close()
}

func QueryRow(query string, account *Account) error {
	err := DB.QueryRowx(query).StructScan(account)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Query Row: %v\n", err)
	}
	return err
}

func InsertRow(query string, account *Account) error {
	_, err := DB.NamedExec(query, account)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to Insert Row: %v\n", err)
	}
	return err

}
