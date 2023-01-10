package controllers

import (
	"redkart/initializers"
	"redkart/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func AddCoupon(c *gin.Context) {
	co_code := c.PostForm("coupon_code")
	//coup_code, _ := strconv.Atoi(co_code)
	co_disc := c.PostForm("coupon_discount")
	coup_disc, _ := strconv.Atoi(co_disc)
	min_val := c.PostForm("minimum_value")
	min_value, _ := strconv.Atoi(min_val)
	//exp_date := c.PostForm("expiry_date")
	//expiry_date, _ := strconv.Atoi(exp_date)

	coupons := models.Coupon{
		Coupon_code: co_code,
		Discount:    uint(coup_disc),
		Min_value:   uint(min_value),
		Exp_date:    time.Now().Add(time.Hour * 24 * 30).Unix(),
		Status:      false,
	}

	record := initializers.DB.Create(&coupons)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"msg": "Coupon adding failed",
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"msg":  "Coupon added successfully",
		"data": coupons,
	})
}

func RedeemCoupon(c *gin.Context) {
	var user models.User
	useremail := c.GetString("user")
	initializers.DB.Raw("select id from users where email=?", useremail).Scan(&user)

	var cartt models.Cart
	initializers.DB.Raw("select * from cart where user_id=?", user.ID).Scan(&cartt)

	coup_code := c.Query("coupon_code")
	var coup models.Coupon
	initializers.DB.Raw("select * from coupons where coupon_code=?", coup_code).Scan(&coup)
	if coup.ID == 0 {
		c.JSON(400, gin.H{"Msg": "Failed to fetch coupon"})
	}
	var grandtotal uint
	initializers.DB.Raw("select sum(total) from carts where user_id=?", user.ID).Scan(&grandtotal)

	if !coup.Status && coup.Exp_date > time.Now().Unix() && coup.Min_value < grandtotal {
		redeemed := grandtotal - coup.Discount
		initializers.DB.Raw("update carts set total=? where user_id=?", redeemed, user.ID).Scan(&cartt)
		initializers.DB.Raw("update coupons set status=? where coupon_code=?", true, coup_code).Scan(&coup)
		c.JSON(200, gin.H{
			"msg":   "Coupon claimed successfully",
			"Total": redeemed,
		})
	} else {
		c.JSON(400, gin.H{
			"msg": "Coupon Invalid",
		})
	}

}

func EditOffer(c *gin.Context) {
	var user models.User
	UserEmail := c.GetString("user")
	initializers.DB.Where("email=?", UserEmail).Find(&user)

	var Offer struct {
		Product_id uint
		Discount   uint
	}
	if err := c.BindJSON(&Offer); err != nil {
		c.JSON(404, gin.H{
			"err": err.Error(),
		})
	}

	//var products models.Product
	var price uint
	initializers.DB.Raw("select price from products where product_id=?", Offer.Product_id).Scan(&price)
	disprice := price - (price * Offer.Discount / 100)
	var Prod []models.Product
	initializers.DB.Raw("update products set selling_price=? where product_id=?", disprice, Offer.Product_id).Scan(&Prod)
	initializers.DB.Raw("update products set discount=? where product_id=?", Offer.Discount, Offer.Product_id).Scan(&Prod)
	c.JSON(200, gin.H{
		"msg": "Offer edited successfully",
	})
}
