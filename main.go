package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"encoding/json"
)

const (
    connHost = "localhost"
    connPort = "8081"
    connType = "tcp"
)

func main() {
    l, err := net.Listen(connType, connHost + ":" + connPort)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening on " + connHost + ":" + connPort)
    for {
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	len, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	buf = buf[0:len]
	message := string(buf)
	fmt.Printf("Message: %s", message)

	response := makeResponse(message)

	conn.Write([]byte(response))
	conn.Close()
}

type outMessage struct {
	Succ 	bool
	Data	string
}

type issue struct {
	ID 		uint64
	Title 	string
	Time 	uint64
	Descr	string
	IsOpen	bool 		//Status
	UserID	uint64
	TpID	uint64
}

type chatMessage struct {
	ID		uint64
	IssueID	uint64
	Time 	uint64
	Content	string
	MType	uint8 		// TODO: type?
}

func makeResponse(message string) string {
	cmdLen := 0
	for i, v := range message {
		if (v == ' ') {
			cmdLen = i
			break
		}
	}
	switch message[:cmdLen] {
	case "TEST":
		data := strings.Trim(message[cmdLen + 1:], " ")
		res := strings.ToUpper(data)
		return res

	case "REGISTER":
		data := strings.Split(message[cmdLen + 1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		login := data[0]
		passwd := data[1]
		if login == "admin" && passwd == "qwerty" {
			return `{"Succ":false}`
		} 
		return `{"Succ":true}`

	case "AUTH":
		data := strings.Split(message[cmdLen + 1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		if (data[0] == "admin" && data[1] == "qwerty") {
			res := outMessage{true, "asfefmiopifjnwoufdsbhnbfhyiasjfdsan"}
			b, err := json.Marshal(res)
			if err != nil {
				return `{"Succ":false}`
			}
			return string(b)
		}
		return `{"Succ":false}`

	case "GET_USER_ISSUES": // <ID>
		// TODO
		return ""

	case "GET_OPEN_ISSUES":
		// TODO
		return ""

	case "GET_ALL_ISSUES":
		// TODO
		return ""

	case "GET_HELPER_ISSUES ": // <ID>
		// TODO
		return ""

	case "GET_ISSUE": // <ID>
		// TODO
		return ""

	case "GET_SHOP_LIST":
		// TODO
		return ""

	case "GET_FAQ":
		// TODO
		return ""

	default:
		return "ERR UKW CMD"
	}
}