package dao

import (
	"strings"

	"gorm.io/gorm"
)

const batchMaxNum = 1000 // 每批次最大数量

type AppendValuesFunc func(values []interface{}, obj interface{}) []interface{}

func BatchInsert(db *gorm.DB, insertSql string, objs []interface{}, appendValuesFunc AppendValuesFunc) error {
	var batches [][]interface{} // 批次数组
	var values []interface{}    // 值数组
	if len(objs) == 0 {
		return nil
	}

	//	计算一次有多少个值，组装从句
	numOfValue := len(appendValuesFunc(values, objs[0]))
	valuesClause := "("
	for i := 0; i < numOfValue; i++ {
		valuesClause += "?,"
	}
	valuesClause = valuesClause[0:len(valuesClause)-1] + ")"

	// 组装批次
	counter := 0
	for _, obj := range objs {
		values = appendValuesFunc(values, obj)
		counter++
		if counter >= batchMaxNum {
			batches = append(batches, values)
			counter = 0
			values = []interface{}{}
		}
	}
	if len(values) != 0 {
		batches = append(batches, values)
	}

	// 分批次执行
	for _, batch := range batches {
		length := len(batch) / numOfValue
		clause := make([]string, length)
		for i := 0; i < length; i++ {
			clause[i] = valuesClause
		}
		err := db.Exec(insertSql+strings.Join(clause, ","), batch...).Error
		if err != nil {
			return err
		}
	}
	return nil
}
