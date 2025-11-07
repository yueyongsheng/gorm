package sqlxtwo

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"gorm.io/gorm"
)

//假设有一个 books 表，包含字段 id 、 title 、 author 、 price 。
// 要求 ：
// 定义一个 Book 结构体，包含与 books 表对应的字段。
// 编写Go代码，使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全。

// Book 书籍结构体
type Book struct {
	ID     int     `db:"id"`
	Title  string  `db:"title"`
	Author string  `db:"author"`
	Price  float64 `db:"price"`
}

// Run 执行sqlx查询示例
func Run(db *gorm.DB) {
	fmt.Println("\n========== 开始执行 Sqlx Books 查询示例 ==========")

	// 从 GORM DB 获取底层的 *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		log.Printf("获取 sql.DB 失败: %v\n", err)
		return
	}

	// 使用 sqlx 包装 sql.DB
	sqlxDB := sqlx.NewDb(sqlDB, "mysql")

	// 确保 books 表存在并初始化数据
	initBooksTable(sqlxDB)

	// 查询价格大于 50 元的书籍
	queryExpensiveBooks(sqlxDB, 50.0)
}

// initBooksTable 初始化书籍表和数据
func initBooksTable(db *sqlx.DB) {
	fmt.Println("\n--- 初始化 books 表 ---")

	// 创建表
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(200) NOT NULL,
		author VARCHAR(100) NOT NULL,
		price DECIMAL(10, 2) NOT NULL
	) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
	`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Printf("创建 books 表失败: %v\n", err)
		return
	}

	// 清空表
	_, err = db.Exec("TRUNCATE TABLE books")
	if err != nil {
		log.Printf("清空 books 表失败: %v\n", err)
	}

	// 插入测试数据
	insertSQL := `
	INSERT INTO books (title, author, price) VALUES
		('射雕英雄传', '金庸', 68.00),
		('天龙八部', '金庸', 78.50),
		('三体', '刘慈欣', 45.00),
		('活着', '余华', 32.00),
		('白夜行', '东野圭吾', 55.00),
		('解忧杂货店', '东野圭吾', 42.00),
		('红楼梦', '曹雪芹', 89.00),
		('平凡的世界', '路遥', 38.00),
		('百年孤独', '加西亚·马尔克斯', 52.00),
		('小王子', '圣埃克苏佩里', 28.00);
	`
	_, err = db.Exec(insertSQL)
	if err != nil {
		log.Printf("插入测试数据失败: %v\n", err)
		return
	}

	fmt.Println("books 表初始化成功，已插入测试数据")
}

// queryExpensiveBooks 查询价格大于指定金额的书籍
func queryExpensiveBooks(db *sqlx.DB, minPrice float64) {
	fmt.Printf("\n--- 查询价格大于 %.2f 元的书籍 ---\n", minPrice)

	var books []Book
	query := "SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC"
	err := db.Select(&books, query, minPrice)
	if err != nil {
		log.Printf("查询书籍失败: %v\n", err)
		return
	}

	fmt.Printf("查询到 %d 本价格大于 %.2f 元的书籍:\n", len(books), minPrice)
	for _, book := range books {
		fmt.Printf("ID: %d, 书名: %s, 作者: %s, 价格: %.2f 元\n",
			book.ID, book.Title, book.Author, book.Price)
	}
}
