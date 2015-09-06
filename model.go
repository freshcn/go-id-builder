package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var mysql *sql.DB

// 数据表名
var tablename string = "idgenerator"

func init() {
	dsn, b := config.Get("mysql", "dsn")
	if !b {
		panic("没有找到mysql dsn配置")
	}
	tmpMysql, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	mysql = tmpMysql
}

// 获取所有的key的名字列表
func idsList() (arr map[string]int64, err error) {
	arr = make(map[string]int64, 0)
	rows, err := mysql.Query("select `name`,`id` from `" + tablename + "` where `is_del` = 0")
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var id int64
		if err = rows.Scan(&name, &id); err == nil {
			arr[name] = id
		} else {
			return
		}
	}
	return
}

// 向name所指定的id中申请新的ID空间
// 向数据库申请成功后返回新申请到的最大数和申请的数量
func updateId(name string) (num int64, preStep int64, err error) {
	num = 0

	preStep = getPreStep()
	_, err = mysql.Exec("update `"+tablename+"` set id=last_insert_id(id+?) where name=?", preStep, name)
	if err != nil {
		return
	}

	row := mysql.QueryRow(fmt.Sprintf("select last_insert_id() as id  from `%s`", tablename))
	if err != nil {
		return
	}

	err = row.Scan(&num)
	if err != nil {
		return
	}
	return
}
