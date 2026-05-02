package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SoroushBeigi/todo-cli-go/constant"
	"github.com/SoroushBeigi/todo-cli-go/contract"
	"github.com/SoroushBeigi/todo-cli-go/entity"
	"github.com/SoroushBeigi/todo-cli-go/filestore"
)

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

var (
	userStorage       []entity.User
	taskStorage       []Task
	categoryStorage   []Category
	authenticatedUser *entity.User
	serializationMode string
)

const userStoragePath = "users.txt"

func main() {

	fmt.Println("TODO start")
	command := flag.String("command", "no command", "command to run")
	sm := flag.String("serialization", constant.ManualMode, "how application writes data to file")
	flag.Parse()

	switch *sm {
	case constant.ManualMode:
		serializationMode = constant.ManualMode
	default:
		serializationMode = constant.JsonMode
	}

	var userFileStore = filestore.NewFileStore(userStoragePath, serializationMode)

	users := userFileStore.Load()
	userStorage = append(userStorage, users...)

	for {
		runCommand(userFileStore, *command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please enter the next command")
		scanner.Scan()
		*command = scanner.Text()
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
		if strings.EqualFold(email, user.Email) && hashPassword(password) == user.Password {
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

func userRegister(store contract.UserWriteStore) {
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

	u := entity.User{
		ID:       id,
		Email:    email,
		Password: hashPassword(password),
		Name:     name,
	}
	userStorage = append(userStorage, u)

	store.Save(u)

}

func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))

	return hex.EncodeToString(hash[:])

}

func runCommand(userFileStore filestore.FileStore, cmd string) {
	if cmd != "user-register" && cmd != "exit" && authenticatedUser == nil {
		userLogin()
		if authenticatedUser == nil {
			return
		}
	}

	switch cmd {
	case "user-register":
		userRegister(&userFileStore)
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
