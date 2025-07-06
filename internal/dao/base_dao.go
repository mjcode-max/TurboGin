package dao

import (
	"gorm.io/gorm"
)

// IBaseDAO 泛型 CRUD 接口
type IBaseDAO[T any] interface {
	Create(entity *T) error
	GetByID(id uint) (*T, error)
	Update(entity *T) error
	Delete(id uint) error
	Find(conditions interface{}, args ...interface{}) ([]T, error)
	DB() *gorm.DB
}

// BaseDAO 泛型 CRUD 实现
type BaseDAO[T any] struct {
	db *gorm.DB
}

func NewBaseDAO[T any](db *gorm.DB) IBaseDAO[T] {
	return &BaseDAO[T]{db: db}
}

func (d *BaseDAO[T]) Create(entity *T) error {
	return d.db.Create(entity).Error
}

func (d *BaseDAO[T]) GetByID(id uint) (*T, error) {
	var entity T
	if err := d.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (d *BaseDAO[T]) Update(entity *T) error {
	return d.db.Save(entity).Error
}

func (d *BaseDAO[T]) Delete(id uint) error {
	var entity T
	return d.db.Delete(&entity, id).Error
}

func (d *BaseDAO[T]) Find(conditions interface{}, args ...interface{}) ([]T, error) {
	var entities []T
	err := d.db.Where(conditions, args...).Find(&entities).Error
	return entities, err
}

func (d *BaseDAO[T]) DB() *gorm.DB {
	return d.db.Model(new(T))
}
