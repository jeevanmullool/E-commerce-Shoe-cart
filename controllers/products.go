package controllers

import (
	"fmt"
	"net/http"
	"redkart/initializers"
	"redkart/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary get all items in the product list
// @ID get-all-products
// @Tags User
// @Produce json
// @Success 200 {object} models.Product
// @Router /user/products [get]
// product listing with pagination , we can input page number to list products in different pages
func ProductList(c *gin.Context) {
	pagestring := c.Query("page")
	page, _ := strconv.Atoi(pagestring)
	offset := (page - 1) * 3
	var product []models.Product
	initializers.DB.Limit(3).Offset(offset).Find(&product)
	c.JSON(http.StatusOK, product)
}

// func ProductList(c *gin.Context) {
// 	var product []models.Product
// 	initializers.DB.Find((&product))
// 	c.JSON(http.StatusOK, product)
// }

var Products []struct {
	Product_ID    uint
	Product_Name  string
	Price         string
	Selling_Price string
	Description   string
	Color         string
	Brands        string
	Stock         uint
	Category      string
	Size          uint
}

func ProductAdding(c *gin.Context) { //Admin

	prodname := c.PostForm("productname")
	price := c.PostForm("price")
	Price, _ := strconv.Atoi(price)
	description := c.PostForm("description")
	color := c.PostForm("color")
	brand := c.PostForm("brandID")
	//brands, _ := strconv.Atoi(brand)
	stock := c.PostForm("stock")
	Stock, _ := strconv.Atoi(stock)
	catogory := c.PostForm("categoryID")
	//catogoryy, _ := strconv.Atoi(catogory)
	size := c.PostForm("sizeID")

	Size, _ := strconv.Atoi(size)
	discont := c.PostForm("discount")
	discount, _ := strconv.Atoi(discont)

	//calculating discounted price
	disprice := Price - (Price * discount / 100)

	var count uint
	initializers.DB.Raw("select count(*) from products where product_name=?", prodname).Scan(&count)
	fmt.Println(count)
	if count > 0 {
		c.JSON(404, gin.H{
			"msg": "A product with same name already exists",
		})
		c.Abort()
		return
	}
	products := models.Product{

		Product_name: prodname,

		Price:         uint(Price),
		Color:         color,
		Description:   description,
		Brand:         brand,
		Category:      catogory,
		ShoeSize:      uint(Size),
		Stock:         uint(Stock),
		Discount:      uint(discount),
		Selling_Price: uint(disprice),
	}

	record := initializers.DB.Create(&products)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"msg": "product already exists",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"msg": "added succesfully",
	})

}

//delete product

// @Summary delete a product by ID
// @ID delete-product-by-id
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "product ID"
// @Success 200 {object} models.Product
// @Faillure 400 {object} message
// @Router /admin/deleteproduct [delete]
func DeleteProductById(c *gin.Context) { //admin
	params := c.Query("id")
	var products models.Product
	var count uint
	initializers.DB.Raw("select count(product_id) from products where product_id=?", params).Scan(&count)
	if count == 0 {
		c.JSON(400, gin.H{
			"msg": "product doesnot exist",
		})
		c.Abort()
		return
	}

	record := initializers.DB.Raw("delete from products where product_id=?", params).Scan(&products)
	if record.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"msg": "deleted successfully"})
}
