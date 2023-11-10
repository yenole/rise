package mapper

import (
	"strings"

	"gorm.io/gorm"
)

type Mapper[T any] struct {
}

func (m *Mapper[T]) Get(db *gorm.DB, dst *T) (bool, error) {
	db = db.Model(new(T)).Where(dst).Find(dst)
	return db.RowsAffected > 0, db.Error
}

func (m *Mapper[T]) Delete(db *gorm.DB, dst *T) error {
	return db.Model(new(T)).Delete(dst).Error
}

func (m *Mapper[T]) filters(filters ...any) func(*gorm.DB) *gorm.DB {
	return func(d *gorm.DB) *gorm.DB {
		var wheres []string
		var values []any
		size := len(filters) / 2
		for i := 0; i < size; i++ {
			if colume, ok := filters[i*2].(string); ok {
				wheres = append(wheres, whereColume(colume))
				values = append(values, filters[i*2+1])
			}
		}
		if len(wheres) == 0 {
			return d
		}
		return d.Where(strings.Join(wheres, " AND "), values...)
	}
}

func (m *Mapper[T]) Find(db *gorm.DB, dist any, filters ...any) error {
	return db.Model(new(T)).Scopes(m.filters(filters...)).Find(dist).Error
}

func (m *Mapper[T]) Update(db *gorm.DB, value *T, columes ...any) error {
	if len(columes) == 0 {
		return db.Model(new(T)).Updates(value).Error
	}
	return db.Model(value).Select(columes[0], columes[1:]...).Updates(value).Error
}

func (m *Mapper[T]) Create(db *gorm.DB, value *T) error {
	return db.Create(value).Error
}

func whereColume(colume string) string {
	if strings.Contains(colume, "__") {
		list := strings.Split(colume, "__")
		return "`" + list[0] + "` " + list[1]
	}
	return "`" + colume + "` = ?"
}
