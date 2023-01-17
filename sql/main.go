package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	//"github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

func main() {
	ExecAllInsertSql("root:Connext@1qaz@WSX@tcp(10.128.0.185:4000)/cloudproject", "db_cmp.t_inf_app_template.sql")
}

var db *sql.DB

func initDatabaseConnect(dataAddr string) {
	dbc, err := sql.Open("mysql", dataAddr)
	if err != nil {
		panic(err)
	}
	err = dbc.Ping()
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connection Success")
	}
	db = dbc
}

func ExecAllInsertSql(dataAddr, fileName string) {
	//初始化数据库连接
	initDatabaseConnect(dataAddr)
	//读取文件内容
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	//按照sql划分数据 ;可能会截取不到完整的sql  sql数据中不能有换行,按照这个逻辑截取;\n
	fileList := strings.Split(string(file), ";\n")
	log.Println(fileList[0], "ssssss")
	//insert Map 根据表名将SQL语句存在map中
	m := make(map[string]string, 0)
	//执行sql标识符 防止最后一张表忘记执行
	now := ""
	//计数器
	ticket := 0
	//计时器
	t := time.Now()
	//数据库设置了 最大单次执行的sql数据大小限制 所以使用count来自适应大小 快要超出就先执行一次
	count := 0
	for i := range fileList {
		ticket++
		if ticket%1000 == 0 {
			log.Println(ticket)
		}
		//找到insert语句
		if strings.Index(fileList[i], "INSERT INTO") != -1 {
			s := fileList[i][strings.Index(fileList[i], "`"):strings.Index(fileList[i], "(")]
			if _, ok := m[s]; !ok {
				//判断是否切换下一张表 如果切换,将上一张表的数据执行
				if s != now {
					//判断是否为初始化数据
					if now != "" {
						ExecSql(m[now], ticket)
						count = 0
					}
					//更新当前表名
					now = s
				}
				//sql语句拼接
				m[s] = "INSERT INTO " + s + fileList[i][strings.Index(fileList[i], "(`"):len(fileList[i])] + ","
			} else {
				//sql语句拼接
				m[s] = m[s] + fileList[i][strings.Index(fileList[i], "VALUES (")+7:len(fileList[i])] + ","
				//数据过大先执行一次
				if count > 1000 {
					ExecSql(m[now], ticket)
					//这里已经执行过了 所以删除就好 会重新拼接一千条sql语句去执行 注意 并未切换下一张表
					delete(m, s)
					count = 0
				} else {
					count++
				}
			}
		} else {
			_, err := db.Exec(fileList[i])
			if err != nil {
				log.Println(err, fileList[i], ticket)
				panic(err)
			}
		}
	}
	//执行最后一张表的数据插入
	ExecSql(m[now], ticket)
	fmt.Println("完成 共插入", ticket, "行数据;", "共用时", time.Now().Sub(t))
}

func ExecSql(sql string, ticket int) {
	sql = sql[:len(sql)-1] + ";"
	DBUnlock()
	_, err := db.Exec(sql)
	if err != nil {
		log.Println(sql, ticket)
		panic(err)
	}
}

// DBUnlock 解锁表格 否则执行会不成功
func DBUnlock() {
	_, err := db.Exec("unlock tables;")
	if err != nil {
		panic(err)
	}

}
