package repository

import (
	"context"
	"fmt"
	"lms-bootcamp/internal/domain/dto"

	"gorm.io/gorm"
)

type GormRepository[T any] struct {
	db *gorm.DB
}

func NewGormRepository[T any](db *gorm.DB) *GormRepository[T] {
	return &GormRepository[T]{db: db}
}

func (r *GormRepository[T]) GetAll(ctx context.Context,table ...string) ([]T, error) {
	var items []T

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(fmt.Sprintf("%s.deleted_at IS NULL", table[0])).Find(&items).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil	
}

func (r *GormRepository[T]) Create(ctx context.Context, item *T) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *GormRepository[T]) BulkCreate(ctx context.Context, items []*T) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(items).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *GormRepository[T]) Update(ctx context.Context, id string, item *T) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(new(T)).Where("id = ? AND deleted_at IS NULL", id).Updates(item).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *GormRepository[T]) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(new(T)).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *GormRepository[T]) GetById(ctx context.Context, id string, table ...string) (*T, error) {
	var item T

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var idQuery string
		if len(table) == 0 || table[0] == "" {
			idQuery = "id = ? AND deleted_at IS NULL"
		} else {
			idQuery = fmt.Sprintf("%s.id = ? AND %s.deleted_at IS NULL", table[0], table[0])
		}

		if err := tx.Where(idQuery, id).First(&item).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *GormRepository[T]) GetByFilter(ctx context.Context, filter map[string]interface{}) ([]T, error) {
	var items []T

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(filter).Find(&items).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return items, nil	
}

func (q *GormRepository[T]) Paginate(ctx context.Context, page, perPage int) (*dto.PaginationResponse[T], error) {
    var items []T
    var total int64

    if err := q.db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
        return nil, err
    }

    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = int(total)
    }
    
    offset := (page - 1) * perPage
    if err := q.db.WithContext(ctx).Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&items).Error; err != nil {
        return nil, err
    }

    return dto.NewPaginationResponse(items, total, page, perPage), nil
}


func (q *GormRepository[T]) PaginateQuery(ctx context.Context, page, perPage int) (*dto.PaginationResponse[T], error) {
    var items []T
    var total int64
    
    // Clone db untuk count
    countDB := q.db.Session(&gorm.Session{})
    if err := countDB.WithContext(ctx).Model(new(T)).Where("deleted_at IS NULL").Count(&total).Error; err != nil {
        return nil, err
    }
    
    if page < 1 {
        page = 1
    }
    if perPage < 1 {
        perPage = 10
    }
    
    offset := (page - 1) * perPage
    // totalPages := int((total + int64(perPage) - 1) / int64(perPage))
    
    if err := q.db.WithContext(ctx).Where("deleted_at IS NULL").Offset(offset).Limit(perPage).Find(&items).Error; err != nil {
        return nil, err
    }

    return dto.NewPaginationResponse(items, total, page, perPage), nil
}

func (r *GormRepository[T]) NewSession(ctx context.Context) (*gorm.DB, error) {
    session := r.db.WithContext(ctx).Session(&gorm.Session{})
    return session, nil
}

func (r *GormRepository[T]) Query() *GormRepository[T] {
    return &GormRepository[T]{db: r.db}
}

func (q *GormRepository[T]) WithContext(ctx context.Context) *GormRepository[T] {
    q.db = q.db.WithContext(ctx)
    return q
}

func (q *GormRepository[T]) Join(table string) *GormRepository[T] {
    q.db = q.db.Joins(table)
    return q
}

func (q *GormRepository[T]) LeftJoin(table string) *GormRepository[T] {
    q.db = q.db.Joins("LEFT JOIN " + table)
    return q
}

func (q *GormRepository[T]) InnerJoin(table string) *GormRepository[T] {
    q.db = q.db.Joins("INNER JOIN " + table)
    return q
}

func (q *GormRepository[T]) Where(query interface{}, args ...interface{}) *GormRepository[T] {
    q.db = q.db.Where(query, args...)
    return q
}

func (q *GormRepository[T]) Preload(query string) *GormRepository[T] {
    q.db = q.db.Preload(query)
    return q
}

func (q *GormRepository[T]) Order(value interface{}) *GormRepository[T] {
    q.db = q.db.Order(value)
    return q
}

func (q *GormRepository[T]) Limit(limit int) *GormRepository[T] {
    q.db = q.db.Limit(limit)
    return q
}

func (q *GormRepository[T]) Offset(offset int) *GormRepository[T] {
    q.db = q.db.Offset(offset)
    return q
}