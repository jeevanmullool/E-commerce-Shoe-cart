package models

import "gorm.io/gorm"

type Coupon struct {
	gorm.Model
	Coupon_code string `json:"coupon_code"`
	Discount    uint   `json:"discount"`
	Min_value   uint   `json:"min_value"`
	Exp_date    int64
	Status      bool
}
