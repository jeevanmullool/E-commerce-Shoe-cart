package controllers

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"redkart/initializers"
	"redkart/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//var UserEmail string

func AddToCart(c *gin.Context) {
	useremail := c.GetString("user")
	fmt.Println("Hellouser")
	fmt.Println(useremail)
	var UsersID int
	var products models.Product
	//initializers.DB.Raw("select id from users where email=?", userEmail).Scan(&user)
	err := initializers.DB.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	//err := initializers.DB.First(&UsersID, "email = ?", user)

	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var ProdtDetails struct {
		Product_id uint
		Quantity   uint
	}

	c.BindJSON(&ProdtDetails)
	fmt.Println(ProdtDetails.Product_id)
	fmt.Println(ProdtDetails.Quantity)
	var stock uint
	initializers.DB.Raw("SELECT stock FROM products WHERE product_id = ?", ProdtDetails.Product_id).Scan(&stock)

	if stock < ProdtDetails.Quantity {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Product is Out Of Stock",
		})
		return
	}
	//geting price for setting totalamount
	initializers.DB.Raw("select selling_price ,stock from products where product_id=?", ProdtDetails.Product_id).Scan(&products)
	total := products.Selling_Price * ProdtDetails.Quantity
	prodid := ProdtDetails.Product_id
	prodqua := ProdtDetails.Quantity
	var name string
	var brand string
	initializers.DB.Raw("SELECT product_name FROM products WHERE product_id = ?", prodid).Scan(&name)
	initializers.DB.Raw("SELECT brand FROM products WHERE product_id = ?", prodid).Scan(&brand)
	cart := models.Cart{
		ProductID:    ProdtDetails.Product_id,
		Quantity:     ProdtDetails.Quantity,
		Product_Name: name,
		Brand_name:   brand,
		UserId:       uint(UsersID),
		Total:        total,
	}
	var Cart []models.Cart
	initializers.DB.Raw("select cart_id,product_id from carts where user_id=?", UsersID).Scan(&Cart) //geting all the cart details associated to user
	//ranging in the cart to find if the product already exists
	for _, l := range Cart {
		fmt.Println("enterd")
		if l.ProductID == prodid {
			fmt.Println("in")
			initializers.DB.Raw("select quantity from carts where product_id=? and user_id=?", ProdtDetails.Product_id, UsersID).Scan(&Cart)
			totl := (prodqua + cart.Quantity) * products.Price
			totqua := prodqua + cart.Quantity
			initializers.DB.Raw("update carts set quantity=?,total=? where product_id=? and user_id=? ", totqua, totl, prodid, UsersID).Scan(&Cart)

			c.JSON(400, gin.H{
				"msg":  "quantity updated",
				"user": UsersID,
			})
			c.Abort()
			return
		}
	}

	record := initializers.DB.Create(&cart)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"err": record.Error.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{
		"msg":      "added to cart",
		"products": cart,
	})
}

func ViewCart(c *gin.Context) {
	var Subtotal int
	useremail := c.GetString("user")
	var UsersID int
	err := initializers.DB.Raw("select id from users where email=?", useremail).Scan(&UsersID)
	fmt.Println("In viewcart")
	fmt.Println(UsersID)
	if errors.Is(err.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, gin.H{
			"message": "user coudnt find",
		})

	}
	var cart []struct {
		ID           int
		Product_id   int
		Product_Name string
		Brand_Name   string
		Total        int
		Quantity     int
	}
	initializers.DB.Select("user_id", "product_id", "product_name", "brand_name", "total", "quantity").Table("carts").Where("user_id=?", UsersID).Find(&cart)
	c.JSON(200, gin.H{
		"Products": cart,
	})

	for _, i := range cart {
		sum := i.Total
		Subtotal = Subtotal + sum
	}
	c.JSON(200, gin.H{
		"status":   true,
		"Subtotal": Subtotal,
	})

}

type Cartsinfo []struct {
	User_id       string
	Product_id    string
	Product_Name  string
	Price         string
	Selling_Price string
	Email         string
	Quantity      string
	Total_Amount  uint
	Total_Price   string
}

func AddAddressCheckout(c *gin.Context) {
	useremail := c.GetString("user")
	var user models.User
	initializers.DB.Raw("select id from users where email=?", useremail).Scan(&user)

	Name := c.PostForm("name")
	Phonenum := c.PostForm("phone_number")
	phonenum, _ := strconv.Atoi(Phonenum)
	pincod := c.PostForm("pincode")
	pincode, _ := strconv.Atoi(pincod)
	area := c.PostForm("area")
	houseadd := c.PostForm("house")
	landmark := c.PostForm("landmark")
	city := c.PostForm("city")
	address := models.Address{
		UserId:       user.ID,
		Name:         Name,
		Phone_number: phonenum,
		Pincode:      pincode,
		Area:         area,
		House:        houseadd,
		Landmark:     landmark,
		City:         city,
	}
	record := initializers.DB.Create(&address)
	if record.Error != nil {
		c.JSON(404, gin.H{"error": record.Error.Error()})
		c.Abort()
		return
	}
	c.JSON(200, gin.H{"msg": "address added"})

}

var Address []struct {
	UserId       uint
	Address_id   uint
	Name         string
	Phone_number uint
	Pincode      uint
	Area         string
	House        string
	Landmark     string
	City         string
}

