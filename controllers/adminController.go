package controllers

import (
	"net/http"
	"os"
	"redkart/initializers"
	"redkart/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

/*type AdminLogins struct {
	Email    string
	Password string
}

var UserDb = map[string]string{
	"email":    "jeevan@gmail.com",
	"password": "123",
}*/

var adminbody struct {
	Name     string
	Email    string
	Password string
	Phone    string
}

// AdminLogin godoc
//
//	@Summary		API to Login for admins
//	@Description	get string by ID
//	@Tags			admin
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		models.Admin	true	"Account ID"
//	@Success		200		{object}	models.Admin
//	@Router			/admin/adminlogin [post]
func AdminLogin(c *gin.Context) { // admin login page post
	var adminbody struct {
		Email    string
		Password string
	}
	//var u AdminLogins
	var admin models.Admin
	if err := c.ShouldBindJSON(&adminbody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		c.Abort()
		return
	}
	//record := initializers.DB.Raw("select * from admins where email=?", adminbody.Email).Scan(&admin)
	//var record models.Admin
	//initializers.DB.First(&record, "email = ?", adminbody.Email)
	initializers.DB.Raw("select * from admins where email=?", adminbody.Email).Scan(&admin)
	/*if record==0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		c.Abort
		return
	}*/
	if admin.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(adminbody.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	//generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create a token",
		})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})

}

func AdminSignup(c *gin.Context) {
	//get email/pass

	if c.Bind(&adminbody) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//Hash password

	hash, err := bcrypt.GenerateFromPassword([]byte(adminbody.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//create

	admin := models.Admin{Name: adminbody.Name, Email: adminbody.Email, Password: string(hash), Phone: adminbody.Phone}
	result := initializers.DB.Create(&admin)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create admin",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{})

	//user login
}

// @Summary block user by ID
// @ID block-user-by-id
// @Tags admin
// @Accept json
// @Produce json
// @Param id path string true "user ID"
// @Success 200 {object} models.User
// @Router /admin//userdata/block/{id} [put]
func BlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	initializers.DB.Raw("UPDATE users SET block_status=true where id=?", params).Scan(&user)
	//initializers.DB.Raw("UPDATE users SET block_status=true where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Blocked succesfully"})
}
func UnBlockUser(c *gin.Context) {
	params := c.Param("id")
	var user models.User
	initializers.DB.Raw("UPDATE users SET block_status=false where id=?", params).Scan(&user)
	c.JSON(http.StatusOK, gin.H{"msg": "Unblocked succesfully"})
}

func ListUsers(c *gin.Context) {
	// var user []models.User
	// initializers.DB.Find((&user))
	// c.JSON(http.StatusOK, user)
	pagestring := c.Query("page")
	page, _ := strconv.Atoi(pagestring)
	offset := (page - 1) * 3
	var user []models.User
	initializers.DB.Limit(3).Offset(offset).Find(&user)
	c.JSON(http.StatusOK, user)
}

func ValidateAdmin(c *gin.Context) {
	var admin models.Admin
	c.JSON(http.StatusOK, gin.H{
		"message": admin,
	})
}
