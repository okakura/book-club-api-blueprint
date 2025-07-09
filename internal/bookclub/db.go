package bookclub

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./bookclub.db")
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	// Create tables if they don't exist
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS clubs (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS members (
		club_id TEXT,
		name TEXT,
		PRIMARY KEY (club_id, name),
		FOREIGN KEY (club_id) REFERENCES clubs(id)
	);
	CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		club_id TEXT,
		title TEXT,
		author TEXT,
		votes INTEGER DEFAULT 0,
		FOREIGN KEY (club_id) REFERENCES clubs(id)
	);
	`
	_, err = DB.Exec(createTablesSQL)
	if err != nil {
		log.Fatal("Failed to create tables:", err)
	}
}
