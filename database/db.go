package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

const (
	USERNAME = "root"
	PASSWORD = "123456"
	HOST     = "localhost"
	PORT     = "3306"
	DATABASE = "yukoreimburse"
	CHARSET  = "utf8mb4"
	DRIVER   = "mysql"
)

// 初始化连接
var Gdb *gorm.DB

func init() {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", USERNAME, PASSWORD, HOST, PORT, DATABASE, CHARSET)

    var err error
    Gdb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		
	})

	if err != nil {
	    panic(err)
	}

	fmt.Println("数据库连接成功")
}