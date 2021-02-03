package xormhelper

import (
	xorm "github.com/laixyz/xormplus"
	"math"
)

type Model struct {
	session *xorm.Session `xorm:"-" json:"-"`
	orderBy string        `xorm:"-" json:"-"`
}

func (obj *Model) SessionInit(session *xorm.Session, tableNameOrBean interface{}) {
	session.Table(tableNameOrBean)
	obj.session = session
}

func (obj *Model) Where(query interface{}, args ...interface{}) *Model {
	obj.session = obj.session.Where(query, args...)
	return obj
}

func (obj *Model) And(column string, args ...interface{}) *Model {
	obj.session = obj.session.And(column, args...)
	return obj
}

func (obj *Model) SQL(query interface{}, args ...interface{}) *Model {
	obj.session = obj.session.SQL(query, args...)
	return obj
}

func (obj *Model) Or(column string, args ...interface{}) *Model {
	obj.session = obj.session.Or(column, args...)
	return obj
}

func (obj *Model) Cols(columns ...string) *Model {
	obj.session = obj.session.Cols(columns...)
	return obj
}

func (obj *Model) AllCols() *Model {
	obj.session = obj.session.AllCols()
	return obj
}

func (obj *Model) Asc(colNames ...string) *Model {
	obj.session = obj.session.Asc(colNames...)
	obj.orderBy = obj.session.GetOrderBy()
	return obj
}

func (obj *Model) Desc(colNames ...string) *Model {
	obj.session = obj.session.Desc(colNames...)
	obj.orderBy = obj.session.GetOrderBy()
	return obj
}

func (obj *Model) OrderBy(order string) *Model {
	obj.session = obj.session.OrderBy(order)
	obj.orderBy = obj.session.GetOrderBy()
	return obj
}

func (obj *Model) GroupBy(keys string) *Model {
	obj.session = obj.session.GroupBy(keys)
	return obj
}

func (obj *Model) Distinct(columns ...string) *Model {
	obj.session = obj.session.Distinct(columns...)
	return obj
}

func (obj *Model) Having(conditions string) *Model {
	obj.session = obj.session.Having(conditions)
	return obj
}

func (obj *Model) IN(column string, args ...interface{}) *Model {
	obj.session = obj.session.In(column, args...)
	return obj
}

func (obj *Model) NotIn(column string, args ...interface{}) *Model {
	obj.session = obj.session.NotIn(column, args...)
	return obj
}

func (obj *Model) Limit(limit int, start ...int) *Model {
	obj.session = obj.session.Limit(limit, start...)
	return obj
}

func (obj *Model) FilterID(id interface{}) *Model {
	obj.session = obj.session.ID(id)
	return obj
}

func (obj *Model) Update() (int64, error) {
	return obj.session.Update(obj)
}

func (obj *Model) LastSQL() (string, []interface{}) {
	return obj.session.LastSQL()
}

func (obj *Model) Save() (int64, error) {
	return obj.session.Insert(obj)
}

func (obj *Model) Delete(bean interface{}) (int64, error) {
	return obj.session.Delete(bean)
}

func (obj *Model) Exist(bean ...interface{}) (bool, error) {
	return obj.session.Exist(bean...)
}

func (obj *Model) Count(bean ...interface{}) (int64, error) {
	return obj.session.Count(bean...)
}

func (obj *Model) FindOne(bean interface{}) (bool, error) {
	if obj.orderBy != "" {
		obj.session.OrderBy(obj.orderBy)
	}
	return obj.session.Get(bean)
}

func (obj *Model) FindAll(rowsSlicePtr interface{}, condiBean ...interface{}) error {
	if obj.orderBy != "" {
		obj.session.OrderBy(obj.orderBy)
	}
	return obj.session.Find(rowsSlicePtr, condiBean...)
}

func (obj *Model) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	if obj.orderBy != "" {
		obj.session.OrderBy(obj.orderBy)
	}
	return obj.session.FindAndCount(rowsSlicePtr, condiBean...)
}

func (obj *Model) Rows(bean interface{}) (*xorm.Rows, error) {
	return obj.session.Rows(bean)
}

func (obj *Model) Select(page, pagesize int, data interface{}) (currentPage, currentPagesize, totalRecords int64, totalPages int64, err error) {
	totalRecords, err = obj.session.Count()
	if err != nil {
		return
	}
	if totalRecords == 0 {
		return
	}
	if page < 1 {
		currentPage = 1
	} else {
		currentPage = int64(page)
	}
	if pagesize < 1 {
		pagesize = 1
	}
	currentPagesize = int64(pagesize)
	limit := pagesize
	start := (page - 1) * pagesize
	totalPages = int64(math.Ceil(float64(totalRecords) / float64(pagesize)))

	obj.session.Limit(limit, start)
	if obj.orderBy != "" {
		obj.session.OrderBy(obj.orderBy)
	}
	err = obj.session.Find(data)
	return currentPage, currentPagesize, totalRecords, totalPages, err
}
