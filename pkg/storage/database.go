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

func UserIsAdmin(db *sql.DB, username string) bool {
	var row *sql.Rows
	var err error
	q := `SELECT is_admin FROM users WHERE name = ?`
	row, err = db.Query(q, username)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var is_admin bool
	for row.Next() {
		if scan_err := row.Scan(&is_admin); err != nil {
			log.Fatal(scan_err)
		}
	}
	return is_admin

}

func GetUserId(db *sql.DB, username string) int {
	var row *sql.Rows
	var err error
	q := `SELECT id_user FROM users WHERE name = ?`
	row, err = db.Query(q, username)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var id_user int
	for row.Next() {
		if scan_err := row.Scan(&id_user); err != nil {
			log.Fatal(scan_err)
		}
	}
	return id_user
}

// check no duplicated name of resource for user
func CheckIfResourceHasNoDupNames(db *sql.DB, resourceName string) bool {
	var count int
	q := `SELECT COUNT(name) FROM resources WHERE name = ?`
	row := db.QueryRow(q, resourceName)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

// create resource

func GetResourceId(db *sql.DB, username, resourceName string) int {
	var row *sql.Rows
	var err error
	q := `SELECT id_resource FROM resources WHERE id_user = ? AND name = ?`
	row, err = db.Query(q, GetUserId(db, username), resourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var id_resource int
	for row.Next() {
		if scan_err := row.Scan(&id_resource); err != nil {
			log.Fatal(scan_err)
		}
	}
	return id_resource
}
