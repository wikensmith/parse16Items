package db

import (
	"log"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	)

var database *gorm.DB

// 连接数据库，并创建表
func InitDB(uri string) {
	var err error
	database, err = gorm.Open("mysql", uri)
	if err != nil {
		log.Fatalf("连接数据库失败，原因:%v\n", err)
	}
}

func DBClose()  {
	_ = database.Close()

}
