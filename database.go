package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

// создаёт запись в дб по заданным параметрам
func register(db *sql.DB, email string, password string, name string) int64 {
	result, err := db.Exec("INSERT INTO users (email,password,dname,tp) VALUES ($1,$2,$3,$4)", email, password, name, 0)
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
func checkPassword(db *sql.DB, email string, password string) int {
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
			return user.ID
		}
	}
	return -1
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
		result, err = db.Query("SELECT * FROM issues WHERE (tp_id=$1 OR status=$2)", id, 1)
	} else {
		result, err = db.Query("SELECT * FROM issues WHERE (user_id=$1 AND status=$2)", id, 1)
	}
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	issues := []issue{}
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
	issues := []issue{}
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
	issues := []issue{}
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

func getMessagesHistory(db *sql.DB, issue_id int) []chatMessage {
	result, err := db.Query("SELECT * FROM messages WHERE issue_id=$1", issue_id)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	messages := []chatMessage{}
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

func removeIssue(db *sql.DB, issue_id int) {
	_, err := db.Exec("DELETE FROM issues WHERE issue_id=$1", issue_id)
	if err != nil {
		fmt.Println(err)
	}
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

func addMessage(db *sql.DB, sender_id int, issue_id int, dtype int, message string, time int) int64 {
	result, err := db.Exec("INSERT INTO messages (sender_id,issue_id,dtype,content,data_create) VALUES ($1,$2,$3,$4,$5)", sender_id, issue_id, dtype, message, time)
	if err != nil {
		fmt.Println(err)
	}
	id, _ := result.LastInsertId()
	return id
}

func getDevices(db *sql.DB, user_id int) []usersDevices {
	result, err := db.Query("SELECT * FROM devices WHERE user_id=$1", user_id)
	if err != nil {
		fmt.Println(err)
	}
	defer result.Close()
	devices := []usersDevices{}
	for result.Next() {
		device := usersDevices{}
		err := result.Scan(&device.ID, &device.UserID, &device.Type, &device.Model, &device.Cost, &device.BuyTime, &device.ValidTime)
		if err != nil {
			fmt.Println(err)
			continue
		}
		devices = append(devices, device)
	}
	return devices
}

func addDevice(db *sql.DB, user_id int, dtype int, model string, cost int, buy_time int, valid_time int) int64 {
	result, err := db.Exec("INSERT INTO devices (user_id,dtype,model,cost,buy_time,valid_time) VALUES ($1,$2,$3,$4,$5,$6)", user_id, dtype, model, cost, buy_time, valid_time)
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
	_, err = db.Exec("CREATE TABLE if not exists users (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, email text, password text, dname text, tp INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	// status: 1 - open + tp, 0 - close
	_, err = db.Exec("CREATE TABLE if not exists issues (issue_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, status INTEGER, dname text, content text, user_id INTEGER, tp_id INTEGER, data_create INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	// dtype: 0 - user, 1 - tp, -1 - bot
	_, err = db.Exec("CREATE TABLE if not exists messages (msg_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, sender_id INTEGER, issue_id INTEGER, dtype INTEGER, content text, data_create INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	_, err = db.Exec("CREATE TABLE if not exists devices (device_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, user_id INTEGER, dtype INTEGER, model text, cost INTEGER, buy_time INTEGER, valid_time INTEGER)")
	if err != nil {
		fmt.Println(err)
	}
	return db
}
