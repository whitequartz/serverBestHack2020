package main

type outMessage struct {
	Succ bool
	Data string
}

type authData struct {
	ID    int64
	Token string
}

type issue struct {
	ID     int
	Title  string
	Time   int64
	Descr  string
	IsOpen bool //Status
	UserID int64
	TpID   int64
}

type chatMessage struct {
	ID       int64
	SenderID int64
	IssueID  int64
	Time     int64
	Content  string
}

type chatMessageRaw struct {
	Sender int64
	Dest   int64
	Data   string
}
