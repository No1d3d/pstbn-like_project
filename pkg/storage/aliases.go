package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"cdecode/pkg/models"
	m "cdecode/pkg/models"
)

const (
	AliasIdColumn       = "id_alias"
	AliasCreatorColumn  = "id_user"
	AliasNameColumn     = "name"
	AliasResourceColumn = "id_resource"

	createAliasTableSQL = `CREATE TABLE IF NOT EXISTS alias (
		"id_alias" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"id_user" integer NOT NULL,		
		"id_resource" integer NOT NULL,		
		"name" TEXT	
	  );` // SQL Statement for Create Table
)

func GetAliases(db *sql.DB, user_id m.UserId) []m.Alias {
	var row *sql.Rows
	var err error
	var q string
	q = `SELECT ` + AliasIdColumn + ", " + AliasCreatorColumn + ", " + AliasResourceColumn + ", " + AliasNameColumn + ` FROM alias WHERE id_user = ?`
	row, err = db.Query(q, user_id)

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

// new create and insert aliases
func CreateAlias(db *sql.DB, creatorId m.UserId, name string, resourceId m.ResourceId) (*m.Alias, error) {
	alias := &m.Alias{
		CreatorId:  creatorId,
		Name:       name,
		ResourceId: resourceId,
	}
	insertAlias(db, alias)

	return alias, nil
}

func insertAlias(db *sql.DB, a *m.Alias) {
	log.Println("Inserting resources record ...")

	query := `INSERT INTO alias
    (` + AliasCreatorColumn + `, ` + AliasNameColumn + ", " + AliasResourceColumn + `) 
    VALUES (?, ?, ?)`
	statement, err := db.Prepare(query) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(a.CreatorId, a.Name, a.ResourceId)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GetAliasById(db *sql.DB, id m.AliasId) (*m.Alias, error) {
	query := "SELECT * FROM alias WHERE " + AliasIdColumn + " = ?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(id)

	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		alias, err := getAliasFromRow(row)
		return alias, err
	}
	return nil, errors.New(fmt.Sprintf("No such alias with id %d", id))
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

func DeleteAlias(db *sql.DB, id models.AliasId) error {
	q := `DELETE FROM alias WHERE ` + AliasIdColumn + ` = ?`
	_, err := db.Exec(q, id)
	return err
}
