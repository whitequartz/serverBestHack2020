package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "8081"
    CONN_TYPE = "tcp"
)

func main() {
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
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

type inMessage struct {
	cmd 	string
	data	string
}

type outMessage struct {
	succ 	bool
	data	string
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
	default:
		return "ERR UKW CMD"
	}
}