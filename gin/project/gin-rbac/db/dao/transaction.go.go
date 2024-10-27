package dao

import "gorm.io/gorm"

// WithTransaction 封装事务处理
func Transaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
