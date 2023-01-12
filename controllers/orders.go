package controllers

import (
	"fmt"
	"net/http"
	"redkart/initializers"
	"redkart/models"

	"github.com/gin-gonic/gin"
)

//view all orders from admin side

func ViewOrders(c *gin.Context) {
	var orders []models.Orders
	initializers.DB.Raw("select * from orders").Scan(&orders)
	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

func Cancelorders(c *gin.Context) {
	var user models.User
	userEmail := c.GetString("user")
	orderid := c.Query("orderID")
	update_status := "order cancelled"
	initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user) //getting user id
	var orders models.Orderd_Items
	// initializers.DB.First(&orders).Where("orders_id=?", orderid)
	initializers.DB.Where("orders_id=?", orderid).Find(&orders)
	fmt.Println(orders.Order_Status, orders.OrdersID)
	fmt.Println(orderid)
	if orders.Order_Status == update_status {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Order already Cancelled",
		})
		return
	}
	initializers.DB.Raw("update orderd_items set order_status=? where orders_id=?", "order cancelled", orderid).Scan(&orders)
	c.JSON(200, gin.H{
		"msg": "order canccelled",
	})

}
