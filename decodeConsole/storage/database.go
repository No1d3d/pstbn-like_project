package database

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3" // Import go-sqlite3 library
)

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

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
		"name" TEXT		
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

func UserExists(db *sql.DB, username string) bool {
	var row *sql.Rows
	var err error
	var q string
	if username == "admin" {
		return true
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
	counter := 0
	for row.Next() {
		if scan_err := row.Scan(&id_user, &name); err != nil {
			log.Fatal(scan_err)
		}
		counter++
	}
	return counter == 1
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

func AdminCreateNewUser(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("Enter a username you want to create: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if UserExists(db, input) {
		fmt.Println("This user already exists")
	} else {
		fmt.Println("Inserting users record ...")
		insertUsersSQL := `INSERT INTO users(name) VALUES (?)`
		statement, err := db.Prepare(insertUsersSQL) // Prepare statement.
		// This is good to avoid SQL injections
		if err != nil {
			log.Fatalln(err.Error())
		}
		_, err = statement.Exec(input)
		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Println("New user created: ", input)
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

// check no duplicated name of alias
func CheckIfAliasHasNoDupes(db *sql.DB, alias string) bool {
	var count int
	q := `SELECT COUNT(name) FROM alias WHERE name = ?`
	row := db.QueryRow(q, alias)
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count == 0
}

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

// region authorisation

func RegisterUser(db *sql.DB, r *bufio.Reader) string {
	fmt.Print("\n---------------\nInitiating registration...\n\n")
	fmt.Print("Enter your preffered username: ")
	username, _ := r.ReadString('\n')
	username = strings.TrimSpace(username)
	if username == "admin" {
		fmt.Print("\nWelcome back, admin!\n")
		return "admin"
	}
	if UserExists(db, username) {
		fmt.Print("\nWelcome back, ", username, "!\n")
	} else {
		InsertUsers(db, username)
		fmt.Print("\nNice to meet you, ", username, "!\n")
	}
	fmt.Print("\n\nAuthorisation complete!\n---------------\n\n")
	return username
}

// region create

func CreateResorce(db *sql.DB, username string, reader *bufio.Reader) {
	if username == "admin" {
		fmt.Println("Why would admin want to create new resources? :)")
		return
	}
	var name, content string
	fmt.Print("Type in the content of your resource. Press Enter to finish typing:\n")
	content, _ = reader.ReadString('\n')
	content = strings.TrimSpace(content)
	fmt.Print("\nType in the name for this resource: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)
	for !CheckIfResourceHasNoDupNames(db, username, name) {
		fmt.Println("You already have this name for a resource! Choose another one!")
		fmt.Print("\nType in the name for this resource: ")
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
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

func CreateAlias(db *sql.DB, username string, reader *bufio.Reader) {
	if username == "admin" {
		fmt.Println("Why would admin want to create new aliases? :)")
		return
	}
	var resource_name, alias string
	fmt.Print("Type in the name of the resource you want to create an alias for.\n")
	resource_name, _ = reader.ReadString('\n')
	resource_name = strings.TrimSpace(resource_name)
	fmt.Print("\n\nType in the alias for this resource: ")
	alias, _ = reader.ReadString('\n')
	alias = strings.TrimSpace(alias)
	for !CheckIfAliasHasNoDupes(db, alias) {
		fmt.Println("This Alias already exists, please choose another one.")
		fmt.Print("\nType in the alias for this resource: ")
		alias, _ = reader.ReadString('\n')
		alias = strings.TrimSpace(alias)
	}
	fmt.Println("Adding new alias...")
	InsertAlias(db, GetUserId(db, username), GetResourceId(db, username, resource_name), alias)
	fmt.Println("New alias added!")
}

func InsertAlias(db *sql.DB, id_user int, id_resource int, name string) {
	fmt.Println("Inserting alias record ...")
	insertAliasSQL := `INSERT INTO alias(id_user, id_resource, name) VALUES (?, ?, ?)`
	statement, err := db.Prepare(insertAliasSQL) // Prepare statement.
	// This is good to avoid SQL injections
	if err != nil {
		log.Fatalln(err.Error())
	}
	_, err = statement.Exec(id_user, id_resource, name)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func CreateUser(db *sql.DB, reader *bufio.Reader) {
	fmt.Print("Type in the username you want to add: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if UserExists(db, name) {
		fmt.Println("This username is occupied.")
	} else {
		InsertUsers(db, name)
	}
}

// region alias

func AliasConnect(db *sql.DB, username string, reader *bufio.Reader) {
	fmt.Print("Type in the alias you want to connect to a resource: ")
	alias, _ := reader.ReadString('\n')
	alias = strings.TrimSpace(alias)
	if CheckIfAliasHasNoDupes(db, alias) {
		fmt.Println("There is no such alias...")
		return
	}
	fmt.Print("Type in the name of the resource you want to connect an alias to: ")
	resource, _ := reader.ReadString('\n')
	resource = strings.TrimSpace(resource)
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

func AliasDisconnect(db *sql.DB, username string, reader *bufio.Reader) {
	fmt.Print("Type in the alias you want to disconnect from a resource: ")
	alias, _ := reader.ReadString('\n')
	alias = strings.TrimSpace(alias)
	if CheckIfAliasHasNoDupes(db, alias) {
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

func ShowUsers(db *sql.DB, cur_user string) {
	var row *sql.Rows
	var err error
	var q string
	if cur_user == "admin" {
		q = `SELECT * FROM users`
		row, err = db.Query(q)
	} else {
		q = `SELECT * FROM users WHERE name = ?`
		row, err = db.Query(q, cur_user)
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

func ShowResources(db *sql.DB, cur_user string) {
	var row *sql.Rows
	var err error
	var q string
	if cur_user == "admin" {
		q = `SELECT * FROM resources`
		row, err = db.Query(q)
	} else {
		q = `SELECT * FROM resources WHERE id_user = ?`
		row, err = db.Query(q, GetUserId(db, cur_user))
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

func ShowAlias(db *sql.DB, cur_user string) {
	var row *sql.Rows
	var err error
	var q string
	if cur_user == "admin" {
		q = `SELECT * FROM alias`
		row, err = db.Query(q)
	} else {
		q = `SELECT * FROM alias WHERE id_user = ?`
		row, err = db.Query(q, GetUserId(db, cur_user))
	}
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()
	var id_alias int
	var id_user int
	var id_resource int
	var name string
	fmt.Println("\tid_a | id_u | id_r | name")
	for row.Next() {
		if scan_err := row.Scan(&id_alias, &id_user, &id_resource, &name); err != nil {
			log.Fatal(scan_err)
		}
		fmt.Println("\t", id_alias, " | ", id_user, " | ", id_resource, " | ", name)
	}
}

// region delete

func DeleteResource(db *sql.DB, username string, reader *bufio.Reader) {
	if username == "admin" {
		fmt.Print("Enter a username for which you want to delete a resource: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}
	fmt.Print("Enter the name of a resource you want to delete: ")
	resource, _ := reader.ReadString('\n')
	resource = strings.TrimSpace(resource)
	log.Println("Deleting resource ", resource)
	q := `DELETE FROM resources WHERE id_user = ? and name = ?`
	_, err := db.Exec(q, GetUserId(db, username), resource)
	if err != nil {
		panic(err)
	}

}

func DeleteUser(db *sql.DB, username string, reader *bufio.Reader) string {
	flag := 0
	if username == "admin" {
		flag = 1
		fmt.Print("Enter a username for which you want to delete entries for: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}
	log.Println("Deleting all resources for user ", username)
	q := `DELETE FROM resources WHERE id_user = ?`
	_, err := db.Exec(q, GetUserId(db, username))
	if err != nil {
		panic(err)
	}
	log.Println("Deleting user ", username)
	q = `DELETE FROM users WHERE name = ?`
	_, err = db.Exec(q, username)
	if err != nil {
		panic(err)
	}
	if flag == 1 {
		username = "admin"
	} else {
		username = RegisterUser(db, reader)
	}
	return username

}

// region read

func ReadResource(db *sql.DB, username string, reader *bufio.Reader) {
	if username == "admin" {
		fmt.Print("Enter a username for which you want to read a resource: ")
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}
	fmt.Print("Enter a name of the resource you want to read: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	var row *sql.Rows
	var err error
	q := `SELECT content FROM resources WHERE id_user = ? AND name = ?`
	row, err = db.Query(q, username, name)
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
		fmt.Println("No resource with name <", name, "> Error...")
	}
}

func ReadAlias(db *sql.DB, reader *bufio.Reader) {
	var row *sql.Rows
	var err error
	fmt.Print("Enter an alias for the resource you want to read: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	q := `SELECT content FROM resources WHERE id_resource = (SELECT id_resource FROM alias WHERE name = ?)`
	row, err = db.Query(q, name)
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
		fmt.Println("No resource was found with alias <", name, "> Error...")
	}
}

//region change

func ChangeUser(db *sql.DB, username string, reader *bufio.Reader) string {
	if username == "admin" {
		fmt.Print("\nUsername changed to <admin>\n")
		return "admin"
	} else {
		if UserExists(db, username) {
			fmt.Println("\nUsername changed to <", username, ">")
			return username
		} else {
			fmt.Print("\nNo such user <", username, ">\nRedirecting to user registration\n\n")
			username = RegisterUser(db, reader)
			fmt.Print("\nUsername changed to <", username, ">\n")
			return username
		}
	}

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

// endregion
