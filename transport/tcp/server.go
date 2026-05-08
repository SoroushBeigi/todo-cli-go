package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/SoroushBeigi/todo-cli-go/repository/memorystore"
	"github.com/SoroushBeigi/todo-cli-go/service/task"
	"github.com/SoroushBeigi/todo-cli-go/transport/param"
)

func main() {
	const address = "127.0.0.1:2003"
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("Cannot listen on ", address)
	}
	fmt.Println("Listening on ", listener.Addr())

	taskMemoryRepo := memorystore.NewTaskStore()

	taskService := task.NewService(taskMemoryRepo)

	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println("Cannot accept ", address)

			continue
		}

		fmt.Println("connection address", connection.RemoteAddr(), connection.LocalAddr())

		var rawReq = make([]byte, 1024)
		numberOfReadBytes, err := connection.Read(rawReq)
		if err != nil {
			log.Println("Cannot read: ", err)

			continue
		}

		req := &param.Request{}
		if err := json.Unmarshal(rawReq[:numberOfReadBytes], req); err != nil {
			log.Println("bad request:", err)

			continue
		}

		switch req.Command {
		case "create-task":
			res, err := taskService.CreateNewTask(task.CreateRequest{
				Title:      req.CreateTaskRequest.Title,
				DueDate:    req.CreateTaskRequest.DueDate,
				CategoryID: req.CreateTaskRequest.CategoryID,
				UserID:     0,
			})

			if err != nil {
				_, wErr := connection.Write([]byte(err.Error()))
				if wErr != nil {
					log.Println("cannot write")

					continue
				}
			}

			data, err := json.Marshal(&res)
			if err != nil {
				_, wErr := connection.Write([]byte(err.Error()))
				if wErr != nil {
					log.Println("cannot write")

					continue
				}
				continue
			}

			_, wErr := connection.Write(data)
			if wErr != nil {
				log.Println("cannot write")

				continue
			}

		}

	}

}
