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
	IssueID  int
	Time     int
	Content  string
	MType    int // 0 - user, 1 - tp, 2 - bot
}

type user struct {
	ID       int
	Email    string
	Password string
	Dname    string
	Tp       int
}

type chatMessageRaw struct {
	Sender int
	Dest   int
	Data   string
}
