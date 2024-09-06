package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

type User struct {
	Id       UserId `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

const (
	HiddenPassword = "****"
)

func (u *User) HidePassword() {
	u.Password = HiddenPassword
}

type Resource struct {
	Id      ResourceId `json:"id"`
	UserId  UserId     `json:"userId"`
	Name    string     `json:"name"`
	Content string     `json:"value"`
}

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

var DbFile = "./storage/mydb.db"

func InitDB() *sql.DB {
	// if db does not exist, creates one
	if !CheckFileExists(DbFile) {
		log.Printf("Creating database file '%s'\n", DbFile)
		file, err := os.Create(DbFile)
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		log.Printf("Database file '%s' created", DbFile)
	}
	// opens db
	db, _ := sql.Open("sqlite3", DbFile)

	// create users table
	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id_user" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT UNIQUE NOT NULL,
    "password" TEXT NOT NULL,
		"is_admin" BOOLEAN DEFAULT FALSE
	  );` // SQL Statement for Create Table
	createTable(db, createUsersTableSQL)

	// create resources table
	createResourcesTableSQL := `CREATE TABLE IF NOT EXISTS resources (
		"id_resource" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"id_user" integer NOT NULL,		
		"name" TEXT,
		"content" TEXT		
	  );` // SQL Statement for Create Table
	createTable(db, createResourcesTableSQL)

	// create alias table
	createAliasTableSQL := `CREATE TABLE IF NOT EXISTS alias (
		"id_alias" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"id_user" integer NOT NULL,		
		"id_resource" integer NOT NULL,		
		"name" TEXT	
	  );` // SQL Statement for Create Table

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
func InsertUser(db *sql.DB, user User) error {
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
func CheckIfResourceHasNoDupNames(db *sql.DB, username string, resource_name string) bool {
	var count int
	q := `SELECT COUNT(name) FROM resources WHERE name = ? AND id_user = ?`
	row := db.QueryRow(q, resource_name, GetUserId(db, username))
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
		user := User{
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

func CreateResource(db *sql.DB, username string, name string, content string) {
	if !CheckIfResourceHasNoDupNames(db, username, name) {
		return // This name is occupied
	}
	log.Println("Adding new resource...")
	InsertResources(db, GetUserId(db, username), name, content)
	log.Println("New resource added!")
}

func InsertResources(db *sql.DB, id_user int, name string, content string) {
	fmt.Println("Inserting resources record ...")
	insertResourcesSQL := `INSERT INTO resources(id_user, name, content) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertResourcesSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id_user, name, content)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func CreateUser(db *sql.DB, name string, password string) (*User, error) {
	if UserExists(db, name) {
		return nil, errors.New("This username is occupied")
	}
	user := User{
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
	if CheckIfResourceHasNoDupNames(db, username, resource) {
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

func getUserFromRow(row *sql.Rows) (*User, error) {
	user := &User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.IsAdmin); err != nil {
		return nil, err
	}

	return user, nil
}

func GetUsers(db *sql.DB) []*User {
	q := `SELECT * FROM users`
	row, err := db.Query(q)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	users := []*User{}
	for row.Next() {
		user, err := getUserFromRow(row)
		if err != nil {
			log.Println(err)
		}
		users = append(users, user)
	}
	return users
}

func GetUserById(db *sql.DB, id int) (*User, error) {
	query := "SELECT * FROM users WHERE id_user = ?"
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
		return user, err
	}
	return nil, errors.New(fmt.Sprintf("No such user with id %d", id))
}

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

func getResourceFromRow(row *sql.Rows) (*Resource, error) {
	resource := &Resource{}
	if err := row.Scan(&resource.Id, &resource.UserId, &resource.Name, &resource.Content); err != nil {
		return nil, err
	}

	return resource, nil
}

func GetResources(db *sql.DB, username string) []Resource {
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
	resources := []Resource{}
	for row.Next() {
		resource, scan_err := getResourceFromRow(row)
		if scan_err != nil {
			log.Println(scan_err)
		}
		resources = append(resources, *resource)
	}
	return resources
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

func DeleteResource(db *sql.DB, username string, target string, resource string) {
	log.Println("Deleting resource ", resource)
	q := `DELETE FROM resources WHERE id_user = ? and name = ?`
	_, err := db.Exec(q, GetUserId(db, target), resource)
	if err != nil {
		panic(err)
	}

}

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

func ReadContentByResource(db *sql.DB, username string, resource_name string) string {
	var row *sql.Rows
	q := `SELECT content FROM resources WHERE id_user = ? AND name = ?`
	row, err := db.Query(q, GetUserId(db, username), resource_name)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	if !row.Next() {
		log.Printf("No resource was found with name '%s'", resource_name)
		return ""
	}
	var content string

	if scan_err := row.Scan(&content); err != nil {
		log.Fatal(scan_err)
	}
	log.Printf("Content: %s", content)
	return content
}

func ReadResource(db *sql.DB, target string, resource_name string) ([]Resource, int) {
	var row *sql.Rows
	var err error
	counter := 0
	q := `SELECT content FROM resources WHERE id_user = ? AND name = ?`
	row, err = db.Query(q, GetUserId(db, target), resource_name)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	resources := []Resource{}
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
