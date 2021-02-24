package dao

import (
	"context"
	"time"
)

type Area struct {
	// 区域id
	Id int64 `gorm:"primaryKey;"`
	// 区域编码
	AreaCode string
	// 区域中文名称
	AreaCnname string
	// 区域英文名称
	AreaEnname string
	// 区域缩写
	AreaAbbr string
	// 区域层级
	Level int32
	// 父级id
	ParentId int64
	// 父级code
	ParentCode string
	// 区域类型，1:普通区域 2:其他区域
	Type int32
	// 创建用户id
	CreatorId int64
	// 创建用户名称
	CreatorName string
	// 最后编辑用户id
	EditorId int64
	// 最后编辑用户名称
	EditorName string
	// 软删除状态，1:停用 2:启用
	DeleteStatus int32
	// 删除时间
	DeleteTime *time.Time

	CreateTime *time.Time `gorm:"autoCreateTime;<-:create;"`
	UpdateTime *time.Time `gorm:"autoUpdateTime;"`
}

type AreaId = int64

type AreaRepository interface {
	Get(ctx context.Context, id int64) (*Area, error)
	Save(ctx context.Context, area *Area) error
	Create(ctx context.Context, area *Area) error
	Update(ctx context.Context, filter Filter, params map[string]interface{}) error
	Delete(ctx context.Context, filter Filter) error
	List(ctx context.Context, filter Filter, paging *Paging) ([]*Area, error)
	Last(ctx context.Context, filter Filter) (*Area, error)
	Count(ctx context.Context, filter Filter) (int64, error)
	BatchInsert(ctx context.Context, areas []*Area) error
	BatchCreate(ctx context.Context, areas []*Area) error
	BatchSave(ctx context.Context, areas []*Area) error
}

func AreaExtractIdSet(list []*Area) map[AreaId]struct{} {
	rs := make(map[AreaId]struct{})
	for _, item := range list {
		rs[item.Id] = struct{}{}
	}
	return rs
}

func AreaExtractIdMap(list []*Area) map[AreaId]*Area {
	rs := make(map[AreaId]*Area)
	for _, item := range list {
		rs[item.Id] = item
	}
	return rs
}
