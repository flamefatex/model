package dao

import (
	"context"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const saveAssociationsKey = "saveAssociations"

var isEnableLog = false

func EnableLog() {
	isEnableLog = true
}

func DisableLog() {
	isEnableLog = false
}

// GormOption gorm 参数
type GormOption struct {
	DriverName                string
	DSN                       string
	SkipInitializeWithVersion bool
	IsEnableLog               bool
}

func NewGorm(dsn string) (db *gorm.DB) {
	option := &GormOption{
		DriverName:                "mysql",
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
		IsEnableLog:               isEnableLog,
	}
	return NewGormWithOption(option)
}

func NewGormWithOption(option *GormOption) (db *gorm.DB) {
	logMode := logger.Silent
	if option.IsEnableLog {
		logMode = logger.Info
	}

	dialector := mysql.New(mysql.Config{
		DriverName:                option.DriverName,
		DSN:                       option.DSN,
		SkipInitializeWithVersion: option.SkipInitializeWithVersion,
	})

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logMode),

		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	}

	db, err := gorm.Open(dialector, config)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetConnMaxLifetime(time.Second * 300)
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(1000)
	return
}

func NewMockGorm() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, _ := sqlmock.New()
	dialector := mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	})
	config := &gorm.Config{
		DisableAutomaticPing: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
	}
	orm, err := gorm.Open(dialector, config)
	return orm, mock, err
}

type txKey struct{}

func NewTxContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, txKey{}, db)
}

// TxFromContext 获取事务连接
func TxFromContext(ctx context.Context) *gorm.DB {
	val := ctx.Value(txKey{})
	if val == nil {
		return nil
	}
	return val.(*gorm.DB)
}

// CurrentDB 如果context里面有事务连接则使用事务连接，否则使用db
func CurrentDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	tx := TxFromContext(ctx)
	if tx != nil {
		return tx
	} else {
		return db
	}
}

type optsKey struct{}

type ORMOption func(db *gorm.DB) *gorm.DB

func NewORMOptsContext(ctx context.Context, opts ...ORMOption) context.Context {
	return context.WithValue(ctx, optsKey{}, opts)
}

// ORMOptsFromContext 获取orm的options
func ORMOptsFromContext(ctx context.Context) []ORMOption {
	val := ctx.Value(optsKey{})
	if val == nil {
		return nil
	}
	return val.([]ORMOption)
}

func SaveAssociations(db *gorm.DB) *gorm.DB {
	return db.Set(saveAssociationsKey, true)
}

func OrderByIdDesc(db *gorm.DB) *gorm.DB {
	return db.Order("id desc")
}


func Preload(query string) ORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(query)
	}
}

func OrderBy(value string) ORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(value)
	}
}


func DealContext(ctx context.Context, db *gorm.DB, op string) *gorm.DB {
	db = db.Set("ctx", ctx)

	// 可能事务
	db = CurrentDB(ctx, db)

	// 额外选项 一般是查询预加载、排序，或关联保存
	for _, v := range ORMOptsFromContext(ctx) {
		db = v(db)
	}

	// 默认不要关联保存
	switch op {
	case "Create", "Save", "BatchCreate", "BatchSave":
		_, ok := db.Get(saveAssociationsKey)
		// 没有设置关联保存
		if !ok {
			db = db.Omit(clause.Associations)
		}
	}

	return db
}
