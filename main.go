package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
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

var (
	userStorage       []User
	taskStorage       []Task
	categoryStorage   []Category
	authenticatedUser *User
	serializationMode string
)

const (
	userStoragePath = "users.txt"
	ManualMode      = "manual"
	JsonMode        = "json"
)

func main() {

	fmt.Println("TODO start")
	command := flag.String("command", "no command", "command to run")
	sm := flag.String("serialization", ManualMode, "how application writes data to file")
	flag.Parse()

	switch *sm {
	case ManualMode:
		serializationMode = ManualMode
	default:
		serializationMode = JsonMode
	}

	loadUserFromStorage(&fileStore{filePath: userStoragePath}, serializationMode)

	for {
		runCommand(*command)
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please enter the next command")
		scanner.Scan()
		*command = scanner.Text()
	}

}
func loadUserFromStorage(store userReadStore, sm string) {
	users := store.Load(sm)
	userStorage = append(userStorage, users...)
}

func manualDecode(uRow string) (User, error) {
	uRow = strings.TrimSpace(uRow)
	if uRow == "" {
		return User{}, fmt.Errorf("Empty user")
	}

	userFields := strings.Split(uRow, ",")
	user := User{}
	for _, field := range userFields {
		field = strings.TrimSpace(field)
		content := strings.Split(field, ": ")
		if len(content) < 2 {

			return User{}, fmt.Errorf("content length <2")
		}
		fieldName := content[0]
		fieldValue := content[1]

		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				return User{}, fmt.Errorf("str convert error")
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
	return user, nil
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

type userWriteStore interface {
	Save(u User)
}

type userReadStore interface {
	Load(serializeMode string) []User
}

type fileStore struct {
	filePath string
}

func (fs *fileStore) Save(u User) {
	writeUserToFile(u, fs.filePath)
}

func (fs *fileStore) Load(sm string) []User {
	var uStore []User
	data, err := os.ReadFile(fs.filePath)
	if err != nil {

		if os.IsNotExist(err) {
			return nil
		}
		fmt.Println("cannot read file: ", err)
		return nil
	}

	dataStr := string(data)
	users := strings.Split(dataStr, "\n")

	for _, uRow := range users {
		var userStruct = User{}
		switch sm {
		case ManualMode:
			userStruct, err = manualDecode(uRow)
			if err != nil {
				fmt.Println("cannot decode user: ", err)

				return nil
			}
		case JsonMode:
			err = json.Unmarshal([]byte(uRow), &userStruct)
			if err != nil {
				fmt.Println("cannot decode in json mode: ", err)

				return nil
			}
		}
		uStore = append(uStore, userStruct)
	}
	return uStore
}

func userRegister(store userWriteStore) {
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

	u := User{
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

func writeUserToFile(u User, filePath string) {
	_, err := os.Stat(filePath)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cannot create or open file: ", err)

		return
	}

	defer file.Close()

	var data []byte
	if serializationMode == ManualMode {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", u.ID, u.Name, u.Email, u.Password))
	} else if serializationMode == JsonMode {
		data, err = json.Marshal(u)
		if err != nil {
			fmt.Println("cannot encode user to json:", err)

			return
		}

		dataStr := string(data)
		dataStr += "\n"

		data = []byte(dataStr)

	} else {
		fmt.Printf("Invalid serialization mode: %v\n", serializationMode)

		return
	}

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
		userRegister(&fileStore{filePath: userStoragePath})
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
