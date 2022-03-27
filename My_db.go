package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"my-db/FanShe"
	"my-db/Logs"
	"my-db/Str_Handle"
	"strconv"
	"strings"
	"time"
)

var (
	db  *sql.DB
	err error
)
var ISlog bool

type my_db struct {
}

type where struct {
	sqlstr string
	where  string
	err    error
	obj    interface{}
}

//set first a where
func (dw where) Where(awhere string) where {
	dw.where = dw.where + awhere
	return dw
}

//set Multiple and where demo: obj.X.AndWhere(`id=0`,`name="lisa"`...)
func (aw where) AndWhere(awhere ...string) where {
	wheres := ""
	if awhere != nil {
		for i := 0; i < len(awhere); i++ {
			wheres = wheres + " and " + awhere[i]
		}
	}
	aw.where = aw.where + wheres
	return aw
}

//set Multiple or where demo: obj.X.orWhere(`id=0`,`name="lisa"`...)
func (ow where) OrWhere(awhere ...string) where {
	wheres := ""
	if awhere != nil {
		for i := 0; i < len(awhere); i++ {
			wheres = wheres + " or " + awhere[i]
		}
	}
	ow.where = ow.where + wheres
	return ow
}

//run
func (ow where) Run() error {
	var data []byte
	if ow.where != "" {
		if ow.where[:4] == " and" || ow.where[:4] == " or " {
			//Logs.Print(Logs.Warm, "检测到第一个条件为 AndWhere 或 OrWhere 已自动更正 稍后请更正！！！")
			ow.sqlstr = ow.sqlstr + " where " + ow.where[4:]

		} else {
			ow.sqlstr = ow.sqlstr + " where " + ow.where
		}
	}
	if ISlog {
		Logs.Print(Logs.Print_log, "sql  is: "+ow.sqlstr)
	}

	if strings.Contains(ow.sqlstr, "select") {

		data, err = dbquery(ow.sqlstr)
		//fmt.Println(string(data))
		if err != nil {
			return err
		}
		if data == nil {
			return nil
		}
		err = json.Unmarshal(data, ow.obj)
		if ISlog {
			fmt.Println("json ", err, string(data))
		}

		return nil
	}

	_, err = db.Exec(ow.sqlstr)
	return err
}

//Only supported "MYSQL" dataSourceName = Database Connection string
func New(dataSourceName string) (*my_db, error) {

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {

		return nil, fmt.Errorf("sql open database or dataSourceName err!")
	}
	err = db.Ping()
	if err != nil {

		return nil, fmt.Errorf("ping database error")

	}
	return new(my_db), nil
}

//Demo obj.Select(&toobj,"tabname","filed1","filed2"...).Where("id=1").AndWhere(`name="lisa"`,`age=20`...).OrWhere(`gender="wenman"`).Run return error
func (my *my_db) Select(toObj interface{}, tablename string, flideName ...string) where {
	filedns := Str_Handle.Update_str(flideName...)
	fmt.Println(filedns)
	if filedns == "" {
		filedns = "*"
	}
	return where{sqlstr: "select " + filedns + " from `" + tablename + "`", err: nil, obj: toObj}
}

// inputObj Can be slice struct Obj or struct Obj
//Demo obj.Insert(`tabname`,&inputObj)  or obj.Insert(`tabname`,inputObj)
func (my *my_db) Insert(tabName string, v interface{}) error {

	sqls := string(FanShe.Getstructinfo(tabName, v))
	if ISlog {
		Logs.Print(Logs.Print_log, "sql insert is: "+sqls)
	}

	_, err = db.Exec(sqls)
	if err != nil {
		Logs.Print(Logs.Error, err)
	}
	return err
}

//Demo obj.Delete(`tabname`).Where(`id=1`).AndWhere(...).OrWhere(...).Run return err
func (my *my_db) Delete(tablename string) where {
	return where{sqlstr: "delete from `" + tablename + "`"}
}

//Demo obj.Update(`tabname`,`id=1`,`name="lisa"`,`age=22`....).Where(`id=1`).AndWhere(...).OrWhere(...).Run return err
func (my my_db) Update(tablename string, set ...string) where {
	return where{sqlstr: "update `" + tablename + "` set " + Str_Handle.Update_str(set...)}
}

