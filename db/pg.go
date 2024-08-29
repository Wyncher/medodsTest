package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"main/config"
	"main/utils"
)

func Conn() *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", config.Host, config.Port, config.User, config.Password, config.Dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	utils.CheckError(err)
	return db
}
