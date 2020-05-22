package main

type outMessage struct {
	Succ bool
	Data string
}

type authData struct {
	ID    int
	Token string
}

type issue struct {
	ID     int
	Title  string
	Time   int
	Descr  string
	IsOpen bool //Status
	UserID int
	TpID   int
}

type chatMessage struct {
	ID       int
	SenderID int
	Time     int
	Content  string
	MType    int // 0 - user, 1 - tp, 2 - bot
}

type user struct {
	id       int
	email    string
	password string
	dname    string
	tp       int
}

type chatMessageRaw struct {
	Sender int
	Dest   int
	Data   string
}
