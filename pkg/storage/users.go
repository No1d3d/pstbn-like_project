package storage

import (
	m "cdecode/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

const (
	UsersTableName      = "users"
	UserIdColumn        = "id_user"
	UserNameColumn      = "name"
	UserPasswordColumn  = "password"
	UserIsAdminColumn   = "is_admin"
	createUsersTableSQL = `CREATE TABLE IF NOT EXISTS ` + UsersTableName + ` (
		"` + UserIdColumn + `" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"` + UserNameColumn + `" TEXT UNIQUE NOT NULL,
    "` + UserPasswordColumn + `" TEXT NOT NULL,
		"` + UserIsAdminColumn + `" BOOLEAN DEFAULT FALSE
	  );`
)

func getUserFromRow(row *sql.Rows) (*m.User, error) {
	user := &m.User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.IsAdmin); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsers(db *sql.DB) []*m.User {
	q := "SELECT * FROM " + UsersTableName
	row, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	users := []*m.User{}
	for row.Next() {
		user, err := getUserFromRow(row)
		if err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}
	return users
}

func GetUserById(db *sql.DB, id int) (*m.User, error) {
	query := "SELECT * FROM " + UsersTableName + " WHERE " + UserIdColumn + " = ?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(id)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		user, err := getUserFromRow(row)
		row.Close()
		return user, err
	}
	return nil, errors.New(fmt.Sprintf("No such user with id %d", id))
}
func GetUserByName(db *sql.DB, name string) (*m.User, error) {
	query := "SELECT * FROM users WHERE name = ?"
	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	row, err := stmt.Query(name)

	if err != nil {
		return nil, err
	}

	for row.Next() {
		user, err := getUserFromRow(row)
		return user, err
	}
	return nil, errors.New(fmt.Sprintf("No such user with name '%s'", name))
}

func MakeAdmin(db *sql.DB, id m.UserId) error {
	q := "UPDATE " + UsersTableName + " SET " + UserIsAdminColumn + " = 1 WHERE " + UserIdColumn + " = ?"
	_, err := db.Exec(q, id)
	return err
}
