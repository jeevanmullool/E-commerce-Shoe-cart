package initializers

import (
	"redkart/models"
)

func SyncDatabase() {
	DB.AutoMigrate((&models.User{}),
		(&models.Admin{}),
		(&models.Product{}),
		(&models.Otp{}),
		(&models.Cart{}),
		(&models.Address{}),
		(&models.Orderd_Items{}),
		(&models.Orders{}),
		(&models.WishList{}),
		(&models.Charge{}),
		(&models.Coupon{}),
	)
}
