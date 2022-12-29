package models

type Product struct {
	Product_id    uint   `json:"product_id" gorm:"primaryKey" `
	Product_name  string `json:"product_name" gorm:"not null"  `
	Price         uint   `json:"price" gorm:"not null"  `
	Discount      uint
	Selling_Price uint
	Stock         uint   `json:"stock"  `
	Color         string `json:"color" gorm:"not null"  `
	Description   string `json:"description"   `

	Category string `json:"category"`
	//CategoryID uint
	Brand string `json:"brand"`
	//Brand_id   uint `json:"brand_id" `
	ShoeSize uint `json:"size"`
	//ShoeSizeID uint
}

/*type Category struct {
	ID       uint `json:"id" gorm:"primaryKey"  `
	Category string
}

type Brand struct {
	ID    uint   `json:"id" gorm:"primaryKey"  `
	Brand string `json:"brand"`
	//Discount uint   `json:"discount"`
}

type ShoeSize struct {
	ID   uint `json:"id" gorm:"primaryKey"  `
	Size uint `json:"size"`
}*/
