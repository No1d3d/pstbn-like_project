package storage

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var DbFile = "./pkg/storage/mydb.db"

type DatabaseResolver = func() AppDatabase

func GetDB() AppDatabase {
	conn, err := sql.Open("sqlite3", DbFile)
	if err != nil {
		log.Fatal(err)
	}

	return AppDatabase{connection: conn}
}

type AppDatabase struct {
	connection *sql.DB
}

func (db *AppDatabase) Init() {
	db.createTable(createUsersTableSQL)
	db.createTable(createResourcesTableSQL)
	db.createTable(createAliasTableSQL)
}

func (db *AppDatabase) Close() {
	log.Println("Closing DB connection")
	if err := db.connection.Close(); err != nil {
		log.Println(err)
	}
}

func (db *AppDatabase) createTable(query string) error {
	statement, err := db.connection.Prepare(query) // Prepare SQL Statement
	if err != nil {
		return err
	}
	_, err = statement.Exec() // Execute SQL Statements
	return err
}
