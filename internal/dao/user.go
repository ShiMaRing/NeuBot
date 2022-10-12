package dao

import "gorm.io/gorm"

// Dao dao,暂时实现假的业务逻辑
type Dao struct {
	gorm.DB
}
