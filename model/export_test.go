package model

import "database/sql"

func SetDB(d *sql.DB) {
	db = d
}
