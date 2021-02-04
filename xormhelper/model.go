package xormhelper

import (
	xorm "github.com/laixyz/xormplus"
	"math"
)

type Model struct {
	session *xorm.Session `xorm:"-" json:"-"`
	orderBy string        `xorm:"-" json:"-"`
}

func (m *Model) SessionInit(session *xorm.Session, tableNameOrBean interface{}) {
	session.Table(tableNameOrBean)
	m.session = session
}

func (m *Model) Where(query interface{}, args ...interface{}) *Model {
	m.session = m.session.Where(query, args...)
	return m
}

func (m *Model) And(column string, args ...interface{}) *Model {
	m.session = m.session.And(column, args...)
	return m
}

func (m *Model) SQL(query interface{}, args ...interface{}) *Model {
	m.session = m.session.SQL(query, args...)
	return m
}

func (m *Model) Or(column string, args ...interface{}) *Model {
	m.session = m.session.Or(column, args...)
	return m
}

func (m *Model) Cols(columns ...string) *Model {
	m.session = m.session.Cols(columns...)
	return m
}

func (m *Model) AllCols() *Model {
	m.session = m.session.AllCols()
	return m
}

func (m *Model) Asc(colNames ...string) *Model {
	m.session = m.session.Asc(colNames...)
	m.orderBy = m.session.GetOrderBy()
	return m
}

func (m *Model) Desc(colNames ...string) *Model {
	m.session = m.session.Desc(colNames...)
	m.orderBy = m.session.GetOrderBy()
	return m
}

func (m *Model) OrderBy(order string) *Model {
	m.session = m.session.OrderBy(order)
	m.orderBy = m.session.GetOrderBy()
	return m
}

func (m *Model) GroupBy(keys string) *Model {
	m.session = m.session.GroupBy(keys)
	return m
}

func (m *Model) Distinct(columns ...string) *Model {
	m.session = m.session.Distinct(columns...)
	return m
}

func (m *Model) Having(conditions string) *Model {
	m.session = m.session.Having(conditions)
	return m
}

func (m *Model) IN(column string, args ...interface{}) *Model {
	m.session = m.session.In(column, args...)
	return m
}

func (m *Model) NotIn(column string, args ...interface{}) *Model {
	m.session = m.session.NotIn(column, args...)
	return m
}

func (m *Model) Limit(limit int, start ...int) *Model {
	m.session = m.session.Limit(limit, start...)
	return m
}

func (m *Model) FilterID(id interface{}) *Model {
	m.session = m.session.ID(id)
	return m
}

func (m *Model) Update() (int64, error) {
	return m.session.Update(m)
}

func (m *Model) LastSQL() (string, []interface{}) {
	return m.session.LastSQL()
}

func (m *Model) Save() (int64, error) {
	return m.session.Insert(m)
}

func (m *Model) Delete(bean interface{}) (int64, error) {
	return m.session.Delete(bean)
}

func (m *Model) Exist(bean ...interface{}) (bool, error) {
	return m.session.Exist(bean...)
}

func (m *Model) Count(bean ...interface{}) (int64, error) {
	return m.session.Count(bean...)
}

func (m *Model) FindOne(bean interface{}) (bool, error) {
	if m.orderBy != "" {
		m.session.OrderBy(m.orderBy)
		defer func() {
			m.orderBy = ""
		}()
	}
	return m.session.Get(bean)
}

func (m *Model) FindAll(rowsSlicePtr interface{}, condiBean ...interface{}) error {
	if m.orderBy != "" {
		m.session.OrderBy(m.orderBy)
		defer func() {
			m.orderBy = ""
		}()
	}
	return m.session.Find(rowsSlicePtr, condiBean...)
}

func (m *Model) FindAndCount(rowsSlicePtr interface{}, condiBean ...interface{}) (int64, error) {
	if m.orderBy != "" {
		m.session.OrderBy(m.orderBy)
		defer func() {
			m.orderBy = ""
		}()
	}

	return m.session.FindAndCount(rowsSlicePtr, condiBean...)
}

func (m *Model) Rows(bean interface{}) (*xorm.Rows, error) {
	if m.orderBy != "" {
		m.session.OrderBy(m.orderBy)
		defer func() {
			m.orderBy = ""
		}()
	}
	return m.session.Rows(bean)
}

func (m *Model) Pagination(page, pagesize int, data interface{}) (currentPage, currentPagesize, totalRecords, totalPages int, err error) {
	total, err := m.session.Count()
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	if page < 1 {
		currentPage = 1
	} else {
		currentPage = page
	}
	if pagesize < 1 {
		pagesize = 1
	}
	totalRecords = int(total)
	currentPagesize = pagesize
	limit := pagesize
	start := (page - 1) * pagesize
	totalPages = int(math.Ceil(float64(totalRecords) / float64(pagesize)))
	if totalRecords > pagesize {
		m.session.Limit(limit, start)
	}
	if m.orderBy != "" {
		m.session.OrderBy(m.orderBy)
		defer func() {
			m.orderBy = ""
		}()
	}
	err = m.session.Find(data)
	return currentPage, currentPagesize, totalRecords, totalPages, err
}