func dbquery(sql_str string) ([]byte, error) {

	if db == nil {
		Logs.Print(2, "没有初始化数据库连接 调用New方法初始化数据湖连接 ")
		return nil, fmt.Errorf("没有初始化数据库连接 调用New方法初始化数据湖连接 ")
	}
	var rows *sql.Rows
	var coln []string
	var byterejson = make([]byte, 0)

	rows, err = db.Query(sql_str)
	if err != nil {
		Logs.Print(2, "sql语句 有误")
		return nil, fmt.Errorf("sql语句 有误")
	}
	coln, err = rows.Columns()
	if err != nil {
		Logs.Print(2, "sql语句 有误")
		return nil, fmt.Errorf("sql语句 有误")
	}
	colt, _ := rows.ColumnTypes()
	if err != nil {
		Logs.Print(2, "sql语句 有误")
		return nil, fmt.Errorf("sql语句 有误")
	}
	var di = make([][]byte, len(coln), len(coln))

	var lsdi = make([]interface{}, len(coln), len(coln))

	for i := 0; i < len(coln); i++ {
		lsdi[i] = &di[i]
	}
	byterejson = append(byterejson, 91)
	js := 0
	for rows.Next() {
		js++
		rows.Scan(lsdi...)
		//34-" 58-: 44-, 123-{ 125-} 102-f 97-a l-108 s-115 e-101
		byterejson = append(byterejson, 123)
		for i := 0; i < len(di); i++ {
			if strings.Contains(colt[i].DatabaseTypeName(), "CHAR") || strings.Contains(colt[i].DatabaseTypeName(), "TEXT") || strings.Contains(colt[i].DatabaseTypeName(), "DATE") || strings.Contains(colt[i].DatabaseTypeName(), "TIME") || strings.Contains(colt[i].DatabaseTypeName(), "YEAR") || strings.Contains(colt[i].DatabaseTypeName(), "BLOB") {
				//byterejson = append(byterejson, 34)
				byterejson = append(byterejson, []byte("\""+coln[i]+"\":\""+string(di[i])+"\",")...)

			} else if colt[i].DatabaseTypeName() == "JSON" {
				if di[i] == nil {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":"+string(di[i])+"{},")...)
				} else {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":"+string(di[i])+",")...)
				}
			} else if colt[i].DatabaseTypeName() == "TINYINT" {
				if string(di[i]) == "" {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":false,")...)
				} else if colt[i].DatabaseTypeName() == "0" {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":false,")...)
				} else {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":true,")...)
				}
			} else if colt[i].DatabaseTypeName() == "INT" {
				if di[i] == nil {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":null,")...)
				} else {
					byterejson = append(byterejson, []byte("\""+coln[i]+"\":"+string(di[i])+",")...)
				}
			} else {
				byterejson = append(byterejson, []byte("\""+coln[i]+"\":"+string(di[i])+",")...)
			}

		}
		byterejson = append(byterejson[:len(byterejson)-1], []byte("},")...)

	}
	if js > 1 {
		byterejson = append(byterejson[:len(byterejson)-1], 93)
	} else {
		if len(byterejson) <= 2 {
			return nil, nil
		}
		byterejson = byterejson[1 : len(byterejson)-1]
	}

	return byterejson, nil
}

type Test struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Isman bool   `json:"isman"`
}

func main() {
	dh, _ := New("root:666666@tcp(127.0.0.1)/li_db")
	var t []Test
	ISlog = true
	dh.Update("test2", `id=520`).Where("id>0").Run()
	//dh.Delete("test2").Run()
	dh.Insert("test2", &Test{})
	dh.Select(&t, "test2", "id", "name").Run()
	fmt.Println(t)
}
func testgrominsert(ts []Test) {
	gdb, err := gorm.Open("mysql", "root:666666@tcp(127.0.0.1)/li_db")
	fmt.Println(err)
	s := time.Now()
	for i := 0; i < len(ts); i++ {
		gdb.Table("test2").Create(&ts[i])
	}
	ends := time.Since(s)
	f64, err := strconv.ParseFloat(fmt.Sprintf("%v", len(ts)), 64)
	fmt.Println("共 ", len(ts), " 条 记录 累计时间/s ", ends)
	fmt.Println("平均时间 1000/s ", ends.Seconds()/f64)

}

func testmydbinsert_SSS(ts []Test) {
	dm, err := New("root:666666@tcp(127.0.0.1)/li_db")
	s := time.Now()
	err = dm.Insert("test2", &ts)
	if err != nil {
		panic("插入出错" + err.Error())
	}
	ends := time.Since(s)
	f64, err := strconv.ParseFloat(fmt.Sprintf("%v", len(ts)), 64)
	fmt.Println("共 ", len(ts), " 条 记录 累计时间/s ", ends)
	fmt.Println("平均时间 1000/s ", ends.Seconds()/f64)

}

func testmydbinsert(ts []Test) {
	dm, err := New("root:666666@tcp(127.0.0.1)/li_db")
	s := time.Now()
	for i := 0; i < len(ts); i++ {
		err = dm.Insert("test2", &ts[i])
		if err != nil {
			panic("插入出错" + err.Error())
		}

	}
	ends := time.Since(s)
	f64, err := strconv.ParseFloat(fmt.Sprintf("%v", len(ts)), 64)
	fmt.Println("共 ", len(ts), " 条 记录 累计时间/s ", ends)
	fmt.Println("平均时间 1000/s ", ends.Seconds()/f64)

}
