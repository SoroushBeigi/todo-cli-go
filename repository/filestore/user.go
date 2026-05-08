package filestore

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/SoroushBeigi/todo-cli-go/constant"
	"github.com/SoroushBeigi/todo-cli-go/entity"
)

func manualDecode(uRow string) (entity.User, error) {
	uRow = strings.TrimSpace(uRow)
	if uRow == "" {
		return entity.User{}, fmt.Errorf("Empty user")
	}

	userFields := strings.Split(uRow, ",")
	user := entity.User{}
	for _, field := range userFields {
		field = strings.TrimSpace(field)
		content := strings.Split(field, ": ")
		if len(content) < 2 {

			return entity.User{}, fmt.Errorf("content length <2")
		}
		fieldName := content[0]
		fieldValue := content[1]

		switch fieldName {
		case "id":
			id, err := strconv.Atoi(fieldValue)
			if err != nil {
				return entity.User{}, fmt.Errorf("str convert error")
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
	fmt.Printf("%+v\n", user)
	return user, nil
}

type FileStore struct {
	filePath          string
	serializationMode string
}

func NewFileStore(path, serializationMode string) FileStore {
	return FileStore{filePath: path, serializationMode: serializationMode}
}

func (fs *FileStore) Save(u entity.User) {
	fs.writeUserToFile(u, fs.filePath)
}

func (fs *FileStore) Load() []entity.User {
	var uStore []entity.User
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
		var userStruct = entity.User{}
		switch fs.serializationMode {
		case constant.ManualMode:
			userStruct, err = manualDecode(uRow)
			if err != nil {
				fmt.Println("cannot decode user: ", err)

				return nil
			}
		case constant.JsonMode:
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

func (fs *FileStore) writeUserToFile(u entity.User, filePath string) {
	_, err := os.Stat(filePath)

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cannot create or open file: ", err)

		return
	}

	defer file.Close()

	var data []byte
	if fs.serializationMode == constant.ManualMode {
		data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n", u.ID, u.Name, u.Email, u.Password))
	} else if fs.serializationMode == constant.JsonMode {
		data, err = json.Marshal(u)
		if err != nil {
			fmt.Println("cannot encode user to json:", err)

			return
		}

		dataStr := string(data)
		dataStr += "\n"

		data = []byte(dataStr)

	} else {
		fmt.Printf("Invalid serialization mode: %v\n", fs.serializationMode)

		return
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		fmt.Println("Error writing file: ", err)
	}
}
