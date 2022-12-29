package controllers

import (
	"fmt"
	"net/http"
	"os"
	"redkart/initializers"
	"redkart/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/charge"
)

func StripePayment(c *gin.Context) {
	var user models.User
	var payment models.Charge
	c.BindJSON(&payment)
	var totalcartvalue int64
	useremail := c.GetString("user")
	initializers.DB.Raw("select id from users where email=?", useremail).Scan(&user)
	initializers.DB.Raw("select sum(total) as total from carts where user_id=?", user.ID).Scan(&totalcartvalue)
	apiKey := os.Getenv("STRIPE_KEY")
	fmt.Println(apiKey + "asads")
	stripe.Key = apiKey
	var body []string
	initializers.DB.Raw("select product_name from carts where user_id=?", user.ID).Scan(&body)
	fmt.Println(user.ID)
	fmt.Println(body)
	a1 := strings.Join(body, ",")
	fmt.Println(a1)

	_, err := charge.New(&stripe.ChargeParams{
		Amount:       &totalcartvalue,
		Currency:     stripe.String(string(stripe.CurrencyUSD)),
		Description:  &a1,
		Source:       &stripe.SourceParams{Token: stripe.String("tok_visa")},
		ReceiptEmail: &useremail,
	})

	if err != nil {
		c.String(http.StatusBadRequest, "Payment Unsuccessfull")
		return
	}

	err1 := SavePayment(&payment)
	if err1 != nil {
		c.String(http.StatusBadRequest, "error occured")
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "Payment Successful",
		})
	}

}

func SavePayment(charge *models.Charge) (err error) {
	if err = initializers.DB.Create(charge).Error; err != nil {
		return err
	}
	return nil

}
