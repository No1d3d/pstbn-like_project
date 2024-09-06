package storage

import (
	"cdecode/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const (
	ResourceIdColumn      = "id_resource"
	ResourceCreatorColumn = "id_user"
	ResourceNameColumn    = "name"
	ResourceContentColumn = "content"

	createResourcesTableSQL = `CREATE TABLE IF NOT EXISTS resources (
		"` + ResourceIdColumn + `" integer NOT NULL PRIMARY KEY AUTOINCREMENT
		,"` + ResourceCreatorColumn + `" integer NOT NULL
		,"` + ResourceNameColumn + `" TEXT UNIQUE
		,"` + ResourceContentColumn + `" TEXT
	  );` // SQL Statement for Create Table
)

func CreateResource(db *sql.DB, creatorId models.UserId, name string, content string) (*models.Resource, error) {
	if !CheckIfResourceHasNoDupNames(db, name) {
		return nil, errors.New(fmt.Sprintf("Resource with name '%s' already exist", name))
	}
	res := &models.Resource{
		UserId:  creatorId,
		Name:    name,
		Content: content,
	}
	insertResources(db, res)

	return res, nil
}

func insertResources(db *sql.DB, r *models.Resource) {
	log.Println("Inserting resources record ...")

	query := `INSERT INTO resources
    (` + ResourceCreatorColumn + `, ` + ResourceNameColumn + `, ` + ResourceContentColumn + `) 
    VALUES (?, ?, ?)`
	statement, err := db.Prepare(query) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(r.UserId, r.Name, r.Content)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GetResources(db *sql.DB, userId models.UserId) []*models.Resource {
	// if UserIsAdmin(db, username) {
	// q = `SELECT * FROM resources`
	// row, err = db.Query(q)
	// } else {
	q := `SELECT * FROM resources WHERE id_user = ?`
	row, err := db.Query(q, userId)
	// }
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	resources := []*models.Resource{}
	for row.Next() {
		resource, scan_err := getResourceFromRow(row)
		if scan_err != nil {
			log.Println(scan_err)
		}
		resources = append(resources, resource)
	}
	return resources
}
