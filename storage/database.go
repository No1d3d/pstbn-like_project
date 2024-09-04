package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	return !errors.Is(error, os.ErrNotExist)
}

func InitDB() *sql.DB {
	// if db does not exist, creates one
	if !CheckFileExists("./storage/mydb.db") {
		fmt.Println("Creating mydb.db...")
		file, err := os.Create("./storage/mydb.db")
		if err != nil {
			log.Fatal(err.Error())
		}
		file.Close()
		fmt.Println("mydb.db created")
	}
	// opens db
	db, _ := sql.Open("sqlite3", "./storage/mydb.db")

	// create users table
	createUsersTableSQL := `CREATE TABLE IF NOT EXISTS users (
		"id_user" integer NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"name" TEXT UNIQUE,
		"is_admin" BOOLEAN
	  );` // SQL Statement for Create Table

	fmt.Println("Create users table...")
	statement, err := db.Prepare(createUsersTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	fmt.Println("users table created")

	// create resources table
	createResourcesTableSQL := `CREATE TABLE IF NOT EXISTS resources (
		"id_resource" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"id_user" integer NOT NULL,		
		"name" TEXT,
		"content" TEXT		
	  );` // SQL Statement for Create Table

	fmt.Println("Create resources table...")
	statement, err = db.Prepare(createResourcesTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	fmt.Println("resources table created")

	// create alias table
	createAliasTableSQL := `CREATE TABLE IF NOT EXISTS alias (
		"id_alias" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"id_user" integer NOT NULL,		
		"id_resource" integer NOT NULL,		
		"name" TEXT	
	  );` // SQL Statement for Create Table

	fmt.Println("Create alias table...")
	statement, err = db.Prepare(createAliasTableSQL) // Prepare SQL Statement
	if err != nil {
		log.Fatal(err.Error())
	}
	statement.Exec() // Execute SQL Statements
	fmt.Println("alias table created")
	return db
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
func InsertUsers(db *sql.DB, name string) {
	fmt.Println("Inserting users record ...")
	insertUsersSQL := `INSERT INTO users(name) VALUES (?)`
	statement, err := db.Prepare(insertUsersSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(name)
	if err != nil {
		log.Fatalln(err.Error())
	}
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

func GetResourceId(db *sql.DB, username, resource_name string) int {
	var row *sql.Rows
	var err error
	q := `SELECT id_resource FROM resources WHERE id_user = ? AND name = ?`
	row, err = db.Query(q, GetUserId(db, username), resource_name)
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
		InsertUsers(db, username)
		fmt.Print("\nNice to meet you, ", username, "!\n")
	}
	fmt.Print("\n\nAuthorisation complete!\n---------------\n\n")
}

// region create

func CreateResorce(db *sql.DB, username string, name string, content string) {
	if !CheckIfResourceHasNoDupNames(db, username, name) {
		return // This name is occupied
	}
	fmt.Println("Adding new resource...")
	InsertResources(db, GetUserId(db, username), name, content)
	fmt.Println("New resource added!")
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

func CreateUser(db *sql.DB, username string, name string) {
	if UserExists(db, name) {
		fmt.Println("This username is occupied.")
	} else {
		InsertUsers(db, name)
	}
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

func ReadResource(db *sql.DB, username string, target string, resource_name string) {
	var row *sql.Rows
	var err error
	q := `SELECT content FROM resources WHERE id_user = ? AND name = ?`
	row, err = db.Query(q, target, resource_name)
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var content string
	counter := 0
	for row.Next() {
		if scan_err := row.Scan(&content); err != nil {
			log.Fatal(scan_err)
		}
		fmt.Println("\t", content)
		counter++
	}
	if counter == 0 {
		fmt.Println("No resource with name <", resource_name, "> Error...")
	}
}

//region change

func ChangeUser(db *sql.DB, username *string, new_username string) {
	if !UserExists(db, new_username) {
		fmt.Print("\nNo such user <", new_username, ">\nRedirecting to user registration\n\n")
		RegisterUser(db, new_username)
	}
	username = &new_username
}

// endregion

// regio Delete
// func DeleteEntry(db *sql.DB, reader *bufio.Reader, username string) string {
// 	cur_user := username
// 	fmt.Println("What would you like to delete (users, resources, alias)?")
// 	fmt.Println("Your choice: ")
// 	choice, _ := reader.ReadString('\n')
// 	choice = strings.TrimSpace(choice)
// 	if choice == "alias" {
// 		fmt.Print("Enter name of alias to delete: ")
// 		inpt, _ := reader.ReadString('\n')
// 		inpt = strings.TrimSpace(inpt)
// 		fmt.Println("Deleting alias")
// 		DeleteAlias(db, inpt, username)
// 		fmt.Println("Alias deleted!")
// 	} else if choice == "resources" {
// 		fmt.Print("Enter name of resource to delete: ")
// 		inpt, _ := reader.ReadString('\n')
// 		inpt = strings.TrimSpace(inpt)
// 		fmt.Println("Deleting resource")
// 		DeleteResource(db, inpt, username)
// 		fmt.Println("Resource deleted")
// 	} else if choice == "users" {
// 		var inpt string
// 		if username == "admin" {
// 			fmt.Print("Enter name of user to delete: ")
// 			inpt, _ = reader.ReadString('\n')
// 			inpt = strings.TrimSpace(inpt)
// 		} else {
// 			inpt = username
// 		}
// 		fmt.Println("Deleting user |", inpt, "|")
// 		DeleteUser(db, inpt)
// 		fmt.Println("User deleted!")
// 		if inpt == username {
// 			fmt.Println("Redirecting to Authorization!")
// 			cur_user = RegisterUser(db, reader)
// 		}
// 	} else {
// 		fmt.Println("Error: no such option:\t", choice)
// 	}
// 	return cur_user
// }

// func DeleteAlias(db *sql.DB, alias string, username string) {
// 	log.Println("Deleting alias ", alias)
// 	q := `DELETE FROM alias WHERE id_user = ? AND name = ?`
// 	_, err := db.Exec(q, GetUserId(db, username), alias)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func ShowResource(db *sql.DB, reader *bufio.Reader) {
// 	fmt.Println("Enter an alias assigned for a resource: ")
// 	alias, _ := reader.ReadString('\n')
// 	alias = strings.TrimSpace(alias)
// 	fmt.Println("Looking for a resource...")
// 	DisplayAlias(db, alias)
// }

// func DisplayAlias(db *sql.DB, resource_alias string) {
// 	var row *sql.Rows
// 	var err error
// 	q := `SELECT content FROM resources WHERE id_resource = (SELECT id_resource FROM alias WHERE name = ?)`
// 	row, err = db.Query(q, resource_alias)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer row.Close()
// 	var content string
// 	counter := 0
// 	for row.Next() {
// 		if scan_err := row.Scan(&content); err != nil {
// 			log.Fatal(scan_err)
// 		}
// 		fmt.Println("\t", content)
// 		counter++
// 	}
// 	if counter == 0 {
// 		fmt.Println("No resource... Error...")
// 	}
// }

// os.Remove("sqlite-database.db") // I delete the file to avoid duplicated records.
//                                 // SQLite is a file based database.

// fmt.Println("Creating sqlite-database.db...")
// file, err := os.Create("sqlite-database.db") // Create SQLite file
// if err != nil {
// 	log.Fatal(err.Error())
// }
// file.Close()
// fmt.Println("sqlite-database.db created")

// sqliteDatabase, _ := sql.Open
// ("sqlite3", "./sqlite-database.db") // Open the created SQLite File
// defer sqliteDatabase.Close() // Defer Closing the database

// create alias
// func ShowResource(db *sql.DB, reader *bufio.Reader) {
// 	fmt.Println("Enter an alias assigned for a resource: ")
// 	alias, _ := reader.ReadString('\n')
// 	alias = strings.TrimSpace(alias)
// 	fmt.Println("Looking for a resource...")
// 	DisplayAlias(db, alias)
// }

// func DisplayAlias(db *sql.DB, resource_alias string) {
// 	var row *sql.Rows
// 	var err error
// 	q := `SELECT content FROM resources WHERE id_resource = (SELECT id_resource FROM alias WHERE name = ?)`
// 	row, err = db.Query(q, resource_alias)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer row.Close()
// 	var content string
// 	counter := 0
// 	for row.Next() {
// 		if scan_err := row.Scan(&content); err != nil {
// 			log.Fatal(scan_err)
// 		}
// 		fmt.Println("\t", content)
// 		counter++
// 	}
// 	if counter == 0 {
// 		fmt.Println("No resource... Error...")
// 	}
// }

// check if user exists
// func UserExists(db *sql.DB, username string) bool {
// 	sqlStmt := `SELECT username FROM userinfo WHERE username = ?`
// 	err := db.QueryRow(sqlStmt, username).Scan(&username)
// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			// a real error happened! you should change your function return
// 			// to "(bool, error)" and return "false, err" here
// 			fmt.Print(err)
// 		}

// 		return false
// 	}

// 	return true
// }

// func UserExists(db *sql.DB, username string) bool {
// 	row := db.QueryRow("select user_email from users where user_email= ?", username)
// 	temp := ""
// 	row.Scan(&temp)
// 	return temp != ""
// }

// endregion
