package main

import "net"

type listenerConn struct {
	chID int64
	conn net.Conn
}

var listeners []listenerConn

func initChatsModule() {
	listeners = make([]listenerConn, 0)
}

func broadcastTo(chID int64, data []byte) {
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
				v.conn.Write(data)
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
