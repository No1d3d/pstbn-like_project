package main

import (
	"bufio"
	database "cdecode/storage"
	"fmt"
	"log"
	"os"
	"strings"
)

func Help() {
	fmt.Println("\nList of commands:")
	fmt.Println("help\n\t- put out a list of commands")
	fmt.Println("create <option>\n\t- resource - creates a resource\n\t- alias - creates an alias")
	fmt.Println("alias <option>\n\t- connect - connect an alias to a resource\n\t- disconnect - disconnect an alias from resource")
	fmt.Println("show <option>\n\t- users - shows users information\n\t- resources - shows resources information\n\t- alias - shows alias information")
	fmt.Println("delete <option>\n\t- resource - deletes a resource\n\t- user - deletes a user")
	fmt.Println("read <option>\n\t- resource - read a resource by it's name\n\t- alias - resource by it's alias")
	fmt.Println("change <username>\n\t- change to another account by username.\n\tIf there is no account by that username - redirects to register a new account")
	fmt.Println("...")
}

func InpError() {
	fmt.Println("-----------------------\n!___Incorrect input___!\n-----------------------")
}

func SpecifyContext() {
	fmt.Println("Please specify context!")
}

func main() {
	//debug.DebugInit(true)
	// reader := bufio.NewReader(os.Stdin)
	// input, _ := reader.ReadString('\n')
	// input = strings.TrimSpace(input)
	// res := strings.Split(input, " ")
	// fmt.Print(res[0])
	db := database.InitDB()
	var current_user string
	reader := bufio.NewReader(os.Stdin)
	current_user = database.RegisterUser(db, reader)
	Help()
	for {
		fmt.Print("Enter desired operation: ")
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		options := strings.Split(input, " ")
		if err != nil {
			log.Fatal(err)
		} else {
			if options[0] == "help" {
				Help()
			} else if options[0] == "create" {
				if len(options) == 1 {
					SpecifyContext()
				} else if options[1] == "resource" {
					database.CreateResorce(db, current_user, reader)
				} else if options[1] == "alias" {
					database.CreateAlias(db, current_user, reader)
				} else if current_user == "admin" && options[1] == "user" {
					database.CreateUser(db, reader)
				} else {
					InpError()
				}
			} else if options[0] == "alias" {
				if len(options) == 1 {
					SpecifyContext()
				} else if options[1] == "connect" {
					database.AliasConnect(db, current_user, reader)
				} else if options[1] == "disconnect" {
					database.AliasDisconnect(db, current_user, reader)
				} else {
					InpError()
				}
			} else if options[0] == "show" {
				if len(options) == 1 {
					SpecifyContext()
				} else if options[1] == "users" {
					database.ShowUsers(db, current_user)
				} else if options[1] == "resources" {
					database.ShowResources(db, current_user)
				} else if options[1] == "alias" {
					database.ShowAlias(db, current_user)
				} else {
					InpError()
				}
			} else if options[0] == "delete" {
				if len(options) == 1 {
					SpecifyContext()
				} else if options[1] == "resources" {
					database.DeleteResource(db, current_user, reader)
				} else if options[1] == "users" {
					current_user = database.DeleteUser(db, current_user, reader)
				} else {
					InpError()
				}
			} else if options[0] == "read" {
				if len(options) == 1 {
					SpecifyContext()
				} else if options[1] == "resource" {
					database.ReadResource(db, current_user, reader)
				} else if options[1] == "alias" {
					database.ReadAlias(db, reader)
				} else {
					InpError()
				}
			} else if options[0] == "change" {
				if len(options) == 1 {
					SpecifyContext()
				} else if options[1] != "" {
					current_user = database.ChangeUser(db, options[1], reader)
				} else {
					InpError()
				}
			} else {
				InpError()
			}

		}
	}
}

// type fileEntry struct {
// 	name    string
// 	content string
// }

// // Файлы будут храниться в мапе в виде ключа/значения алиас/структура_файл
// // структура файл: name - string, content - string

// func ReadFile(r *bufio.Reader, container map[string]fileEntry) {
// 	fmt.Print("Enter file alias: ")
// 	input, _ := r.ReadString('\n')
// 	input = strings.TrimSpace(input)
// 	if val, ok := container[input]; ok {
// 		fmt.Print("name: ", val.name, "\n")
// 		fmt.Print("content: ", val.content, "\n")
// 	} else {
// 		fmt.Print("ERROR: NO SUCH ALIAS: ", input)
// 	}
// }

// func createEntry(r *bufio.Reader, container map[string]fileEntry) {
// 	var e fileEntry
// 	fmt.Print("Enter file name: ")
// 	e.name = CreateFileName(r)
// 	fmt.Print("Enter file content: ")
// 	e.content = CreateFileContent(r)
// 	fmt.Print("Enter file alias: ")
// 	a := CreateFileAlias(r)
// 	container[a] = e
// }

// func CreateFileAlias(r *bufio.Reader) string {
// 	input, _ := r.ReadString('\n')
// 	input = strings.TrimSpace(input)
// 	return input
// }

// func CreateFileName(r *bufio.Reader) string {
// 	input, _ := r.ReadString('\n')
// 	input = strings.TrimSpace(input)
// 	return input
// }

// func CreateFileContent(r *bufio.Reader) string {
// 	input, _ := r.ReadString('\n')
// 	input = strings.TrimSpace(input)
// 	return input
// }

// func ShowFiles(container map[string]fileEntry) {
// 	for key, value := range container {
// 		fmt.Print("-----------\n")
// 		fmt.Print("alias: ", key, "\n")
// 		fmt.Print("\tname: ", value.name, "\n")
// 		fmt.Print("\tcontent: ", value.content, "\n")
// 		fmt.Print("-----------\n\n")
// 	}
// }
