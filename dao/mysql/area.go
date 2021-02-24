package mysql

import (
	"context"

	"gorm.io/gorm"
	"github.com/flamefatex/model/dao"
)

const (
	areaSqlInsert = "INSERT INTO `area` (" +
		//"`id`," + // 若该实体是同步来的，需要加上id
		"`area_code`," +
		"`area_cnname`," +
		"`area_enname`," +
		"`area_abbr`," +
		"`level`," +
		"`parent_id`," +
		"`parent_code`," +
		"`type`," +
		"`creator_id`," +
		"`creator_name`," +
		"`editor_id`," +
		"`editor_name`," +
		"`delete_status`," +
		"`delete_time`) " +
		"VALUES"
)

type areaRepo struct {
	db *gorm.DB
}

func NewAreaRepo(db *gorm.DB) *areaRepo {
	return &areaRepo{
		db: db,
	}
}

func (n *areaRepo) Get(ctx context.Context, id int64) (*dao.Area, error) {
	op := "Get"
	db := dao.DealContext(ctx, n.db, op)
	model := &dao.Area{}
	err := db.Where("id = ?", id).First(model).Error
	if err != nil {
	    return nil, err
	}
	return model, err
}

func (n *areaRepo) Save(ctx context.Context, area *dao.Area) error {
	op := "Save"
	db := dao.DealContext(ctx, n.db, op)
	return db.Save(area).Error
}

func (n *areaRepo) Create(ctx context.Context, area *dao.Area) error {
	op := "Create"
	db := dao.DealContext(ctx, n.db, op)
	return db.Create(area).Error
}

func (n *areaRepo) Update(ctx context.Context, filter dao.Filter, params map[string]interface{}) error {
	op := "Update"
	db := dao.DealContext(ctx, n.db, op)
	db = db.Model(&dao.Area{})
	db = dao.DealUpdateOrDeleteFilter(db, filter)
	return db.Updates(params).Error
}

func (n *areaRepo) Delete(ctx context.Context, filter dao.Filter) error {
	op := "Delete"
	db := dao.DealContext(ctx, n.db, op)
	db = dao.DealUpdateOrDeleteFilter(db, filter)
	return db.Delete(&dao.Area{}).Error
}

func (n *areaRepo) List(ctx context.Context, filter dao.Filter, paging *dao.Paging) (rs []*dao.Area, err error) {
	op := "List"
	db := dao.DealContext(ctx, n.db, op)
	db = db.Model(&dao.Area{})
	db = dao.DealFilter(db, filter)
	if paging != nil {
		db = dao.SetPage(db, paging)
	}
	err = db.Find(&rs).Error
	return rs, err
}

func (n *areaRepo) Last(ctx context.Context, filter dao.Filter) (*dao.Area, error) {
	op := "Last"
	db := dao.DealContext(ctx, n.db, op)
	model := &dao.Area{}
	db = dao.DealFilter(db, filter)
	err := db.Last(model).Error
	if err != nil {
	    return nil, err
	}
	return model, err
}

func (n *areaRepo) Count(ctx context.Context, filter dao.Filter) (int64, error) {
	op := "Count"
	db := dao.DealContext(ctx, n.db, op)
	db = db.Model(&dao.Area{})
	db = dao.DealFilter(db, filter)
	var count int64
	err := db.Count(&count).Error
	return count, err
}

func (n *areaRepo) BatchInsert(ctx context.Context, areas []*dao.Area) error {
	op := "BatchInsert"
    db := dao.DealContext(ctx, n.db, op)
	var objs []interface{}
	for _, obj := range areas {
		objs = append(objs, obj)
	}
	return dao.BatchInsert(db, areaSqlInsert,
		objs, func(values []interface{}, obj interface{}) []interface{} {
			v := obj.(*dao.Area)
			values = append(values,
				// v.Id, // 若该实体是同步来的，需要加上id
				v.AreaCode, v.AreaCnname, v.AreaEnname, v.AreaAbbr, v.Level, v.ParentId, v.ParentCode, v.Type, v.CreatorId, v.CreatorName, v.EditorId, v.EditorName, v.DeleteStatus, v.DeleteTime)
			return values
		})
}

func (n *areaRepo) BatchCreate(ctx context.Context, areas []*dao.Area) error {
	op := "BatchCreate"
	db := dao.DealContext(ctx, n.db, op)
	return db.Create(areas).Error
}


func (n *areaRepo) BatchSave(ctx context.Context, areas []*dao.Area) error {
	if len(areas) == 0 {
	    return nil
	}
	op := "BatchSave"
	db := dao.DealContext(ctx, n.db, op)
	return db.Save(areas).Error
}
