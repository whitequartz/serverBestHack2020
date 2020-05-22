package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

var db = database()

func makeResponse(message string) (string, int64) {
	cmdLen := 0
	for i, v := range message {
		if v == ' ' {
			cmdLen = i
			break
		}
	}
	switch message[:cmdLen] {
	case "TEST":
		data := strings.Trim(message[cmdLen+1:], " ")
		res := strings.ToUpper(data)
		out := outMessage{true, res}
		b, err := json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return string(b), -1

	case "REGISTER":
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		login := data[0]
		passwd := data[1]
		name := data[2]
		if getUserId(db, login) == -1 {
			register(db, login, passwd, name)
			return `{"Succ":true}`, -1
		}
		return `{"Succ":false}`, -1

	case "AUTH":
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		id := checkPassword(db, data[0], data[1])
		if id == -1 {
			return `{"Succ":false}`, -1 // неправильный пароль
		}
		encryptMsg, _ := encrypt(key, string(id))
		res := authData{id, encryptMsg}
		b, err := json.Marshal(res)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		out := outMessage{true, string(b)}
		b, err = json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return string(b), -1

	case "CHECK_TOKEN":
		data := strings.Trim(message[cmdLen+1:], " \n")
		msg, err := decrypt(key, data)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return `{"Succ":true,"Data":"` + msg + `"}`, -1

	case "GET_CUR_ISSUES": // <ID>
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		id, _ := strconv.Atoi(data[0])
		issues := getCurIssues(db, id)
		b, err := json.Marshal(issues)
		fmt.Println(string(b))
		if err != nil {
			return `{"Succ":false}`, -1
		}
		out := outMessage{true, string(b)}
		b, err = json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return string(b), -1

	case "GET_USER_ISSUES": // <ID>
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		id, _ := strconv.Atoi(data[0])
		issues := getUserIssues(db, id)
		b, err := json.Marshal(issues)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		out := outMessage{true, string(b)}
		b, err = json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return string(b), -1

	case "GET_OPEN_ISSUES":
		// TODO:
		return "", -1

	case "GET_ALL_ISSUES":
		issues := getAllIssues(db)
		b, err := json.Marshal(issues)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		out := outMessage{true, string(b)}
		b, err = json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return string(b), -1

	case "GET_HELPER_ISSUES ": // <ID>
		// TODO
		return "", -1

	case "GET_ISSUE": // <ID>
		// TODO
		return "", -1

	case "GET_SHOP_LIST":
		// TODO
		return "", -1

	case "GET_FAQ":
		// TODO
		return "", -1

	case "LISTEN":
		data := strings.Split(message[cmdLen+1:], " ")
		data[0] = strings.Trim(data[0], " \n\t")
		ch, err := strconv.ParseInt(data[0], 10, 64)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return "", ch

	case "SEND_MSG":
		data := []byte(strings.Trim(message[cmdLen+1:], " "))
		for i, v := range data {
			if v == '\n' {
				data[i] = ' '
			}
		}
		raw := chatMessageRaw{}
		json.Unmarshal(data, &raw)
		broadcastTo(raw.Dest, raw)
		if raw.Sender == -1 {
			addMessage(db, raw.Sender, raw.Dest, -1, raw.Data, int(time.Now().Unix()))
		} else if isTp(db, raw.Sender) {
			addMessage(db, raw.Sender, raw.Dest, 1, raw.Data, int(time.Now().Unix()))
		} else {
			addMessage(db, raw.Sender, raw.Dest, 0, raw.Data, int(time.Now().Unix()))
		}
		return `{"Succ":true}`, -1

	case "HISTORY":
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		id, _ := strconv.Atoi(data[0])
		messages := getMessagesHistory(db, id)
		b, err := json.Marshal(messages)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		out := outMessage{true, string(b)}
		b, err = json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`, -1
		}
		return string(b), -1

	default:
		return "ERR UKW CMD", -1
	}
}
