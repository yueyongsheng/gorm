package advancedone

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

// ============================================
// 题目1：模型定义
// ============================================
// 假设你要开发一个博客系统，有以下几个实体： User （用户）、 Post （文章）、 Comment （评论）。
// 要求 ：
// 使用Gorm定义 User 、 Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）。
// 编写Go代码，使用Gorm创建这些模型对应的数据库表

// ============================================
// 题目2：关联查询
// ============================================
// 基于上述博客系统的模型定义。
// 要求 ：
// 编写Go代码，使用Gorm查询某个用户发布的所有文章及其对应的评论信息。
// 编写Go代码，使用Gorm查询评论数量最多的文章信息。

// ============================================
// 题目3：钩子函数
// ============================================
// 继续使用博客系统的模型。
// 要求 ：
// 为 Post 模型添加一个钩子函数，在文章创建时自动更新用户的文章数量统计字段。
// 为 Comment 模型添加一个钩子函数，在评论删除时检查文章的评论数量，如果评论数量为 0，则更新文章的评论状态为 "无评论"。

// ============================================
// 模型定义
// ============================================

// User 用户模型
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Email     string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"` // json:"-" 表示不序列化密码
	Nickname  string    `gorm:"type:varchar(50)" json:"nickname"`
	PostCount int       `gorm:"default:0" json:"post_count"` // 题目3：文章数量统计字段
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 一对多关系：一个用户可以发布多篇文章
	Posts []Post `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title         string    `gorm:"type:varchar(200);not null;index" json:"title"`
	Content       string    `gorm:"type:text;not null" json:"content"`
	UserID        uint      `gorm:"not null;index" json:"user_id"` // 外键：关联用户
	ViewCount     int       `gorm:"default:0" json:"view_count"`
	CommentStatus string    `gorm:"type:varchar(20);default:'有评论'" json:"comment_status"` // 题目3：评论状态
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 多对一关系：多篇文章属于一个用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`

	// 一对多关系：一篇文章可以有多个评论
	Comments []Comment `gorm:"foreignKey:PostID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

// Comment 评论模型
type Comment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	PostID    uint      `gorm:"not null;index" json:"post_id"` // 外键：关联文章
	UserID    uint      `gorm:"not null;index" json:"user_id"` // 外键：关联用户（评论者）
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 多对一关系：多个评论属于一篇文章
	Post Post `gorm:"foreignKey:PostID" json:"post,omitempty"`

	// 多对一关系：多个评论属于一个用户
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// ============================================
// 题目3：钩子函数
// ============================================

// BeforeCreate Post 创建前钩子：自动更新用户的文章数量
func (p *Post) BeforeCreate(tx *gorm.DB) error {
	// 更新用户的文章数量
	if err := tx.Model(&User{}).Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error; err != nil {
		log.Printf("更新用户文章数量失败: %v", err)
		return err
	}
	log.Printf("钩子函数: 用户 ID=%d 的文章数量 +1", p.UserID)
	return nil
}

// AfterDelete Comment 删除后钩子：检查文章评论数量，更新评论状态
func (c *Comment) AfterDelete(tx *gorm.DB) error {
	// 查询该文章的评论数量
	var count int64
	if err := tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count).Error; err != nil {
		log.Printf("查询评论数量失败: %v", err)
		return err
	}

	// 如果评论数量为 0，更新文章的评论状态为 "无评论"
	if count == 0 {
		if err := tx.Model(&Post{}).Where("id = ?", c.PostID).
			Update("comment_status", "无评论").Error; err != nil {
			log.Printf("更新文章评论状态失败: %v", err)
			return err
		}
		log.Printf("钩子函数: 文章 ID=%d 的评论状态更新为 '无评论'", c.PostID)
	}

	return nil
}

// ============================================
// TableName 方法：自定义表名（可选）
// ============================================

func (User) TableName() string {
	return "users"
}

func (Post) TableName() string {
	return "posts"
}

func (Comment) TableName() string {
	return "comments"
}

// ============================================
// Run 执行博客系统模型创建
// ============================================

func Run(db *gorm.DB) {
	fmt.Println("\n========== 开始执行 GORM 高级模型定义示例 ==========")

	// 1. 自动迁移：创建表
	createTables(db)

	// 2. 插入测试数据
	insertTestData(db)

	// 3. 演示查询功能
	demonstrateQueries(db)

	fmt.Println("========== GORM 高级模型定义示例执行完毕 ==========\n")
}

// createTables 创建数据库表
func createTables(db *gorm.DB) {
	fmt.Println("\n--- 创建数据库表 ---")

	// 自动迁移（按顺序迁移，确保外键关系正确）
	err := db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Printf("自动迁移失败: %v\n", err)
		return
	}

	fmt.Println("数据库表创建成功：users, posts, comments")
}

// insertTestData 插入测试数据
func insertTestData(db *gorm.DB) {
	fmt.Println("\n--- 插入测试数据 ---")

	// 清空旧数据（按顺序删除，先删除依赖表）
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	// 创建用户
	user1 := User{
		Username: "alice",
		Email:    "alice@example.com",
		Password: "hashed_password_123",
		Nickname: "Alice",
	}

	user2 := User{
		Username: "bob",
		Email:    "bob@example.com",
		Password: "hashed_password_456",
		Nickname: "Bob",
	}

	if err := db.Create(&user1).Error; err != nil {
		log.Printf("创建用户1失败: %v\n", err)
		return
	}

	if err := db.Create(&user2).Error; err != nil {
		log.Printf("创建用户2失败: %v\n", err)
		return
	}

	fmt.Printf("创建用户成功: %s (ID: %d), %s (ID: %d)\n", user1.Username, user1.ID, user2.Username, user2.ID)

	// 创建文章
	post1 := Post{
		Title:   "GORM 入门教程",
		Content: "这是一篇关于 GORM 的入门教程，介绍了如何使用 GORM 进行数据库操作...",
		UserID:  user1.ID,
	}

	post2 := Post{
		Title:   "Go 语言最佳实践",
		Content: "本文分享了一些 Go 语言开发的最佳实践，包括代码组织、错误处理等...",
		UserID:  user1.ID,
	}

	post3 := Post{
		Title:   "微服务架构设计",
		Content: "探讨微服务架构的设计原则和实践经验...",
		UserID:  user2.ID,
	}

	if err := db.Create(&post1).Error; err != nil {
		log.Printf("创建文章1失败: %v\n", err)
		return
	}

	if err := db.Create(&post2).Error; err != nil {
		log.Printf("创建文章2失败: %v\n", err)
		return
	}

	if err := db.Create(&post3).Error; err != nil {
		log.Printf("创建文章3失败: %v\n", err)
		return
	}

	fmt.Printf("创建文章成功: %d 篇\n", 3)

	// 创建评论
	comment1 := Comment{
		Content: "这篇文章写得很好，学到了很多！",
		PostID:  post1.ID,
		UserID:  user2.ID, // Bob 评论 Alice 的文章
	}

	comment2 := Comment{
		Content: "感谢分享，期待更多内容！",
		PostID:  post1.ID,
		UserID:  user2.ID,
	}

	comment3 := Comment{
		Content: "赞同作者的观点！",
		PostID:  post2.ID,
		UserID:  user2.ID,
	}

	comment4 := Comment{
		Content: "非常实用的架构设计思路！",
		PostID:  post3.ID,
		UserID:  user1.ID, // Alice 评论 Bob 的文章
	}

	comments := []Comment{comment1, comment2, comment3, comment4}
	if err := db.Create(&comments).Error; err != nil {
		log.Printf("创建评论失败: %v\n", err)
		return
	}

	fmt.Printf("创建评论成功: %d 条\n", len(comments))
}

// demonstrateQueries 演示查询功能
func demonstrateQueries(db *gorm.DB) {
	fmt.Println("\n========================================")
	fmt.Println("题目2：关联查询")
	fmt.Println("========================================")

	// 题目2-1：查询某个用户发布的所有文章及其对应的评论信息
	queryUserPostsWithComments(db, "alice")

	// 题目2-2：查询评论数量最多的文章信息
	queryMostCommentedPost(db)

	fmt.Println("\n========================================")
	fmt.Println("题目3：钩子函数演示")
	fmt.Println("========================================")

	// 题目3-1：演示 Post 创建钩子（自动更新用户文章数量）
	demonstratePostCreateHook(db)

	// 题目3-2：演示 Comment 删除钩子（更新文章评论状态）
	demonstrateCommentDeleteHook(db)
}

// ============================================
// 题目2：关联查询实现
// ============================================

// queryUserPostsWithComments 查询某个用户发布的所有文章及其对应的评论信息
func queryUserPostsWithComments(db *gorm.DB, username string) {
	fmt.Printf("\n>>> 【题目2-1】查询用户 '%s' 发布的所有文章及其评论信息\n", username)

	var user User
	// 预加载文章，并预加载文章的评论和评论的用户信息
	if err := db.Preload("Posts.Comments.User").Where("username = ?", username).First(&user).Error; err != nil {
		log.Printf("查询失败: %v\n", err)
		return
	}

	fmt.Printf("用户: %s (昵称: %s)\n", user.Username, user.Nickname)
	fmt.Printf("文章数量: %d (PostCount字段: %d)\n", len(user.Posts), user.PostCount)
	fmt.Println("---")

	for i, post := range user.Posts {
		fmt.Printf("[文章 %d] %s\n", i+1, post.Title)
		fmt.Printf("  评论状态: %s\n", post.CommentStatus)
		fmt.Printf("  评论数量: %d\n", len(post.Comments))
		for j, comment := range post.Comments {
			fmt.Printf("    [评论 %d] %s: %s\n", j+1, comment.User.Username, comment.Content)
		}
		fmt.Println()
	}
}

// queryMostCommentedPost 查询评论数量最多的文章信息
func queryMostCommentedPost(db *gorm.DB) {
	fmt.Println("\n>>> 【题目2-2】查询评论数量最多的文章信息")

	// 方法1：使用子查询和 Group By
	type PostCommentCount struct {
		PostID       uint
		CommentCount int64
	}

	var result PostCommentCount
	err := db.Model(&Comment{}).
		Select("post_id, COUNT(*) as comment_count").
		Group("post_id").
		Order("comment_count DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		log.Printf("查询失败: %v\n", err)
		return
	}

	// 根据 PostID 查询文章详细信息
	var post Post
	if err := db.Preload("User").Preload("Comments.User").First(&post, result.PostID).Error; err != nil {
		log.Printf("查询文章详情失败: %v\n", err)
		return
	}

	fmt.Printf("评论数量最多的文章:\n")
	fmt.Printf("  文章ID: %d\n", post.ID)
	fmt.Printf("  标题: %s\n", post.Title)
	fmt.Printf("  作者: %s\n", post.User.Username)
	fmt.Printf("  评论数量: %d\n", len(post.Comments))
	fmt.Printf("  评论状态: %s\n", post.CommentStatus)
	fmt.Println("  所有评论:")
	for i, comment := range post.Comments {
		fmt.Printf("    [%d] %s: %s\n", i+1, comment.User.Username, comment.Content)
	}
}

// ============================================
// 题目3：钩子函数演示
// ============================================

// demonstratePostCreateHook 演示 Post 创建钩子
func demonstratePostCreateHook(db *gorm.DB) {
	fmt.Println("\n>>> 【题目3-1】演示 Post 创建钩子：自动更新用户文章数量")

	// 查询 alice 用户当前的文章数量
	var user User
	db.Where("username = ?", "alice").First(&user)
	fmt.Printf("创建前: alice 的文章数量 = %d\n", user.PostCount)

	// 创建新文章（会触发 BeforeCreate 钩子）
	newPost := Post{
		Title:   "钩子函数测试文章",
		Content: "这是一篇用来测试钩子函数的文章",
		UserID:  user.ID,
	}

	if err := db.Create(&newPost).Error; err != nil {
		log.Printf("创建文章失败: %v\n", err)
		return
	}

	// 重新查询用户，查看文章数量是否自动更新
	db.Where("username = ?", "alice").First(&user)
	fmt.Printf("创建后: alice 的文章数量 = %d (自动 +1)\n", user.PostCount)
	fmt.Printf("新文章ID: %d, 标题: %s\n", newPost.ID, newPost.Title)
}

// demonstrateCommentDeleteHook 演示 Comment 删除钩子
func demonstrateCommentDeleteHook(db *gorm.DB) {
	fmt.Println("\n>>> 【题目3-2】演示 Comment 删除钩子：更新文章评论状态")

	// 查找一篇只有一条评论的文章
	// 先创建一篇新文章
	var user User
	db.Where("username = ?", "bob").First(&user)

	testPost := Post{
		Title:   "测试评论删除钩子的文章",
		Content: "这篇文章用来测试评论删除钩子",
		UserID:  user.ID,
	}
	db.Create(&testPost)

	// 创建一条评论
	testComment := Comment{
		Content: "这是唯一的评论",
		PostID:  testPost.ID,
		UserID:  user.ID,
	}
	db.Create(&testComment)

	// 更新文章的评论状态为"有评论"
	db.Model(&testPost).Update("comment_status", "有评论")

	fmt.Printf("删除前: 文章 '%s' 的评论状态 = '%s'\n", testPost.Title, testPost.CommentStatus)

	// 查询评论数量
	var count int64
	db.Model(&Comment{}).Where("post_id = ?", testPost.ID).Count(&count)
	fmt.Printf("删除前: 该文章的评论数量 = %d\n", count)

	// 删除评论（会触发 AfterDelete 钩子）
	if err := db.Delete(&testComment).Error; err != nil {
		log.Printf("删除评论失败: %v\n", err)
		return
	}

	// 重新查询文章，查看评论状态是否自动更新
	db.First(&testPost, testPost.ID)
	db.Model(&Comment{}).Where("post_id = ?", testPost.ID).Count(&count)
	fmt.Printf("删除后: 文章 '%s' 的评论状态 = '%s' (自动更新)\n", testPost.Title, testPost.CommentStatus)
	fmt.Printf("删除后: 该文章的评论数量 = %d\n", count)
}
