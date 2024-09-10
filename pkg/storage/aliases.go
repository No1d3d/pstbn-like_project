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

func (db *AppDatabase) GetAliases(user_id m.UserId) []m.Alias {
	var row *sql.Rows
	var err error
	var q string
	q = `SELECT ` + AliasIdColumn + ", " + AliasCreatorColumn + ", " + AliasResourceColumn + ", " + AliasNameColumn + ` FROM alias WHERE id_user = ?`
	row, err = db.connection.Query(q, user_id)

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
func (db *AppDatabase) checkIfAliasHasNoDupes(alias string) bool {
	var count int
	q := `SELECT COUNT(name) FROM alias WHERE name = ?`
	row := db.connection.QueryRow(q, alias)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

// new create and insert aliases
func (db *AppDatabase) CreateAlias(creatorId m.UserId, name string, resourceId m.ResourceId) (*m.Alias, error) {
	alias := &m.Alias{
		CreatorId:  creatorId,
		Name:       name,
		ResourceId: resourceId,
	}
	db.insertAlias(alias)

	return alias, nil
}

func (db *AppDatabase) insertAlias(a *m.Alias) {
	log.Println("Inserting resources record ...")

	query := `INSERT INTO alias
    (` + AliasCreatorColumn + `, ` + AliasNameColumn + ", " + AliasResourceColumn + `) 
    VALUES (?, ?, ?)`
	statement, err := db.connection.Prepare(query) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(a.CreatorId, a.Name, a.ResourceId)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (db *AppDatabase) GetAliasById(id m.AliasId) (*m.Alias, error) {
	query := "SELECT * FROM alias WHERE " + AliasIdColumn + " = ?"
	stmt, err := db.connection.Prepare(query)

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

func (db *AppDatabase) GetAliasByName(name string) (*m.Alias, error) {
	query := "SELECT * FROM alias WHERE " + AliasNameColumn + " = ?"
	stmt, err := db.connection.Prepare(query)

	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(name)

	if err != nil {
		return nil, err
	}
	defer row.Close()
	for row.Next() {
		alias, err := getAliasFromRow(row)
		return alias, err
	}
	return nil, errors.New(fmt.Sprintf("No such alias with name '%s'", name))
}

func (db *AppDatabase) ReadContentByAlias(alias string) string {
	var row *sql.Rows
	q := `SELECT content FROM resources WHERE id_resource = (SELECT id_resource FROM alias WHERE name = ?)`
	row, err := db.connection.Query(q, alias)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		log.Printf("No resource was found with alias '%s'", alias)
		return ""
	}
	var content string

	if err = row.Scan(&content); err != nil {
		log.Fatal(err)
	}
	log.Printf("Content: %s", content)
	return content
}

func (db *AppDatabase) DeleteAlias(id models.AliasId) error {
	q := `DELETE FROM alias WHERE ` + AliasIdColumn + ` = ?`
	_, err := db.connection.Exec(q, id)
	return err
}
