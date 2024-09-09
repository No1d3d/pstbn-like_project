package storage

import (
	"database/sql"
	"log"

	m "cdecode/models"
)

const (
	createAliasTableSQL = `CREATE TABLE IF NOT EXISTS alias (
		"id_alias" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"id_user" integer NOT NULL,		
		"id_resource" integer NOT NULL,		
		"name" TEXT	
	  );` // SQL Statement for Create Table
)

func GetAliases(db *sql.DB, username string) []m.Alias {
	var row *sql.Rows
	var err error
	var q string
	if UserIsAdmin(db, username) {
		q = `SELECT * FROM alias`
		row, err = db.Query(q)
	} else {
		q = `SELECT * FROM alias WHERE id_user = ?`
		row, err = db.Query(q, GetUserId(db, username))
	}
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	aliases := []m.Alias{}
	for row.Next() {
		alias, scan_err := getAliasFromRow(row)
		if scan_err != nil {
			log.Println(scan_err)
		}
		aliases = append(aliases, *alias)
	}

	return aliases
}

func getAliasFromRow(row *sql.Rows) (*m.Alias, error) {
	alias := &m.Alias{}
	if err := row.Scan(&alias.Id, &alias.CreatorId, &alias.ResourceId, &alias.Name); err != nil {
		return nil, err
	}

	return alias, nil
}

// Check no duplicated name of alias
func checkIfAliasHasNoDupes(db *sql.DB, alias string) bool {
	var count int
	q := `SELECT COUNT(name) FROM alias WHERE name = ?`
	row := db.QueryRow(q, alias)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

// Creates new alias
func CreateAlias(db *sql.DB, username string, resourceName string, alias string) {
	if !checkIfAliasHasNoDupes(db, alias) {
		log.Printf("Alias '%s' already exist", alias)
		return //alias name already exists
	}
	newAlias := m.Alias{
		CreatorId:  GetUserId(db, username),
		ResourceId: GetResourceId(db, username, resourceName),
		Name:       alias,
	}
	insertAlias(db, newAlias)
}

// Creates new entity in database
func insertAlias(db *sql.DB, alias m.Alias) {
	log.Println("Adding new alias...")
	insertAliasSQL := "INSERT INTO alias(id_user, id_resource, name) VALUES (?, ?, ?)"
	statement, err := db.Prepare(insertAliasSQL) // Prepare statement.
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(alias.CreatorId, alias.ResourceId, alias.Name)
	if err != nil {
		log.Fatalln(err.Error())
	}
	log.Printf("Alias '%s' for resource '%d' added", alias.Name, alias.ResourceId)
}

func ReadContentByAlias(db *sql.DB, alias string) string {
	var row *sql.Rows
	q := `SELECT content FROM resources WHERE id_resource = (SELECT id_resource FROM alias WHERE name = ?)`
	row, err := db.Query(q, alias)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		log.Printf("No resource was found with alias '%s'", alias)
		return ""
	}
	var content string

	if scan_err := row.Scan(&content); err != nil {
		log.Fatal(scan_err)
	}
	log.Printf("Content: %s", content)
	return content
}
