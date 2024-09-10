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

func (db *AppDatabase) UserExists(username string) bool {
	q := `SELECT COUNT(name) FROM users WHERE name = ?`
	row := db.connection.QueryRow(q, username)
	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 1
}

func (db *AppDatabase) GetUsers() []*m.User {
	q := "SELECT * FROM " + UsersTableName
	row, err := db.connection.Query(q)
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

func (db *AppDatabase) GetUserById(id int) (*m.User, error) {
	query := "SELECT * FROM " + UsersTableName + " WHERE " + UserIdColumn + " = ?"
	stmt, err := db.connection.Prepare(query)

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
func (db *AppDatabase) GetUserByName(name string) (*m.User, error) {
	query := "SELECT * FROM users WHERE name = ?"
	stmt, err := db.connection.Prepare(query)

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

func (db *AppDatabase) MakeAdmin(id m.UserId) error {
	q := "UPDATE " + UsersTableName + " SET " + UserIsAdminColumn + " = 1 WHERE " + UserIdColumn + " = ?"
	_, err := db.connection.Exec(q, id)
	return err
}
func (db *AppDatabase) CreateUser(name string, password string) (*m.User, error) {
	if db.UserExists(name) {
		return nil, errors.New("This username is occupied")
	}
	user := m.User{
		Name:     name,
		Password: password,
		IsAdmin:  false,
	}
	err := db.insertUser(user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (db *AppDatabase) insertUser(user m.User) error {
	insertUsersSQL := `INSERT INTO users(name, password, is_admin) VALUES (?, ?, ?)`
	statement, err := db.connection.Prepare(insertUsersSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		return err
	}
	_, err = statement.Exec(user.Name, user.Password, false)
	if err != nil {
		return err
	}
	return nil
}

// First to refactor
// TODO: refactor this
func DeleteUser(db *sql.DB, username *string, target string, new_username string) {
	log.Println("Deleting all resources for user ", target)
	q := `DELETE FROM resources WHERE id_user = ?`
	_, err := db.Exec(q, GetUserId(db, target))
	if err != nil {
		panic(err)
	}
	log.Println("Deleting user ", target)
	q = `DELETE FROM users WHERE name = ?`
	_, err = db.Exec(q, target)
	if err != nil {
		panic(err)
	}
	if target == *username {
		// RegisterUser(db, new_username)
		username = &new_username
	}
}
