package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

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
		if login == "admin" && passwd == "qwerty" {
			return `{"Succ":false}`, -1
		}
		return `{"Succ":true}`, -1

	case "AUTH":
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		if data[0] == "admin" && data[1] == "qwerty" {
			res := authData{31, "asfefmiopifjnwoufdsbhnbfhyiasjfdsan"} // TODO: Id
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
		}
		return `{"Succ":false}`, -1

	case "CHECK_TOKEN":
		data := strings.Trim(message[cmdLen+1:], " ")
		if data == "asfefmiopifjnwoufdsbhnbfhyiasjfdsan" {
			return `{"Succ":true,"Data":"` + string(31) + `"}`, -1 // TODO: Id
		}
		return `{"Succ":false}`, -1

	case "GET_CUR_ISSUES": // <ID>
		// TODO: Id
		issues := []issue{
			{1, "Title 1", 0, "Solve my problem!", true, 1, 10001},
			{2, "Title 2", 0, "Solve my problem too!", true, 2, 10001},
			{3, "Title 3", 0, "Description", true, 3, 10001},
		}
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

	case "GET_USER_ISSUES": // <ID>
		// TODO:
		return "", -1

	case "GET_OPEN_ISSUES":
		// TODO:
		return "", -1

	case "GET_ALL_ISSUES":
		// TODO:
		return "", -1

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
		return `{"Succ":true}`, -1

	default:
		return "ERR UKW CMD", -1
	}
}
