package gormSqlTwo

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// 题目2：事务语句
// 假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）
// 和 transactions表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
// 要求 ：
// 编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
// 如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。

// Account 账户表
type Account struct {
	ID      uint    `gorm:"primaryKey;autoIncrement"`
	Balance float64 `gorm:"type:decimal(10,2)"`
}

// Transaction 交易记录表
type Transaction struct {
	ID            uint    `gorm:"primaryKey;autoIncrement"`
	FromAccountID uint    `gorm:"column:from_account_id"`
	ToAccountID   uint    `gorm:"column:to_account_id"`
	Amount        float64 `gorm:"type:decimal(10,2)"`
}

// Run 执行转账事务示例
func Run(db *gorm.DB) {
	// 自动迁移创建表
	db.AutoMigrate(&Account{}, &Transaction{})

	fmt.Println("=== 银行转账事务示例 ===")

	// 创建测试账户
	accountA := Account{Balance: 500.00} // 账户A初始余额500元
	accountB := Account{Balance: 300.00} // 账户B初始余额300元

	db.Create(&accountA)
	db.Create(&accountB)

	fmt.Printf("转账前账户余额:\n")
	fmt.Printf("账户A (ID: %d): %.2f 元\n", accountA.ID, accountA.Balance)
	fmt.Printf("账户B (ID: %d): %.2f 元\n", accountB.ID, accountB.Balance)

	// 执行转账事务：从账户A向账户B转账100元
	fmt.Println("\n执行转账事务：从账户A向账户B转账100元...")
	err := transferMoney(db, accountA.ID, accountB.ID, 100.00)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功!")
	}

	// 查询并显示转账后余额
	var updatedAccountA, updatedAccountB Account
	db.First(&updatedAccountA, accountA.ID)
	db.First(&updatedAccountB, accountB.ID)

	fmt.Printf("\n转账后账户余额:\n")
	fmt.Printf("账户A (ID: %d): %.2f 元\n", updatedAccountA.ID, updatedAccountA.Balance)
	fmt.Printf("账户B (ID: %d): %.2f 元\n", updatedAccountB.ID, updatedAccountB.Balance)

	// 显示交易记录
	var transactions []Transaction
	db.Find(&transactions)
	fmt.Println("\n交易记录:")
	for _, t := range transactions {
		fmt.Printf("ID: %d, 从账户%d向账户%d转账%.2f元\n",
			t.ID, t.FromAccountID, t.ToAccountID, t.Amount)
	}

	// 演示余额不足的情况
	fmt.Println("\n尝试从账户B向账户A转账500元（余额不足）...")
	err = transferMoney(db, updatedAccountB.ID, updatedAccountA.ID, 500.00)
	if err != nil {
		fmt.Printf("转账失败: %v\n", err)
	} else {
		fmt.Println("转账成功!")
	}
}

// transferMoney 执行转账事务
func transferMoney(db *gorm.DB, fromAccountID, toAccountID uint, amount float64) error {
	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 1. 查询转出账户余额
	var fromAccount Account
	if err := tx.First(&fromAccount, fromAccountID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询转出账户失败: %v", err)
	}

	// 2. 检查余额是否足够
	if fromAccount.Balance < amount {
		tx.Rollback()
		return errors.New("余额不足，无法完成转账")
	}

	// 3. 扣除转出账户余额
	if err := tx.Model(&fromAccount).Update("balance", fromAccount.Balance-amount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("扣除转出账户余额失败: %v", err)
	}

	// 4. 增加转入账户余额
	var toAccount Account
	if err := tx.First(&toAccount, toAccountID).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("查询转入账户失败: %v", err)
	}

	if err := tx.Model(&toAccount).Update("balance", toAccount.Balance+amount).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("增加转入账户余额失败: %v", err)
	}

	// 5. 记录交易信息
	transaction := Transaction{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        amount,
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("记录交易信息失败: %v", err)
	}

	// 提交事务
	return tx.Commit().Error
}
