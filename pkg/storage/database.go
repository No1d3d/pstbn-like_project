package storage

import (
	m "cdecode/pkg/models"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

var DbFile = "./pkg/storage/mydb.db"

func InitDB() *sql.DB {

	db, _ := sql.Open("sqlite3", DbFile)

	createTable(db, createUsersTableSQL)

	createTable(db, createResourcesTableSQL)

	createTable(db, createAliasTableSQL)

	return db
}

func createTable(db *sql.DB, query string) {
	log.Println("Creating database table...")
	log.Printf("Create table query:\n'%s'\n", query)

	statement, err := db.Prepare(query) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}

	_, err = statement.Exec() // Execute SQL Statements
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Table created")
}

func MakeAdmin(db *sql.DB, username string) {
	q := `UPDATE users SET is_admin = 1 WHERE id_user = ?;`
	_, err := db.Exec(q, GetUserId(db, username))
	if err != nil {
		panic(err)
	}
	fmt.Println("User <", username, "> now has admin rights!")
}

func UserExists(db *sql.DB, username string) bool {
	var count int
	q := `SELECT COUNT(name) FROM users WHERE name = ?`
	row := db.QueryRow(q, username)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 1
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

// insert new entry in users table
func InsertUser(db *sql.DB, user m.User) error {
	fmt.Println("Inserting user record...")
	insertUsersSQL := `INSERT INTO users(name, password, is_admin) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertUsersSQL) // Prepare statement.
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

func AdminCreateNewUser(db *sql.DB, username string) {
	if UserExists(db, username) {
		fmt.Println("This user already exists")
	} else {
		fmt.Println("Inserting users record ...")
		insertUsersSQL := `INSERT INTO users(name) VALUES (?)`
		statement, err := db.Prepare(insertUsersSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(username)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("New user created: ", username)
	}
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

// region authorisation

// Авторизация:
// пользователь вводит имя
// имя передается в БД и проверяется на наличие в ней
// Если имя существует - логин
// если нет - Авторизация

func RegisterUser(db *sql.DB, username string) {
	fmt.Print("\n---------------\nInitiating registration...\n\n")
	if UserExists(db, username) {
		fmt.Print("\nWelcome back, ", username, "!\n")
	} else {
		user := m.User{
			Name:    username,
			IsAdmin: false,
		}
		InsertUser(db, user)
		fmt.Print("\nNice to meet you, ", username, "!\n")
	}
	fmt.Print("\n\nAuthorisation complete!\n---------------\n\n")
	// }
}

// region create

func CreateUser(db *sql.DB, name string, password string) (*m.User, error) {
	if UserExists(db, name) {
		return nil, errors.New("This username is occupied")
	}
	user := m.User{
		Name:     name,
		Password: password,
		IsAdmin:  false,
	}
	err := InsertUser(db, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// region alias

func AliasConnect(db *sql.DB, username string, resource string, alias string) {
	if checkIfAliasHasNoDupes(db, alias) {
		fmt.Println("There is no such alias...")
		return
	}
	if CheckIfResourceHasNoDupNames(db, resource) {
		fmt.Println("There is no such resource...")
		return
	}
	q := `UPDATE alias SET id_resource = ? WHERE name = ?;`
	_, err := db.Exec(q, GetResourceId(db, username, resource), alias)
	if err != nil {
		panic(err)
	}
	fmt.Println("Alias <", alias, "> was connected to resource <", resource, ">")
}

func AliasDisconnect(db *sql.DB, username string, alias string) {
	if checkIfAliasHasNoDupes(db, alias) {
		fmt.Println("There is no such alias...")
		return
	}
	q := `UPDATE alias SET id_resource = 0 WHERE name = ?;`
	_, err := db.Exec(q, alias)
	if err != nil {
		panic(err)
	}
	fmt.Println("Alias <", alias, "> was disconnected from a resource!")
}

// region show

func ShowUsers(db *sql.DB, username string) {
	var row *sql.Rows
	var err error
	var q string
	if UserIsAdmin(db, username) {
		q = `SELECT * FROM users`
		row, err = db.Query(q)
	} else {
		q = `SELECT * FROM users WHERE name = ?`
		row, err = db.Query(q, username)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var id_user int
	var name string
	fmt.Println("    id    |    username")
	for row.Next() {
		if scan_err := row.Scan(&id_user, &name); err != nil {
			log.Fatal(scan_err)
		}
		fmt.Printf("  %4d    |   %4s\n", id_user, name)
	}
}

func ShowResources(db *sql.DB, username string) {
	var row *sql.Rows
	var err error
	var q string
	if UserIsAdmin(db, username) {
		q = `SELECT * FROM resources`
		row, err = db.Query(q)
	} else {
		q = `SELECT * FROM resources WHERE id_user = ?`
		row, err = db.Query(q, GetUserId(db, username))
	}
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var id_resource int
	var id_user int
	var name string
	var content string
	fmt.Println("\tid | id_u | name | content")
	for row.Next() {
		if scan_err := row.Scan(&id_resource, &id_user, &name, &content); err != nil {
			log.Fatal(scan_err)
		}
		fmt.Println("\t", id_resource, " | ", id_user, " | ", name, " | ", content)
	}
}

// region delete

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
		RegisterUser(db, new_username)
		username = &new_username
	}
}

// region read

func ReadContentByResource(db *sql.DB, username string, resourceName string) string {
	var row *sql.Rows
	q := `SELECT content FROM resources WHERE id_user = ? AND name = ?`
	row, err := db.Query(q, GetUserId(db, username), resourceName)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		log.Printf("No resource was found with name '%s'", resourceName)
		return ""
	}
	var content string

	if scan_err := row.Scan(&content); err != nil {
		log.Fatal(scan_err)
	}
	log.Printf("Content: %s", content)
	return content
}

func ReadResource(db *sql.DB, target string, resource_name string) ([]m.Resource, int) {
	var row *sql.Rows
	var err error
	counter := 0
	q := `SELECT content FROM resources WHERE id_user = ? AND name = ?`
	row, err = db.Query(q, GetUserId(db, target), resource_name)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	resources := []m.Resource{}
	for row.Next() {
		resource, scan_err := getResourceFromRow(row)
		if scan_err != nil {
			log.Println(scan_err)
		}
		resources = append(resources, *resource)
		counter++
	}
	return resources, counter
}

//region change

func ChangeUser(db *sql.DB, username *string, new_username string) {
	if !UserExists(db, new_username) {
		fmt.Print("\nNo such user <", new_username, ">\nRedirecting to user registration\n\n")
		RegisterUser(db, new_username)
	}
	username = &new_username
}
