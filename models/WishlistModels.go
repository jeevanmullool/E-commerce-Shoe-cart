package models

import "gorm.io/gorm"

type WishList struct {
	gorm.Model
	UserID     uint
	Product_id uint
}
