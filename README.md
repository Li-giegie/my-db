#Hello stranger 
#你好陌生人
    nice to meet you!!
    很高兴见到你
## why is it that "my-db"
##为什么会有 "my-db"
### 1. simple 
###简单好用
    Extremely simple    ：非常简单
    High performance    ：高性能
### 2. support Database data type "JSON" Yes, that's great !
###支持数据库数据类型“JSON”是的，太棒了！
    struct field obj Can be mapped to Database Data Type "JSON"
### 3. support slice struct obj handle
###支持切片结构对象
####Sample common part:示例公共部分

    type A struct {
        Id    int    `json:"id"`
        Name  string `json:"name"`
        Isman bool   `json:"isman"`
    }
    var a =make([]A,10)
    a[0]=A{Id:0,Name:"lisa",Isman:false}
    ...
    a[9]=A{Id:0,Name:"lisa",Isman:false}
    dh,_:=New("root:666666@tcp(127.0.0.1)/li_db") //my test DB string
    
####Demo Insert:插入 例子
    //Encouraged operation
    //鼓励的操作 性能非常高 应对大量数据性能非常好
    dh.Insert("tabname",&a)

    //性能不如第一个 应对大量数据性能差
    //Discouraged operation
    for _,v:range a{
        dh.Insert("tabname",v) //encourage handle
    }
####Demo Select:查询 例子
    var b A
    dh.Select(&b,`tabname`,`filed1`,`filed2`...).Where(`id=0`).AndWhere(`name="lisa"`,`age=1`...).Run()
####Demo Update:更新 例子
    var b A
    dh.Update(`tabname`,`name="test"`,`age=22`...).Where(`id=0`).OrWhere(`name="lisa"`,`age=1`...).Run()
####Demo Delete:删除 例子
    dh.Delete(`tabname`).Where(`id=0`).OrWhere(`name="lisa"`,`age=1`...).Run()
### 4.print SQL
### 输出SQL语句
    dh.ISLog=true