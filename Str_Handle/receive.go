package Str_Handle

import (
	"my-db/Logs"
	"strings"
)

func Receive(sql_str string, args ...string) string {
	if args == nil {
		return sql_str
	}
	if strings.Count(sql_str, "?") != len(args) {
		Logs.Print(Logs.Error, "sql_str and args number")
		return ""
	}
	for _, v := range args {
		//Logs.Print(Logs.Print_log, v)
		sql_str = strings.Replace(sql_str, "?", v, 1)
	}
	return sql_str
}
