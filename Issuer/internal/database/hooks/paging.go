package hooks

import "gorm.io/gorm"

func Paging(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		pageIndex := page
		if pageIndex <= 0 {
			pageIndex = 1
		}

		pageSize := size
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageIndex - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
