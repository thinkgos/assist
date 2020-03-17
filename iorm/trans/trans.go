// trans gorm 事务封装
package trans

import (
	"context"
	"errors"

	"github.com/jinzhu/gorm"
	icontext "github.com/thinkgos/assist/iorm/context"
)

// TransFunc 定义事务执行函数
type TransFunc func(context.Context) error

// Trans 事务管理
type Trans struct {
	db *gorm.DB
}

// NewTrans 创建事务管理实例
func NewTrans(db *gorm.DB) *Trans {
	return &Trans{db}
}

// Begin 开启事务,返回事务句柄
func (a *Trans) Begin() (interface{}, error) {
	tx := a.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}

// Commit 提交事务
func (a *Trans) Commit(trans interface{}) error {
	tx, ok := trans.(*gorm.DB)
	if !ok {
		return errors.New("unknown trans")
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

// Rollback 回滚事务
func (a *Trans) Rollback(trans interface{}) error {
	tx, ok := trans.(*gorm.DB)
	if !ok {
		return errors.New("unknown trans")
	}
	if err := tx.Rollback().Error; err != nil {
		return err
	}
	return nil
}

// ExecTrans 执行事务
func ExecTrans(ctx context.Context, db *gorm.DB, cb TransFunc) error {
	if trans := icontext.FromTrans(ctx); trans != nil {
		return cb(ctx)
	}

	transModel := NewTrans(db)
	trans, err := transModel.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = transModel.Rollback(trans)
			panic(r)
		}
	}()

	ctx = icontext.NewTrans(ctx, trans)
	err = cb(ctx)
	if err != nil {
		_ = transModel.Rollback(trans)
		return err
	}
	return transModel.Commit(trans)
}

// ExecTransWithLock 执行事务（加锁）
func ExecTransWithLock(ctx context.Context, db *gorm.DB, cb TransFunc) error {
	if !icontext.FromTransLock(ctx) {
		ctx = icontext.NewTransLock(ctx)
	}
	return ExecTrans(ctx, db, cb)
}