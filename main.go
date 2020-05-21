package main

import (
	"fmt"
	"net"
	"os"
)

const (
	connHost = ""
	connPort = "8081"
	connType = "tcp"
)

func main() {
	l, err := net.Listen(connType, connHost+":"+connPort)
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
