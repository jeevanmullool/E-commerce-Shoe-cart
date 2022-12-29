package models

type User struct {
	ID           uint   `json:"id" gorm:"primaryKey;unique"  `
	First_Name   string `json:"first_name"  gorm:"not null" validate:"required,min=2,max=50"  `
	Last_Name    string `json:"last_name"    gorm:"not null"    validate:"required,min=1,max=50"  `
	Email        string `json:"email"   gorm:"not null;unique"  validate:"email,required"`
	Password     string `json:"password" gorm:"not null"  validate:"required"`
	Phone        string `json:"phone"   gorm:"not null;unique" validate:"required"`
	Block_status bool   `json:"block_status" gorm:"not null"   `
	Country      string `json:"country "   `
	City         string `json:"city "   `
	Pincode      uint   `json:"pincode "   `

	Address    Address
	Address_id uint `json:"address_id" `
	//Cart_id      uint   `json:"cart_id" `
	//Address_id   uint   `json:"address_id" `
	//Orders_ID    uint   `json:"orders_id" `
}

type Address struct {
	Address_id   uint   `json:"address_id" gorm:"primaryKey"  `
	UserId       uint   `json:"user_id"  gorm:"not null" `
	Name         string `json:"name"  gorm:"not null" `
	Phone_number int    `json:"phone_number"  gorm:"not null" `
	Pincode      int    `json:"pincode"  gorm:"not null" `
	House        string `json:"house"   `
	Area         string `json:"area"   `
	Landmark     string `json:"landmark"  gorm:"not null" `
	City         string `json:"city"  gorm:"not null" `
}
