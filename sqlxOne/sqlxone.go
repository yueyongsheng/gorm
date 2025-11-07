package sqlxone

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

// 题目1：使用SQL扩展库进行查询
// 假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
// 要求 ：
// 编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
// 编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。

// Employee 员工结构体
type Employee struct {
	ID         int     `db:"id"`
	Name       string  `db:"name"`
	Department string  `db:"department"`
	Salary     float64 `db:"salary"`
}

// Run 执行sqlx查询示例
func Run(db *gorm.DB) {
	fmt.Println("\n========== 开始执行 Sqlx 查询示例 ==========")

	// 从 GORM DB 获取底层的 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("获取 sql.DB 失败: %v\n", err)
		return
	}

	// 使用 sqlx 包装 sql.DB
	sqlxDB := sqlx.NewDb(sqlDB, "mysql")

	// 确保 employees 表存在并初始化数据
	initEmployeesTable(sqlxDB)

	// 1. 查询技术部的所有员工
	queryTechDepartment(sqlxDB)

	// 2. 查询工资最高的员工
	queryHighestSalary(sqlxDB)

	fmt.Println("========== Sqlx 查询示例执行完毕 ==========")
}

// initEmployeesTable 初始化员工表和数据
func initEmployeesTable(db *sqlx.DB) {
	fmt.Println("\n--- 初始化 employees 表 ---")

	// 创建表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		department VARCHAR(100) NOT NULL,
		salary DECIMAL(10, 2) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建 employees 表失败: %v\n", err)
		return
	}

	// 清空表
	_, err = db.Exec("TRUNCATE TABLE employees")
	if err != nil {
		log.Printf("清空 employees 表失败: %v\n", err)
	}

	// 插入测试数据
	insertSQL := `
	INSERT INTO employees (name, department, salary) VALUES
		('张三', '技术部', 8000.00),
		('李四', '技术部', 9500.00),
		('王五', '销售部', 7000.00),
		('赵六', '技术部', 12000.00),
		('钱七', '人力资源部', 6500.00),
		('孙八', '技术部', 10500.00),
		('周九', '销售部', 8500.00);
	`
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Printf("插入测试数据失败: %v\n", err)
		return
	}

	fmt.Println("employees 表初始化成功，已插入测试数据")
}

// queryTechDepartment 查询技术部的所有员工
func queryTechDepartment(db *sqlx.DB) {
	fmt.Println("\n--- 查询技术部的所有员工 ---")

	var employees []Employee
	query := "SELECT id, name, department, salary FROM employees WHERE department = ?"
	err := db.Select(&employees, query, "技术部")
	if err != nil {
		log.Printf("查询技术部员工失败: %v\n", err)
		return
	}

	fmt.Printf("查询到 %d 名技术部员工:\n", len(employees))
	for _, emp := range employees {
		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
			emp.ID, emp.Name, emp.Department, emp.Salary)
	}
}

// queryHighestSalary 查询工资最高的员工
func queryHighestSalary(db *sqlx.DB) {
	fmt.Println("\n--- 查询工资最高的员工 ---")

	var employee Employee
	query := "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1"
	err := db.Get(&employee, query)
	if err != nil {
		log.Printf("查询工资最高的员工失败: %v\n", err)
		return
	}

	fmt.Printf("工资最高的员工:\nID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
		employee.ID, employee.Name, employee.Department, employee.Salary)
}
