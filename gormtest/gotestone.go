package gormtest

import (
	"database/sql"
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           // Standard field for the primary key
	Name         string         // A regular string field
	Email        *string        // A pointer to a string, allowing for null values
	Age          uint8          // An unsigned 8-bit integer
	Birthday     *time.Time     // A pointer to time.Time, can be null
	MemberNumber sql.NullString // Uses sql.NullString to handle nullable strings
	ActivatedAt  sql.NullTime   // Uses sql.NullTime for nullable time fields
	CreatedAt    time.Time      // Automatically managed by GORM for creation time
	UpdatedAt    time.Time      // Automatically managed by GORM for update time
	ignored      string         // fields that aren't exported are ignored
}
type Member struct {
	gorm.Model
	Name string
	Age  uint8
}
type Author struct {
	Name  string
	Email string
}

type Blog struct {
	Author
	ID      int
	Upvotes int32
}

func Run(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&User{}, &Member{}, &Blog{})

	// 创建一个用户
	user := User{}
	user.MemberNumber.Valid = true

	// 创建一个会员
	member := Member{}
	//数据库创建
	db.Create(&member)
	db.Create(&user)

}
