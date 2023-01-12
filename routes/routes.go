package routes

import (
	"context"
	"redkart/controllers"
	"redkart/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(c *gin.Engine) {

	user := c.Group("/user")
	{
		user.POST("/signup", controllers.Signup)
		user.POST("/login", controllers.Login)
		user.GET("/validateuser", middleware.UserAuth, controllers.Validate)
		user.GET("/products", controllers.ProductList)
		user.POST("/login/otpverify", controllers.SendOtp)
		user.POST("/login/otpcheck", controllers.CheckOtp)
		user.POST("/addtocart", middleware.UserAuth, controllers.AddToCart)
		user.GET("/viewcart", middleware.UserAuth, controllers.ViewCart)
		user.POST("/redeemcoupon", middleware.UserAuth, controllers.RedeemCoupon)
		user.POST("/addaddresscheckout", middleware.UserAuth, controllers.AddAddressCheckout)
		user.POST("/placeorder", middleware.UserAuth, controllers.PlaceOrder)
		user.POST("/addtowishlist", middleware.UserAuth, controllers.AddtoWishlist)
		user.GET("/wishlistview", middleware.UserAuth, controllers.WishlistView)
		user.DELETE("/removefromwishlist", middleware.UserAuth, controllers.RemoveFromWishlist)
		user.GET("/wishlisttocart", middleware.UserAuth, controllers.WishlistToCart)
		user.POST("/stripepayment", middleware.UserAuth, controllers.StripePayment)
		user.POST("/cancelorder", middleware.UserAuth, controllers.Cancelorders)
		user.POST("/paypal", middleware.UserAuth, controllers.Orders(context.Background()))

	}
}

func AdminRoutes(c *gin.Engine) {
	admin := c.Group("/admin")
	{
		admin.POST("/adminsignup", controllers.AdminSignup)
		admin.POST("/adminlogin", controllers.AdminLogin)
		admin.GET("/validateadmin", middleware.AdminAuth, controllers.ValidateAdmin)
		admin.PUT("/userdata/block/:id", middleware.AdminAuth, controllers.BlockUser)
		admin.PUT("/userdata/unblock/:id", middleware.AdminAuth, controllers.UnBlockUser)
		admin.GET("/listuser", middleware.AdminAuth, controllers.ListUsers)
		//admin.POST("/addproduct", middleware.AdminAuth, controllers.AddProduct)
		admin.POST("/addproduct", middleware.AdminAuth, controllers.ProductAdding)
		admin.DELETE("/deleteproduct", middleware.AdminAuth, controllers.DeleteProductById)
		admin.GET("/vieworders", middleware.AdminAuth, controllers.ViewOrders)
		admin.POST("/addcoupon", middleware.AdminAuth, controllers.AddCoupon)
		admin.PATCH("/editoffer", middleware.AdminAuth, controllers.EditOffer)
	}
}
