package gormSql

import (
	"fmt"

	"gorm.io/gorm"
)

// 题目1：基本CRUD操作
// 假设有一个名为 students 的表，包含字段 id （主键，自增）、 name （学生姓名，字符串类型）、 age （学生年龄，整数类型）、 grade （学生年级，字符串类型）。

type Students struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"type:varchar(100)"`
	Age   int    `gorm:"type:int"`
	Grade string `gorm:"type:varchar(50)"`
}

func Run(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&Students{})

	// 1. 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"。
	student := Students{Name: "张三", Age: 20, Grade: "三年级"}
	result := db.Create(&student)
	if result.Error != nil {
		fmt.Printf("插入学生记录失败: %v\n", result.Error)
	} else {
		fmt.Printf("成功插入学生记录，ID: %d\n", student.ID)
	}

	// 2. 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息。
	var studentsAbove18 []Students
	db.Where("age > ?", 18).Find(&studentsAbove18)
	for _, s := range studentsAbove18 {
		fmt.Printf("ID: %d, 姓名: %s, 年龄: %d, 年级: %s\n", s.ID, s.Name, s.Age, s.Grade)
	}

	// // 3. 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"。
	// result := db.Model(&Students{}).Where("name = ?", "张三").Update("grade", "四年级")
	// if result.Error != nil {
	// 	fmt.Printf("更新学生年级失败: %v\n", result.Error)
	// } else if result.RowsAffected > 0 {
	// 	fmt.Printf("成功更新%d条学生记录\n", result.RowsAffected)
	// } else {
	// 	fmt.Println("未找到匹配的学生记录")
	// }

	// // 4. 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录。
	// fmt.Println("\n4. 删除年龄小于15岁的学生...")
	// result := db.Where("age < ?", 15).Delete(&Students{})
	// if result.Error != nil {
	// 	fmt.Printf("删除学生记录失败: %v\n", result.Error)
	// } else {
	// 	fmt.Printf("成功删除%d条学生记录\n", result.RowsAffected)
	// }

}
