package storage

import (
	"database/sql"
	"log"
)

type UserId = int
type AliasId = int
type ResourceId = int

type Alias struct {
	Id         AliasId
	CreatorId  UserId
	ResourceId ResourceId
	Name       string
}

func GetAliases(db *sql.DB, username string) []Alias {
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
	aliases := []Alias{}
	for row.Next() {
		alias, scan_err := getAliasFromRow(row)
		if scan_err != nil {
			log.Println(scan_err)
		}
		aliases = append(aliases, *alias)
	}

	return aliases
}

func getAliasFromRow(row *sql.Rows) (*Alias, error) {
	alias := &Alias{}
	if err := row.Scan(&alias.Id, &alias.CreatorId, &alias.ResourceId, &alias.Name); err != nil {
		return nil, err
	}

	return alias, nil
}
