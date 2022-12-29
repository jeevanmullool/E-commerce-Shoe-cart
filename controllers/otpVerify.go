package controllers

import (
	"fmt"
	"net/http"
	"os"
	"redkart/initializers"
	"redkart/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/twilio/twilio-go"
	verify "github.com/twilio/twilio-go/rest/verify/v2"
)

var (
	accountSid string
	serviceSid string
	authToken  string
	fromPhone  string

	client *twilio.RestClient
)

func SendOtp(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	serviceSid = os.Getenv("VERIFY_SERVICE_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("phone")

	result := ChekNumber(Mob)
	fmt.Println(result)

	if !result {
		c.JSON(400, gin.H{
			"status":  false,
			"message": "Mobile number doesnt exist! Please SignUp",
		})
		return
	}

	mobile := "+91" + Mob

	params := &verify.CreateVerificationParams{}
	params.SetTo(mobile)
	params.SetChannel("sms")
	fmt.Println(accountSid)
	fmt.Println(mobile)
	fmt.Println(params)
	resp, err := client.VerifyV2.CreateVerification(serviceSid, params)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(400, gin.H{
			"status":  false,
			"message": "error sending OTP",
		})
	} else {
		fmt.Printf("Sent verification '%s'\n", *resp.Sid)
		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "OTP Sent Succesfully",
		})
	}

}

// Checking number already used

func ChekNumber(str string) bool {

	mobilenumber := str
	var checkOtp models.User
	initializers.DB.Raw("SELECT phone FROM users WHERE phone=?", mobilenumber).Scan(&checkOtp)
	return checkOtp.Phone == mobilenumber

}

func CheckOtp(c *gin.Context) {
	accountSid = os.Getenv("ACCOUNT_SID")
	serviceSid = os.Getenv("VERIFY_SERVICE_SID")
	authToken = os.Getenv("AUTH_TOKEN")
	fromPhone = os.Getenv("FROM_PHONE")
	client = twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	Mob := c.Query("number")
	code := c.Query("otps")

	ChekNumber(Mob)
	var user models.User
	initializers.DB.First(&user, "phone = ?", Mob)

	mobile := "+91" + Mob
	fromPhone = os.Getenv("FROM_PHONE")
	fmt.Println(mobile)
	params := &verify.CreateVerificationCheckParams{}
	params.SetTo(mobile)
	params.SetCode(code)

	resp, err := client.VerifyV2.CreateVerificationCheck(fromPhone, params)

	if err != nil {
		fmt.Println(err.Error())
	} else if *resp.Status == "approved" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.Email,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})
		// Sent it back
		/*fmt.Println(tokenstring)
		token := tokenstring["access_token"]*/
		tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to create a token",
			})
			return
		}
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"status":  true,
			"message": "ok",
			"data":    tokenString,
		})
	} else {

		c.JSON(404, gin.H{
			"msg": "otp is invalid",
		})
	}
}
