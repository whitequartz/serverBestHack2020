package main

import (
	"encoding/json"
	"strings"
)

func makeResponse(message string) string {
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
		return res

	case "REGISTER":
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		login := data[0]
		passwd := data[1]
		if login == "admin" && passwd == "qwerty" {
			return `{"Succ":false}`
		}
		return `{"Succ":true}`

	case "AUTH":
		data := strings.Split(message[cmdLen+1:], " ")
		for i := range data {
			data[i] = strings.Trim(data[i], " \n\t")
		}
		if data[0] == "admin" && data[1] == "qwerty" {
			res := outMessage{true, "asfefmiopifjnwoufdsbhnbfhyiasjfdsan"}
			b, err := json.Marshal(res)
			if err != nil {
				return `{"Succ":false}`
			}
			return string(b)
		}
		return `{"Succ":false}`

	case "GET_CUR_ISSUES": // <ID>
		// TODO: Id
		issues := []issue{
			{1, "Title 1", 0, "Solve my problem!", true, 1, 10001},
			{2, "Title 2", 0, "Solve my problem too!", true, 2, 10001},
			{3, "Title 3", 0, "Description", true, 3, 10001},
		}
		b, err := json.Marshal(issues)
		if err != nil {
			return `{"Succ":false}`
		}
		out := outMessage{true, string(b)}
		b, err = json.Marshal(out)
		if err != nil {
			return `{"Succ":false}`
		}
		return string(b)

	case "GET_USER_ISSUES": // <ID>
		// TODO
		return ""

	case "GET_OPEN_ISSUES":
		// TODO
		return ""

	case "GET_ALL_ISSUES":
		// TODO
		return ""

	case "GET_HELPER_ISSUES ": // <ID>
		// TODO
		return ""

	case "GET_ISSUE": // <ID>
		// TODO
		return ""

	case "GET_SHOP_LIST":
		// TODO
		return ""

	case "GET_FAQ":
		// TODO
		return ""

	default:
		return "ERR UKW CMD"
	}
}
