package dao

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/url"
	"time"
)

var DB *sql.DB

func init() {
	dataSourceName := fmt.Sprintf("root:password@tcp(localhost:3306)/blog?charset=utf8mb4&loc=%s&parseTime=true", url.QueryEscape("Asia/Shanghai"))
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic("连接数据库失败" + err.Error())
	}

	//最大空闲连接数，默认不配置，是2个最大空闲连接
	db.SetMaxIdleConns(5)
	//最大连接数，默认不配置，是不限制最大连接数
	db.SetMaxOpenConns(100)
	// 连接最大存活时间
	db.SetConnMaxLifetime(time.Minute * 3)
	//空闲连接最大存活时间
	db.SetConnMaxIdleTime(time.Minute * 1)

	DB = db
}
