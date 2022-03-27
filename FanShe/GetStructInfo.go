package FanShe

import (
	"encoding/json"
	"fmt"
	"my-db/Logs"
	"reflect"
	"strings"
)

func Getstructinfo(tabName string, si interface{}) []byte {
	v := reflect.ValueOf(si)

	if strings.Contains(v.Type().String(), "*") {
		v = v.Elem()
	}
	if v.Kind() == reflect.Slice {
		return getsttu_infs(tabName, v)

	} else {
		return getsttu_inf(tabName, v)
	}
}
func getsttu_inf(tabName string, v reflect.Value) []byte {
	t := v.Type()
	sqlstr := make([]byte, 0)
	sqlstr = append(sqlstr, []byte("insert into `"+tabName+"`(")...)
	sqlfliedname := make([]byte, 0)
	sqlvalues := make([]byte, 0)
	for i := 0; i < v.NumField(); i++ {
		sqlvalues = append(sqlvalues, getsttu_inf_getvalue(v.Field(i))...)
		sqlfliedname = append(sqlfliedname, getsttu_inf_getfiledname(t.Field(i))...)
	}
	sqlstr = append(sqlstr, sqlfliedname[:len(sqlfliedname)-1]...)
	sqlstr = append(sqlstr, 41, 32, 118, 97, 108, 117, 101, 115, 40) //41 32 118 97 108 117 101 115 40 =) values(
	sqlstr = append(sqlstr, sqlvalues[:len(sqlvalues)-1]...)
	sqlstr = append(sqlstr, 41) //)

	return sqlstr

}
func getsttu_infs(tabName string, v reflect.Value) []byte {
	var sqlbyte = make([]byte, 0)
	for i := 0; i < v.Len(); i++ {

		if i > 0 {
			sqlbyte = append(sqlbyte, 44)
			sqlbyte = append(sqlbyte, getsttu_inf(tabName, v.Index(i))[strings.Index(string(getsttu_inf(tabName, v.Index(i))), ") values(")+8:]...)
		} else {
			sqlbyte = append(sqlbyte, getsttu_inf(tabName, v.Index(i))...)
		}
	}
	return sqlbyte
}
func getsttu_inf_getfiledname(t reflect.StructField) []byte {
	if t.Tag.Get("json") != "" {
		return []byte("`" + t.Tag.Get("json") + "`,")
	} else {
		return []byte("`" + strings.ToLower(t.Name) + "`,")
	}
}

func getsttu_inf_getvalue(v reflect.Value) []byte {
	t := v.Type()
	var rb = make([]byte, 0)

	if t.Kind().String() == "string" {

		return []byte("\"" + v.String() + "\",")

	} else if t.Kind().String() == "struct" {
		b, _ := json.Marshal(v.Interface())
		return []byte("'" + string(b) + "',")

	} else if t.Kind().String() == "slice" || t.Kind().String() == "array" {
		if v.Interface() != nil {
			if v.Len() > 0 {
				if v.Index(0).Kind().String() == "struct" {
					rb = append(rb, 39)
					jb, e := json.Marshal(v.Interface())
					if e != nil {
						return []byte("'[]',")
					}
					rb = append(rb, jb...)
					rb = append(rb, 39, 44)

					return rb
				} else {
					Logs.Print(Logs.Warm, "检测到数组数组类型为MYSQL 不支持的类型 将导致出错")
					return nil
				}
			} else {
				Logs.Print(Logs.Warm, "检测到空数组或空JSON对象 数组类型为MYSQL 不支持的类型 将导致出错！")
				return []byte("'[]',")
			}

		} else {
			Logs.Print(Logs.Warm, "Interface error 检测到空数组或空JSON对象 数组类型为MYSQL 不支持的类型 将导致出错！")
			return []byte("'[]',")
		}
	} else {
		return []byte(fmt.Sprintf("%v", v.Interface()) + ",")
	}
}
