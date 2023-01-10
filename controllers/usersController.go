package controllers

import (
	"net/http"
	"os"
	"redkart/initializers"
	"redkart/models"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	//get email/pass
	var body struct {
		First_Name string
		Last_Name  string
		Email      string
		Password   string
		Phone      string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//Hash password

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	//create user

	user := models.User{First_Name: body.First_Name, Last_Name: body.Last_Name, Email: body.Email, Password: string(hash), Phone: body.Phone}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	//respond
	c.JSON(http.StatusOK, gin.H{
		"msg": "Signup successfully",
	})

	//user login
}

// UserLogin godoc
//
//	@Summary		API to Login for users
//	@Description	get string by ID
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			admin	body		models.User	true	"User ID"
//	@Success		200		{object}	models.User
//	@Router			/user/login [post]
func Login(c *gin.Context) {
	//get the email& password of req body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})

		return
	}

	//loookup the user
	var user models.User
	initializers.DB.First(&user, "email = ?", body.Email)
	// SELECT * FROM users WHERE email = "email in database";
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}
	if user.Block_status {
		c.JSON(404, gin.H{
			"msg":     "Can't login, user has been blocked By admin",
			"Message": user.ID})
		c.Abort()
		return
	}

	//Compare the passwords
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})

		return
	}

	//generating jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})
		return
	}

	//send it back
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{
		"message":     "User logged in successfully",
		"data":        user,
		"accessToken": tokenString,
	})
}

func Validate(c *gin.Context) {
	var user models.User
	c.SetSameSite(http.SameSiteLaxMode)
	c.JSON(http.StatusOK, gin.H{
		"message": user,
		"msg":     user.Email,
	})
}
