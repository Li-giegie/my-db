package Str_Handle

import "strings"

func Update_str(set ...string) string {
	sqlstr := ""
	if set == nil {
		return ""
	} else {
		for i := 0; i < len(set); i++ {
			if !strings.Contains(set[i], "`") && !strings.Contains(set[i], "=") {
				sqlstr = sqlstr + "`" + set[i] + "`, "
			} else {
				sqlstr = sqlstr + set[i] + ", "
			}

		}
		return sqlstr[:len(sqlstr)-2]
	}
}
