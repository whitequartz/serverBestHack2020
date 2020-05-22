package main

import (
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type listenerConn struct {
	chID int64
	conn net.Conn
}

var listeners []listenerConn

func initChatsModule() {
	listeners = make([]listenerConn, 0)
}

func broadcastTo(chID int64, msg chatMessageRaw) {
	ctime := time.Now().Unix()
	newMsg := chatMessage{0, msg.Sender, msg.Dest, ctime, msg.Data}
	data, err := json.Marshal(newMsg)
	if err != nil {
		fmt.Println("Error 8423521")
		return
	}
	tstr := string(data) + "\n"
	toRemove := make([]int, 0)
	for i, v := range listeners {
		if v.chID == chID {
			(func() {
				good := false
				defer (func() {
					if !good {
						v.conn.Close()
						toRemove = append(toRemove, i)
						recover() //! WARN
					}
				})()
				v.conn.Write([]byte(tstr))
				good = true
			})()
		}
	}
	newListeners := make([]listenerConn, 0)
	for i, v := range listeners {
		match := false
		for _, j := range toRemove {
			if i == j {
				match = true
				break
			}
		}
		if !match {
			newListeners = append(newListeners, v)
		}
	}
	listeners = newListeners
}
