package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

type User struct {
	ID       int
	Email    string
	Password string
	Name     string
}

var userStorage []User
var authenticatedUser *User

func main() {
	fmt.Println("TODO start")
	command := flag.String("command", "no command", "command to run")
	flag.Parse()

	for {
		runCommand(*command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please enter the next command")
		scanner.Scan()
		*command = scanner.Text()
	}

}

func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var name, dueDate, category string

	fmt.Println("enter task title")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("enter due date")
	scanner.Scan()
	dueDate = scanner.Text()

	fmt.Println("enter task category")
	scanner.Scan()
	category = scanner.Text()

	fmt.Println("task: ", name, category, dueDate)
}
func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var catColor, catName string

	fmt.Println("enter category name")
	scanner.Scan()
	catName = scanner.Text()

	fmt.Println("enter category color")
	scanner.Scan()
	catColor = scanner.Text()

	fmt.Println("category: ", catName, catColor)
}
func userLogin() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password string

	fmt.Println("enter user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("enter user password")
	scanner.Scan()
	password = scanner.Text()

	for _, user := range userStorage {
		if strings.EqualFold(email, user.Email) && password == user.Password {
			fmt.Println("You are logged in!")
			authenticatedUser = &user

			break
		}

	}
	if authenticatedUser == nil {
		fmt.Println("Incorrect email or password")

		return
	}

	fmt.Println("user login: ", email, password)
}
func userRegister() {
	scanner := bufio.NewScanner(os.Stdin)
	var email, password, name string

	fmt.Println("enter your name")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("enter user email")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("enter user password")
	scanner.Scan()
	password = scanner.Text()

	id := len(userStorage) + 1

	fmt.Println("user: ", id, email, password, name)

	user := User{
		ID:       id,
		Email:    email,
		Password: password,
		Name:     name,
	}
	userStorage = append(userStorage, user)
}
func runCommand(cmd string) {
	if cmd != "user-register" && cmd != "exit" && authenticatedUser == nil {
		userLogin()
	}
	switch cmd {
	case "user-register":
		userRegister()
	case "user-login":
		userLogin()
	case "create-task":
		createTask()
	case "create-category":
		createCategory()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command not valid: ", cmd)
	}
}
