package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// создаёт запись в дб по заданным параметрам
func register(db *sql.DB, email string, password string, name string, tp int) int64 {
	result, err := db.Exec("INSERT INTO users (email,password,dname,tp) VALUES ($1,$2,$3,$4)", email, password, name, tp)
	if err != nil {
		panic(err)
	}
	id, _ := result.LastInsertId()
	return id
}

// получить id по email. Если записи в дб не существует, то возвращает -1
func getUserId(db *sql.DB, email string) int {
	result, err := db.Query("SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	if result.Next() {
		user := user{}
		err = result.Scan(&user.ID, &user.Email, &user.Password, &user.Dname, &user.Tp)
		if err != nil {
			fmt.Println(err)
		}
		return user.ID
	}
	return -1
}

// если пароль для email верный, то возращает имя, иначе пустую строку
func checkPassword(db *sql.DB, email string, password string) string {
	result, err := db.Query("SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	if result.Next() {
		user := user{}
		err = result.Scan(&user.ID, &user.Email, &user.Password, &user.Dname, &user.Tp)
		if err != nil {
			fmt.Println(err)
		}
		if user.Password == password {
			return user.Dname
		}
		return ""
	}
	return ""
}

func isTp(db *sql.DB, id int) bool {
	result, err := db.Query("SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	if result.Next() {
		user := user{}
		err = result.Scan(&user.ID, &user.Email, &user.Password, &user.Dname, &user.Tp)
		if err != nil {
			fmt.Println(err)
		}
		return user.Tp == 1
	}
	return false
}

func getCurIssues(db *sql.DB, id int) []issue {
	var result *sql.Rows
	var err error
	if isTp(db, id) {
		result, err = db.Query("SELECT * FROM issues WHERE (tp_id=$1 OR status=$2)", id, 0)
	} else {
		result, err = db.Query("SELECT * FROM issues WHERE user_id=$1", id)
	}
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	var issues []issue
	for result.Next() {
		issue := issue{}
		var status int
		err := result.Scan(&issue.ID, &status, &issue.Title, &issue.Descr, &issue.UserID, &issue.TpID, &issue.Time)
		issue.IsOpen = status == 1
		if err != nil {
			fmt.Println(err)
			continue
		}
		issues = append(issues, issue)
	}
	return issues
}

func getUserIssues(db *sql.DB, id int) []issue {
	result, err := db.Query("SELECT * FROM issues WHERE user_id=$1", id)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	var issues []issue
	for result.Next() {
		issue := issue{}
		var status int
		err := result.Scan(&issue.ID, &status, &issue.Title, &issue.Descr, &issue.UserID, &issue.TpID, &issue.Time)
		issue.IsOpen = status == 1
		if err != nil {
			fmt.Println(err)
			continue
		}
		issues = append(issues, issue)
	}
	return issues
}

func getAllIssues(db *sql.DB) []issue {
	result, err := db.Query("SELECT * FROM issues")
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	var issues []issue
	for result.Next() {
		issue := issue{}
		var status int
		err := result.Scan(&issue.ID, &status, &issue.Title, &issue.Descr, &issue.UserID, &issue.TpID, &issue.Time)
		issue.IsOpen = status == 1
		if err != nil {
			fmt.Println(err)
			continue
		}
		issues = append(issues, issue)
	}
	return issues
}

func getMessagesHistory(db *sql.DB, sender_id int) []chatMessage {
	result, err := db.Query("SELECT * FROM messages WHERE sender_id=$1", sender_id)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	var messages []chatMessage
	for result.Next() {
		chatMessage := chatMessage{}
		err := result.Scan(&chatMessage.ID, &chatMessage.SenderID, &chatMessage.IssueID, &chatMessage.MType, &chatMessage.Content, &chatMessage.Time)
		if err != nil {
			fmt.Println(err)
			continue
		}
		messages = append(messages, chatMessage)
	}
	return messages
}

func addIssue(db *sql.DB, title string, message string, user_id int, time int) int64 {
	result, err := db.Exec("INSERT INTO issues (status,dname,content,user_id,tp_id,data_create) VALUES ($1,$2,$3,$4,$5,$6)", 1, title, message, user_id, -1, time)
	if err != nil {
		fmt.Println(err)
	}
	id, _ := result.LastInsertId()
	return id
}

func addTpForIssue(db *sql.DB, issue_id int, tp_id int) {
	_, err := db.Exec("UPDATE issues SET tp_id=$1 WHERE issue_id=$2", tp_id, issue_id)
	if err != nil {
		fmt.Println(err)
	}
}

func closeIssue(db *sql.DB, issue_id int) {
	_, err := db.Exec("UPDATE issues SET status=$1 WHERE issue_id=$2", 0, issue_id)
	if err != nil {
		fmt.Println(err)
	}
}

func addMessage(db *sql.DB, sender_id int, issue_id, dtype int, message string, time int) int64 {
	result, err := db.Exec("INSERT INTO message (sender_id,issue_id,dtype,content,data_create) VALUES ($1,$2,$3,$4)", sender_id, issue_id, dtype, message, time)
	if err != nil {
		fmt.Println(err)
	}
	id, _ := result.LastInsertId()
	return id
}

func database() *sql.DB {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		fmt.Println(err)
	}
	// tp: 0 - user, 1 - tp
	_, err = db.Exec("CREATE TABLE if not exists users (id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL, email text, password text, dname text, tp INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	// status: 1 - open + tp, 0 - close
	_, err = db.Exec("CREATE TABLE if not exists issues (issue_id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL, status INTEGER, dname text, content text, user_id INTEGER, tp_id INTEGER, data_create INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	// dtype: 0 - user, 1 - tp, 2 - bot
	_, err = db.Exec("CREATE TABLE if not exists messages (msg_id INTEGER PRIMARY KEY AUTO_INCREMENT NOT NULL, sender_id INTEGER, issue_id INTEGER, dtype INTEGER, content text, data_create INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	return db
}
