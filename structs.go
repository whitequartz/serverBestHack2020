package main

type outMessage struct {
	Succ bool
	Data string
}

type issue struct {
	ID     uint64
	Title  string
	Time   uint64
	Descr  string
	IsOpen bool //Status
	UserID uint64
	TpID   uint64
}

type chatMessage struct {
	ID      uint64
	IssueID uint64
	Time    uint64
	Content string
	MType   uint8 // TODO: type?
}
