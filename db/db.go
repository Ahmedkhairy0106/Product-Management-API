package db

import (
	"database/sql"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

var DB *sql.DB

func InitializeDB() (*sql.DB, error) {
	var err error

	connString := "server=localhost;user id=sa;password=Password_123#;port=1433;database=master"

	DB, err = sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	if err = DB.Ping(); err != nil {
		return nil, err
	}

	log.Println("Connected to database successfully!")
	return DB, nil
}