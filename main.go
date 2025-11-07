package main

import (
	"fmt"
	advancedone "gorm/gormAdvanced"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 尝试连接数据库
	db, err := connectDatabase()
	if err != nil {
		fmt.Printf("警告: 数据库连接失败: %v\n", err)
	} else {

		//gormSql.Run(db)//基本CRUD操作
		//gormSqlTwo.Run(db)//事务语句
		//sqlxone.Run(db)//Sqlx题目一
		//sqlxtwo.Run(db) //Sqlx题目二
		advancedone.Run(db) //进阶gorm
		//fmt.Println("数据库操作执行完毕")
	}
}

// connectDatabase 尝试连接到数据库
func connectDatabase() (*gorm.DB, error) {
	dsn := "root:123456@tcp(localhost:3306)/grom?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
