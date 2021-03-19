package mysql

import (
	"fmt"
	"gorm.io/gorm"
)

// 分页返回数
type PageList struct {
	CurrentPage int64       `json:"current_page"`
	FirstPage   int64       `json:"first_page"`
	LastPage    int64       `json:"last_page"`
	PageSize    int64       `json:"page_size"`
	Total       int64       `json:"total"`
	Data        interface{} `json:"data"`
}

func SetPageList(data interface{}, currentPage int64, pageSize ...int64) *PageList {
	pageList := &PageList{CurrentPage: currentPage, FirstPage: 1, Data: data}
	if len(pageSize) > 0 {
		pageList.PageSize = pageSize[0]
	}
	return pageList
}

type Where struct {
	Way   string
	Value interface{}
}

type Orm struct {
	DB *gorm.DB
}

func GetInstance() *Orm {
	return &Orm{DB: DB}
}

// 设置查询关联
func (db *Orm) Relationship(relationship []string) *Orm {
	for _, v := range relationship {
		db.DB = db.DB.Preload(v)
	}
	return db
}

func (db *Orm) getQueryValues(where map[string]interface{}) (string, []interface{}) {
	query := ""
	var values []interface{}
	for k, v := range where {
		if w, ok := v.(Where); ok {
			if query == "" {
				query = fmt.Sprintf("%s %s ?", k, w.Way)
			} else {
				query = fmt.Sprintf("%s AND %s %s ?", query, k, w.Way)
			}
			values = append(values, w.Value)
		} else {
			if query == "" {
				query = fmt.Sprintf("%s = ?", k)
			} else {
				query = fmt.Sprintf("%s AND %s = ?", query, k)
			}
			values = append(values, v)
		}
	}
	return query, values
}

func (db *Orm) Where(where map[string]interface{}) *Orm {
	query, values := db.getQueryValues(where)
	db.DB = db.DB.Where(query, values...)
	return db
}

func (db *Orm) Or(where map[string]interface{}) *Orm {
	query, values := db.getQueryValues(where)
	db.DB = db.DB.Or(query, values...)
	return db
}

// 设置查询分页信息
func (db *Orm) Limit(offset int, limit int) *Orm {
	db.DB = db.DB.Offset(offset).Limit(limit)
	return db
}

// 设置查询排序
func (db *Orm) Order(orderBy string) *Orm {
	db.DB = db.DB.Order(orderBy)
	return db
}

// 创建数据
func (db *Orm) Create(model interface{}) error {
	return db.DB.Create(model).Error
}

// 保存更新数据
func (db *Orm) Save(model interface{}) error {
	return db.DB.Save(model).Error
}

// 删除数据
func (db *Orm) Delete(model interface{}) error {
	return db.DB.Delete(model).Error
}

// 主键查询一条数据
func (db *Orm) First(model interface{}, id interface{}, relationship ...string) error {
	return db.Relationship(relationship).DB.First(model, id).Error
}

// 查询多条数据
func (db *Orm) One(model interface{}, relationship ...string) error {
	return db.Relationship(relationship).DB.First(model).Error
}

// 查询多条数据
func (db *Orm) Get(model interface{}, relationship ...string) {
	db.Relationship(relationship).DB.Find(model)
}

// 查询分页数据
func (db *Orm) Paginate(lists *PageList, relationship ...string) {
	if lists.PageSize == 0 {
		lists.PageSize = 15
	}
	var count int64
	db.DB.Model(lists.Data).Count(&count)
	lists.LastPage = (count / lists.PageSize) + 1
	lists.Total = count
	if count > 0 && lists.LastPage >= lists.CurrentPage {
		offset := (lists.CurrentPage - 1) * lists.PageSize
		db.Relationship(relationship).Limit(int(offset), int(lists.PageSize)).Get(lists.Data)
	}
}
