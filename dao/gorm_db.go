package dao

import (
	"gorm.io/gorm"
)

var defaultDB *gorm.DB

func SetDefaultDB(option *GormOption) {
	if defaultDB == nil {
		defaultDB = NewGormWithOption(option)
		return
	}
	panic("不能重复初始化")
}

func GetDefaultDB() *gorm.DB {
	return defaultDB
}
