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

type Task struct {
	ID       int
	UserID   int
	Title    string
	DueDate  string
	Category string
	IsDone   bool
}

var userStorage []User
var taskStorage []Task
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
	var title, dueDate, category string

	fmt.Println("enter task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("enter due date")
	scanner.Scan()
	dueDate = scanner.Text()

	fmt.Println("enter task category")
	scanner.Scan()
	category = scanner.Text()

	if authenticatedUser != nil {
		task := Task{
			ID:       len(taskStorage) + 1,
			UserID:   authenticatedUser.ID,
			IsDone:   false,
			Title:    title,
			Category: category,
			DueDate:  dueDate,
		}

		taskStorage = append(taskStorage, task)
	}

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
		if authenticatedUser == nil {
			return
		}
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
