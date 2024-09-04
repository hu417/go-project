package curd

import (
	"bluebell/models/req"

	"gorm.io/gorm"
)

func Paginate(p *req.ParamPage) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
	   page := p.Page
	   if page == 0 {
		  page = 1
	   }
	   pageSize := p.PageSize
	   switch {
	   case pageSize > 100:
		  pageSize = 100
	   case pageSize <= 0:
		  pageSize = 10
	   }
	   offset := (page - 1) * pageSize
	   return db.Offset(offset).Limit(pageSize)
	}
 }
 