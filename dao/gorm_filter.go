package dao

import (
	"reflect"

	"gorm.io/gorm"
)

type whereArg struct {
	arg interface{}
}
type whereArgs struct {
	args []interface{}
}
type whereMinMax struct {
	min interface{}
	max interface{}
}

// 条件 用于 where
type WhereLike whereArg
type WhereEQ whereArg
type WhereNE whereArg
type WhereGT whereArg
type WhereGE whereArg
type WhereLT whereArg
type WhereLE whereArg
type WhereIn whereArgs
type WhereNotIn whereArgs
type WhereBetween whereMinMax

func ConditionLike(arg interface{}) WhereLike {
	return WhereLike{
		arg: arg,
	}
}

func ConditionEQ(arg interface{}) WhereEQ {
	return WhereEQ{
		arg: arg,
	}
}

func ConditionNE(arg interface{}) WhereNE {
	return WhereNE{
		arg: arg,
	}
}

func ConditionGT(arg interface{}) WhereGT {
	return WhereGT{
		arg: arg,
	}
}

func ConditionGE(arg interface{}) WhereGE {
	return WhereGE{
		arg: arg,
	}
}

func ConditionLT(arg interface{}) WhereLT {
	return WhereLT{
		arg: arg,
	}
}

func ConditionLE(arg interface{}) WhereLE {
	return WhereLE{
		arg: arg,
	}
}

func ConditionIn(args ...interface{}) WhereIn {
	return WhereIn{
		args: args,
	}
}

func ConditionNotIn(args ...interface{}) WhereNotIn {
	return WhereNotIn{
		args: args,
	}
}

func ConditionBetween(min interface{}, max interface{}) WhereBetween {
	return WhereBetween{
		min: min,
		max: max,
	}
}

type Filter = map[string]interface{}
type ListFilter = map[string][]interface{}

func DealUpdateOrDeleteFilter(db *gorm.DB, filter Filter) *gorm.DB {
	if filter == nil || len(filter) == 0 {
		db = db.Where("1 = 1")
		return db
	}
	return DealFilter(db, filter)
}

func DealFilter(db *gorm.DB, filter Filter) *gorm.DB {
	for k, v := range filter {
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Slice, reflect.Array:
			switch t.Name() {
			default:
				db = db.Where(k+" IN (?)", v)
			}
		default:
			switch t.Name() {
			case "WhereLike":
				s := v.(WhereLike)
				db = db.Where(k+" LIKE ?", s.arg)
			case "WhereEQ":
				s := v.(WhereEQ)
				db = db.Where(k+" = ?", s.arg)
			case "WhereNE":
				s := v.(WhereNE)
				db = db.Where(k+" != ?", s.arg)
			case "WhereGT":
				s := v.(WhereGT)
				db = db.Where(k+" > ?", s.arg)
			case "WhereGE":
				s := v.(WhereGE)
				db = db.Where(k+" >= ?", s.arg)
			case "WhereLT":
				s := v.(WhereLT)
				db = db.Where(k+" < ?", s.arg)
			case "WhereLE":
				s := v.(WhereLE)
				db = db.Where(k+" <= ?", s.arg)
			case "WhereIn":
				s := v.(WhereIn)
				db = db.Where(k+" IN (?)", s.args)
			case "WhereNotIn":
				s := v.(WhereNotIn)
				db = db.Where(k+" NOT IN (?)", s.args)
			case "WhereBetween":
				s := v.(WhereBetween)
				db = db.Where(k+" BETWEEN ? AND ?", s.min, s.max)

			default:
				db = db.Where(k+" = ?", v)
			}

		}
	}
	return db
}

type Paging struct {
	Page       int64
	PageSize   int64
	TotalCount int64
}

// 把分页信息加入到查询条件中，如果p = nil则什么都不做
// 如果page跟pagesize都等于0，也什么都不做
func SetPage(db *gorm.DB, p *Paging) *gorm.DB {
	if p == nil {
		return db
	}
	if p.Page == 0 && p.PageSize == 0 {
		return db
	}
	db.Count(&p.TotalCount)
	return db.Offset(int((p.Page - 1) * p.PageSize)).Limit(int(p.PageSize))
}
