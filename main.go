package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Email    string
	Password string
	Name     string
}

type Task struct {
	ID         int
	UserID     int
	Title      string
	DueDate    string
	CategoryID int
	IsDone     bool
}

type Category struct {
	ID     int
	Title  string
	Color  string
	UserID int
}

var userStorage []User
var taskStorage []Task
var categoryStorage []Category
var authenticatedUser *User

const userStoragePath = "users.txt"

func main() {
	initStorage()
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

func initStorage() {
	data, err := os.ReadFile(userStoragePath)
	if err != nil {

		if os.IsNotExist(err) {
			return
		}
		fmt.Println("cannot read file: ", err)
		return
	}

	dataStr := string(data)
	users := strings.Split(dataStr, "\n")

	for _, uRow := range users {
		uRow = strings.TrimSpace(uRow)
		if uRow == "" {
			continue
		}

		userFields := strings.Split(uRow, ",")
		user := User{}
		for _, field := range userFields {
			field = strings.TrimSpace(field)
			content := strings.Split(field, ": ")
			if len(content) < 2 {

				continue
			}
			fieldName := content[0]
			fieldValue := content[1]

			switch fieldName {
			case "id":
				id, err := strconv.Atoi(fieldValue)
				if err != nil {
					fmt.Println("ERROR: ", err)

					continue
				}
				user.ID = id
				break
			case "name":
				user.Name = fieldValue
				break
			case "email":
				user.Email = fieldValue
				break
			case "password":
				user.Password = fieldValue
				break
			}

		}
		userStorage = append(userStorage, user)
		fmt.Printf("%+v\n", user)
	}

}

func createTask() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, dueDate string

	fmt.Println("enter task title")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("enter due date")
	scanner.Scan()
	dueDate = scanner.Text()

	fmt.Println("enter the category id")
	scanner.Scan()
	categoryId, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Printf("Category ID is not a valid integer: %v\n", err)

		return
	}

	found := false
	for _, c := range categoryStorage {
		if c.ID == categoryId && c.UserID == authenticatedUser.ID {
			found = true

			break
		}
	}

	if !found {
		fmt.Println("Category ID is not defined for User")

		return
	}

	if authenticatedUser != nil {
		t := Task{
			ID:         len(taskStorage) + 1,
			UserID:     authenticatedUser.ID,
			IsDone:     false,
			Title:      title,
			CategoryID: categoryId,
			DueDate:    dueDate,
		}

		taskStorage = append(taskStorage, t)
	}

}
func createCategory() {
	scanner := bufio.NewScanner(os.Stdin)
	var title, color string

	fmt.Println("enter category name")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("enter category color")
	scanner.Scan()
	color = scanner.Text()

	c := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}

	categoryStorage = append(categoryStorage, c)
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

	u := User{
		ID:       id,
		Email:    email,
		Password: password,
		Name:     name,
	}
	userStorage = append(userStorage, u)

	writeUserToFile(u)

}

func writeUserToFile(u User) {
	_, err := os.Stat(userStoragePath)

	file, err := os.OpenFile(userStoragePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cannot create or open file: ", err)

		return
	}

	defer file.Close()

	data := fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", u.ID, u.Name, u.Email, u.Password)

	_, err = file.Write([]byte(data))
	if err != nil {
		fmt.Println("Error writing file: ", err)
	}
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
	case "list-task":
		listTask()
	case "create-category":
		createCategory()
	case "exit":
		os.Exit(0)
	default:
		fmt.Println("command not valid: ", cmd)
	}
}

func listTask() {
	for i, task := range taskStorage {
		if task.UserID == authenticatedUser.ID {
			fmt.Printf("Task #%v: %+v\n", i, task)
		}
	}
}
