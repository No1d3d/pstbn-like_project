package storage

import (
	m "cdecode/pkg/models"
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
	  );`
)

func (db *AppDatabase) CreateResource(creatorId m.UserId, content string) (*m.Resource, error) {
	res := &m.Resource{
		UserId:  creatorId,
		Content: content,
	}
	db.insertResources(res)

	return res, nil
}

func (db *AppDatabase) insertResources(r *m.Resource) {
	log.Println("Inserting resources record ...")

	query := `INSERT INTO resources
    (` + ResourceCreatorColumn + `, ` + ResourceContentColumn + `) 
    VALUES (?, ?)`
	statement, err := db.connection.Prepare(query) // Prepare statement.
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(r.UserId, r.Content)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func (db *AppDatabase) GetResources(userId m.UserId) []*m.Resource {
	q := "SELECT " + ResourceIdColumn + ", " + ResourceCreatorColumn + ", " + ResourceContentColumn + ` FROM resources WHERE id_user = ?`
	row, err := db.connection.Query(q, userId)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	resources := []*m.Resource{}
	for row.Next() {
		resource, err := getResourceFromRow(row)
		if err != nil {
			log.Println(err)
		}
		resources = append(resources, resource)
	}
	return resources
}

func getResourceFromRow(row *sql.Rows) (*m.Resource, error) {
	resource := &m.Resource{}
	if err := row.Scan(&resource.Id, &resource.UserId, &resource.Content); err != nil {
		return nil, err
	}

	return resource, nil
}
func (db *AppDatabase) GetResourceById(id m.ResourceId) (*m.Resource, error) {
	query := "SELECT * FROM resources WHERE " + ResourceIdColumn + " = ?"
	stmt, err := db.connection.Prepare(query)

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

func (db *AppDatabase) DeleteResource(id m.ResourceId) error {
	q := `DELETE FROM resources WHERE ` + ResourceIdColumn + ` = ?`
	_, err := db.connection.Exec(q, id)
	return err
}

func (db *AppDatabase) UpdateResource(r m.Resource) error {
	q := `UPDATE resources SET ` + ResourceContentColumn + ` = ?  WHERE ` + ResourceIdColumn + ` = ?`
	_, err := db.connection.Exec(q, r.Content, r.Id)
	return err
}
