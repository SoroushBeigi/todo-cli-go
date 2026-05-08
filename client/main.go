package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/SoroushBeigi/todo-cli-go/transport/param"
)

func main() {

	message := ""

	if len(os.Args) > 1 {
		message = os.Args[1]
	}

	const address = "127.0.0.1:2003"
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln("Cannot Dial ", address)
	}

	defer conn.Close()

	fmt.Println("local address", conn.LocalAddr())

	req := param.Request{Command: message}

	if req.Command == "create-task" {
		req.CreateTaskRequest = param.CreateTaskRequest{
			Title:      "tst",
			DueDate:    "tst",
			CategoryID: 1,
		}
	}

	serializedData, err := json.Marshal(&req)
	if err != nil {
		log.Fatalln("cannot marshal request")
	}

	numberOfWritten, err := conn.Write([]byte(serializedData))
	if err != nil {
		log.Fatalln("Cannot write ", err)
	}

	fmt.Println("numberOfWrittenBytes: ", numberOfWritten)

	var data = make([]byte, 1024)
	_, rErr := conn.Read(data)
	if rErr != nil {
		log.Fatalln("can't read data from connection", rErr)
	}

	fmt.Println("server response:", string(data))

}
