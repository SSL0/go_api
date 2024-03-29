package database

import (
	"fmt"
	"os"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Account struct {
	ID        uint
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  []byte    `json:"-"`
	CreatedAt time.Time `db:"created_at"`
	LastLogin time.Time `db:"last_login"`
	Balance   uint      `json:"balance"`
}

const (
	host     = "postgres"
	port     = 5432
	user     = "am"
	password = "postgres"
	dbname   = "godb"
)

var db_params string = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
	host, port, user, password, dbname)

var DB *sqlx.DB

func AutoMigrate() {
	query := `CREATE TABLE accounts (
				id Serial PRIMARY KEY,
				name text UNIQUE,
				email text UNIQUE,
				password bytea,
				created_at timestamp with time zone DEFAULT NOW(),
				last_login timestamp with time zone DEFAULT NOW(),
				balance int DEFAULT 0
			);`
	if DB == nil {
		fmt.Fprintf(os.Stderr, "ERR: No connection to database")
		os.Exit(1)
	}
	_, err := DB.Exec(query)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
}

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
