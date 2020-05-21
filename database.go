package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type user struct {
	id       int
	email    string
	password string
	dname    string
	tp       int
}

type issue struct {
	issue_id    int
	status      int
	dname       string
	content     string
	user_id     int
	tp_id       int
	data_create int
}

type message struct {
	msg_id      int
	issue_id    int
	dtype       string
	content     string
	data_create int
}

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
		panic(err)
	}
	defer result.Close()
	if result.Next() {
		user := user{}
		err = result.Scan(&user.id, &user.email, &user.password, &user.dname, &user.tp)
		if err != nil {
			panic(err)
		}
		return user.id
	}
	return -1
}

// если пароль для email верный, то возращает имя, иначе пустую строку
func checkPassword(db *sql.DB, email string, password string) string {
	result, err := db.Query("SELECT * FROM users WHERE email=$1", email)
	if err != nil {
		panic(err)
	}
	defer result.Close()
	if result.Next() {
		user := user{}
		err = result.Scan(&user.id, &user.email, &user.password, &user.dname, &user.tp)
		if err != nil {
			panic(err)
		}
		if user.password == password {
			return user.dname
		}
		return ""
	}
	return ""
}

func isTp(db *sql.DB, id int) bool {
	result, err := db.Query("SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		panic(err)
	}
	defer result.Close()
	if result.Next() {
		user := user{}
		err = result.Scan(&user.id, &user.email, &user.password, &user.dname, &user.tp)
		if err != nil {
			panic(err)
		}
		return user.tp == 1
	}
	return false
}

func getIssues(db *sql.DB, id int) []issue {
	var result *sql.Rows
	var err error
	if isTp(db, id) {
		result, err = db.Query("SELECT * FROM issues WHERE (tp_id=$1 OR status=$2)", id, 0)
	} else {
		result, err = db.Query("SELECT * FROM issues WHERE user_id=$1", id)
	}
	if err != nil {
		panic(err)
	}
	defer result.Close()
	var issues []issue
	for result.Next() {
		issue := issue{}
		err := result.Scan(&issue.issue_id, &issue.status, &issue.dname, &issue.content, &issue.user_id, &issue.tp_id, &issue.data_create)
		if err != nil {
			panic(err)
			continue
		}
		issues = append(issues, issue)
	}
	return issues
}

func database() *sql.DB {
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		panic(err)
	}
	// tp: 0 - user, 1 - tp
	_, err = db.Exec("CREATE TABLE if not exists users (id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, email text, password text, dname text, tp INTEGER)")
	if err != nil {
		panic(err)
	}
	// status: 0 - bot, 1 - open + tp, 2 - close
	_, err = db.Exec("CREATE TABLE if not exists issues (issue_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, status INTEGER, dname text, content text, user_id INTEGER, tp_id INTEGER, data_create INTEGER)")
	if err != nil {
		panic(err)
	}
	_, err = db.Exec("CREATE TABLE if not exists messages (msg_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, issue_id INTEGER, dtype text, content text, data_create INTEGER)")
	if err != nil {
		panic(err)
	}
	return db
}