func PlaceOrder(c *gin.Context) {
	var user models.User
	//var products models.Product
	var cart models.Cart
	var cartinf Cartsinfo
	useremail := c.GetString("user")

	initializers.DB.Raw("select id from users where email=?", useremail).Scan(&user)
	record := initializers.DB.Raw("select  products.product_id, products.product_name,products.selling_price,carts.user_id,users.email ,carts.quantity,total from carts join products on products.product_id=carts.product_id join users on carts.user_id=users.id where users.email=? ", useremail).Scan(&cartinf)
	if record.Error != nil {
		c.JSON(404, gin.H{
			"error": record.Error.Error(),
		})
		c.Abort()
		return
	}

	var totalcartvalue uint
	var address models.Address
	addres := c.Query("addressID")
	addressID, _ := strconv.Atoi(addres)

	initializers.DB.Raw("select sum(total) as total from carts where user_id=?", user.ID).Scan(&totalcartvalue)

	payMeth := c.Query("payment_method")
	orderidd := CreateOrderId()
	if payMeth == "COD" {
		fmt.Println(" In cash on delivery")

		for _, l := range cartinf {
			ui := l.User_id
			jui, _ := strconv.Atoi(ui)
			pid := l.Product_id
			pidd, _ := strconv.Atoi(pid)
			pname := l.Product_Name
			pprice := l.Selling_Price
			pPrice, _ := strconv.Atoi(pprice)
			pquantity := l.Quantity
			pQuantity, _ := strconv.Atoi(pquantity)
			totamount := uint(pQuantity) * uint(pPrice)
			ordereditems := models.Orderd_Items{UserId: uint(jui), Product_id: uint(pidd),
				Product_Name: pname, Price: pprice, OrdersID: orderidd,
				Order_Status: "confirmed", Payment_Status: "COD", Total_amount: totamount,
			}
			initializers.DB.Create(&ordereditems)
		}
		//getting details from address
		precord := initializers.DB.Raw("select address_id, user_id,name,phone_number,pincode,house,area,landmark,city from addresses where user_id=?", user.ID).Scan(&Address)
		if precord.Error != nil {
			c.JSON(404, gin.H{
				"err": precord.Error.Error(),
			})
			c.Abort()
			return
		}

		initializers.DB.Raw("select address_id,user_id,name from addresses where address_id=?", addressID).Scan(&address)

		c.JSON(300, gin.H{
			"address":          Address,
			"total cart value": totalcartvalue,
		})

		//if addressID == int(address.Address_id) && address.UserId == user.ID {
		orders := models.Orders{
			UserId:         user.ID,
			Addresss_id:    uint(addressID),
			Total_Amount:   totalcartvalue,
			Order_id:       orderidd,
			Order_Status:   "order Placed",
			Payment_Status: "COD",
		}
		result := initializers.DB.Create(&orders)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"err": result.Error.Error(),
			})
			c.Abort()
			return

		}
	} else if payMeth == "STRIPE" {
		for _, l := range cartinf {
			ui := l.User_id
			jui, _ := strconv.Atoi(ui)
			pid := l.Product_id
			pidd, _ := strconv.Atoi(pid)
			pname := l.Product_Name
			pprice := l.Selling_Price
			pPrice, _ := strconv.Atoi(pprice)
			pquantity := l.Quantity
			pQuantity, _ := strconv.Atoi(pquantity)
			totamount := uint(pQuantity) * uint(pPrice)
			ordereditems := models.Orderd_Items{UserId: uint(jui), Product_id: uint(pidd),
				Product_Name: pname, OrdersID: CreateOrderId(),
				Order_Status: "confirmed", Payment_Status: "Stripe Payment", Total_amount: totamount,
			}
			initializers.DB.Create(&ordereditems)
		}
		//getting details from address
		precord := initializers.DB.Raw("select address_id, user_id,name,phone_number,pincode,house,area,landmark,city from addresses where user_id=?", user.ID).Scan(&Address)
		if precord.Error != nil {
			c.JSON(404, gin.H{
				"err": precord.Error.Error(),
			})
			c.Abort()
			return
		}

		initializers.DB.Raw("select address_id,user_id,name from addresses where address_id=?", addressID).Scan(&address)

		c.JSON(300, gin.H{
			"address":          Address,
			"total cart value": totalcartvalue,
		})

		//if addressID == int(address.Address_id) && address.UserId == user.ID {
		orders := models.Orders{
			UserId:         user.ID,
			Addresss_id:    uint(addressID),
			Total_Amount:   totalcartvalue,
			Order_id:       orderidd,
			Order_Status:   "order Placed",
			Payment_Status: "Stripe Payment",
		}
		result := initializers.DB.Create(&orders)
		if result.Error != nil {
			c.JSON(404, gin.H{
				"err": result.Error.Error(),
			})
			c.Abort()
			return

		}
		StripePayment(c)

	}
	var proid uint
	initializers.DB.Raw("select product_id from carts where user_id=?", user.ID).Scan(&proid)
	var stock uint
	initializers.DB.Raw("select products.stock from products where product_id=?", proid).Scan(&stock)
	var cartQuantity uint
	initializers.DB.Raw("select quantity from carts where user_id=?", user.ID).Scan(&cartQuantity)
	newQuantity := stock - cartQuantity
	fmt.Println(newQuantity)
	fmt.Println(cartQuantity)
	fmt.Println(proid)
	initializers.DB.Raw("delete from carts where user_id=?", user.ID).Scan(&cart)
	var Prod []models.Product
	initializers.DB.Raw("update products set stock=? where product_id=?", newQuantity, proid).Scan(&Prod)

	//c.JSON(200, gin.H{"order": "order placed"})
}

func CreateOrderId() string {
	rand.Seed(time.Now().UnixNano())
	value := rand.Intn(9999999999-1000000000) + 1000000000
	id := strconv.Itoa(value)
	orderID := "RED" + id
	return orderID
}
