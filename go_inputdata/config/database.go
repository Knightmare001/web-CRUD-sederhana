package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func DBConn() (db *sql.DB, err error) {

	db, err = sql.Open("mysql", "root:@/mahasiswa")
	return
}
