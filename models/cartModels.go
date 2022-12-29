package models

import "gorm.io/gorm"

type Cart struct {
	Cart_id      uint   `json:"cart_id" gorm:"primaryKey"  `
	UserId       uint   `json:"user_id"`
	ProductID    uint   `json:"product_id"`
	Product_Name string `json:"product_name"`
	Brand_name   string `json:"brand_name"`
	Quantity     uint   `json:"quantity"`
	Total        uint   `json:"total"`
}

type Cartdetails struct {
	gorm.Model
	user_id     string
	Product_id  string
	ProductName string
	Price       string
	Email       string
	Quantity    string
	Total_Price string
}
type Otp struct {
	gorm.Model
	Mobile string
	Otp    string
}

type PaymentMethod struct {
	COD bool
}
type Orders struct {
	gorm.Model
	UserId       uint   `json:"user_id"  gorm:"not null" `
	Order_id     string `json:"order_id"  gorm:"not null" `
	Total_Amount uint   `json:"total_amount"  gorm:"not null" `
	//PaymentMethod   string `json:"paymentmethod"  gorm:"not null" `
	Payment_Status string `json:"payment_status"   `
	Order_Status   string `json:"order_status"   `
	Addresss_id    uint
}

type Orderd_Items struct {
	gorm.Model
	UserId         uint `json:"user_id"  gorm:"not null" `
	Product_id     uint
	OrdersID       string
	Product_Name   string
	Price          string
	Order_Status   string
	Payment_Status string
	Total_amount   uint
}
