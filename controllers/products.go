package controllers

import (
	"fmt"
	"net/http"
	"redkart/initializers"
	"redkart/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*var productbody struct {
	Product_name string
	brand        string
	Price        uint
	Stock        uint
	Color        string
	Description  string
	Category     string
	Size         uint
}*/

/*func AddProduct(c *gin.Context) {

	if c.Bind(&productbody) == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	products := models.Product{Product_name: productbody.Product_name, Brand: productbody.Brand, Price: productbody.Price, Stock: productbody.Stock, Color: productbody.Color, Description: productbody.Description, Category: productbody.Category, Size: productbody.Size}
	result := initializers.DB.Create(&products)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to add product",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{})
}*/

func ProductList(c *gin.Context) {
	var product []models.Product
	initializers.DB.Find((&product))
	c.JSON(http.StatusOK, product)
}

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

func DeleteProductById(c *gin.Context) { //admin
	params := c.Param("id")
	var products models.Product
	var count uint
	initializers.DB.Raw("select count(product_id) from products where product_id=?", params).Scan(&count)
	if count <= 0 {
		c.JSON(404, gin.H{
			"msg": "product doesnot exist",
		})
		c.Abort()
		return
	}

	record := initializers.DB.Raw("delete from products where product_id=?", params).Scan(&products)
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"msg": "deleted successfully"})
}
