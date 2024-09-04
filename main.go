package main

import (
	"bufio"
	database "cdecode/storage"
	"fmt"
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
	fmt.Println("\n!___Incorrect input___!")
}

func SpecifyContext() {
	fmt.Println("Please specify context!")
}

func AskForInput(reader *bufio.Reader, context string) string {
	fmt.Print(context)
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.TrimSpace(input)
	return input
}

func main() {
	//debug.DebugInit(true)
	// reader := bufio.NewReader(os.Stdin)
	// input, _ := reader.ReadString('\n')
	// input = strings.TrimSpace(input)
	// res := strings.Split(input, " ")
	// fmt.Print(res[0])
	db := database.InitDB()
	reader := bufio.NewReader(os.Stdin)
	current_user := AskForInput(reader, "Enter your desired username: ")
	database.RegisterUser(db, current_user)
	Help()
	for {
		input := AskForInput(reader, "Enter desired operation: ")
		options := strings.Split(input, " ")
		if options[0] == "help" {
			Help()
		} else if options[0] == "create" {
			if len(options) == 1 {
				SpecifyContext()
			} else if options[1] == "resource" {
				content := AskForInput(reader, "Enter content for your desired resource: ")
				name := AskForInput(reader, "Enter a name for this resource: ")
				database.CreateResorce(db, current_user, name, content)
			} else if options[1] == "alias" {
				resource_name := AskForInput(reader, "Enter the name of the resource you want to create an alias for: ")
				alias := AskForInput(reader, "Enter an alias for this resource: ")
				database.CreateAlias(db, current_user, resource_name, alias)
			} else if options[1] == "user" {
				if database.UserIsAdmin(db, current_user) {
					name := AskForInput(reader, "Enter a username for a user you want to create: ")
					database.CreateUser(db, current_user, name)
				} else {
					InpError()
				}
			} else {
				InpError()
			}
		} else if options[0] == "alias" {
			if len(options) == 1 {
				SpecifyContext()
			} else if options[1] == "connect" {
				resource_name := AskForInput(reader, "Enter the name of the resource you want to create an alias for: ")
				alias := AskForInput(reader, "Enter an alias you want to connect to this resource: ")
				database.AliasConnect(db, current_user, resource_name, alias)
			} else if options[1] == "disconnect" {
				alias := AskForInput(reader, "Enter an alias you want to disconnect from a resource: ")
				database.AliasDisconnect(db, current_user, alias)
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
				aliases := database.GetAliases(db, current_user)

				fmt.Println("\tid_a | id_u | id_r | name")
				for _, a := range aliases {
					fmt.Println("\t", a.Id, " | ", a.CreatorId, " | ", a.ResourceId, " | ", a.Name)
				}

			} else {
				InpError()
			}
		} else if options[0] == "delete" {
			if len(options) == 1 {
				SpecifyContext()
			} else if options[1] == "resources" {
				target := current_user
				if database.UserIsAdmin(db, current_user) {
					target = AskForInput(reader, "Enter a username, whose resource you want to delete: ")
				}
				resource_name := AskForInput(reader, "Enter a name of the resource you want to delete: ")
				database.DeleteResource(db, current_user, target, resource_name)
			} else if options[1] == "users" {
				target := current_user
				new_username := ""
				if database.UserIsAdmin(db, current_user) {
					target = AskForInput(reader, "Enter a username you want to delete: ")
				}
				if target == current_user {
					new_username = AskForInput(reader, "Enter a your new account name: ")
				}
				database.DeleteUser(db, &current_user, target, new_username)
			} else {
				InpError()
			}
		} else if options[0] == "read" {
			if len(options) == 1 {
				SpecifyContext()
			} else if options[1] == "resource" {
				target := current_user
				if database.UserIsAdmin(db, current_user) {
					target = AskForInput(reader, "Enter a username, whose resource you want to read: ")
				}
				resource_name := AskForInput(reader, "Enter a name of the resource you want to read: ")
				database.ReadResource(db, current_user, target, resource_name)
			} else if options[1] == "alias" {
				alias := AskForInput(reader, "Enter an alias to read a resource assigned to it: ")
				content := database.ReadContentByAlias(db, alias)
				fmt.Printf("Content:\n\n'%s'\n", content)
			} else {
				InpError()
			}
		} else if options[0] == "change" {
			if len(options) == 1 {
				SpecifyContext()
			} else if options[1] != "" {
				database.ChangeUser(db, &current_user, options[1])
			} else {
				InpError()
			}
		} else {
			InpError()
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
