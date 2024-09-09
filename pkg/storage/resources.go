package storage

import (
	"cdecode/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const (
	ResourceIdColumn      = "id_resource"
	ResourceCreatorColumn = "id_user"
	ResourceContentColumn = "content"

	createResourcesTableSQL = `CREATE TABLE IF NOT EXISTS resources (
		"` + ResourceIdColumn + `" integer NOT NULL PRIMARY KEY AUTOINCREMENT
		,"` + ResourceCreatorColumn + `" integer NOT NULL
		,"` + ResourceContentColumn + `" TEXT
	  );` // SQL Statement for Create Table
)

func CreateResource(db *sql.DB, creatorId models.UserId, content string) (*models.Resource, error) {
	res := &models.Resource{
		UserId:  creatorId,
		Content: content,
	}
	insertResources(db, res)

	return res, nil
}

func insertResources(db *sql.DB, r *models.Resource) {
	log.Println("Inserting resources record ...")

	query := `INSERT INTO resources
    (` + ResourceCreatorColumn + `, ` + ResourceContentColumn + `) 
    VALUES (?, ?)`
	statement, err := db.Prepare(query) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(r.UserId, r.Content)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func GetResources(db *sql.DB, userId models.UserId) []*models.Resource {
	// if UserIsAdmin(db, username) {
	// q = `SELECT * FROM resources`
	// row, err = db.Query(q)
	// } else {
	q := "SELECT " + ResourceIdColumn + ", " + ResourceCreatorColumn + ", " + ResourceContentColumn + ` FROM resources WHERE id_user = ?`
	row, err := db.Query(q, userId)
	// }
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	resources := []*models.Resource{}
	for row.Next() {
		resource, err := getResourceFromRow(row)
		if err != nil {
			log.Println(err)
		}
		resources = append(resources, resource)
	}
	return resources
}

func getResourceFromRow(row *sql.Rows) (*models.Resource, error) {
	resource := &models.Resource{}
	if err := row.Scan(&resource.Id, &resource.UserId, &resource.Content); err != nil {
		return nil, err
	}

	return resource, nil
}
func GetResourceById(db *sql.DB, id models.ResourceId) (*models.Resource, error) {
	query := "SELECT * FROM resources WHERE " + ResourceIdColumn + " = ?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(id)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		res, err := getResourceFromRow(row)
		row.Close()
		return res, err
	}
	return nil, errors.New(fmt.Sprintf("No such resource with id %d", id))
}

func DeleteResource(db *sql.DB, id models.ResourceId) error {
	q := `DELETE FROM resources WHERE ` + ResourceIdColumn + ` = ?`
	_, err := db.Exec(q, id)
	return err
}

func UpdateResource(db *sql.DB, r models.Resource) error {
	q := `UPDATE resources SET ` + ResourceContentColumn + ` = ?  WHERE ` + ResourceIdColumn + ` = ?`
	_, err := db.Exec(q, r.Content, r.Id)
	return err
}
