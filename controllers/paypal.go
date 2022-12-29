package controllers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	paypal "github.com/logpacker/PayPal-Go-SDK"
	paypalsdk "github.com/netlify/PayPal-Go-SDK"
)

// func Paypal(C *gin.Context) {
// 	clientID := os.Getenv("clientID")
// 	secretID := os.Getenv("Paypal_Secret")
// 	// Initialize client
// 	c, err := paypal.NewClient(clientID, secretID, paypal.APIBaseSandBox)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Retrieve access token
// 	_, err = c.GetAccessToken(context.Background())
// 	if err != nil {
// 		panic(err)
// 	}

// 	payout := paypal.Payout{
// 		SenderBatchHeader: &paypal.SenderBatchHeader{
// 			EmailSubject: "Subject will be displayed on PayPal",
// 		},
// 		Items: []paypal.PayoutItem{
// 			paypal.PayoutItem{
// 				RecipientType: "EMAIL",
// 				Receiver:      "single-email-payout@mail.com",
// 				Amount: &paypal.AmountPayout{
// 					Value:    "15.11",
// 					Currency: "USD",
// 				},
// 				Note:         "Optional note",
// 				SenderItemID: "Optional Item ID",
// 			},
// 		},
// 	}

// 	payoutResp, err := c.CreateSinglePayout(payout)
// }

func Paypal(C *gin.Context) {
	clientID := os.Getenv("clientID")
	secretID := os.Getenv("Paypal_Secret")
	c, err := paypalsdk.NewClient(clientID, secretID, paypalsdk.APIBaseSandBox)
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed paypal attempt 1 ",
		})
		return
	}
	fmt.Println(secretID)
	fmt.Println(clientID)
	// accessToken, err :=
	c.GetAccessToken()
	// fmt.Println(accessToken)

	amount := paypalsdk.Amount{
		Total:    "7.00",
		Currency: "USD",
	}
	redirectURI := "https://www.youtube.com/"
	cancelURI := "https://www.youtube.com/"
	description := "Description for this payment"
	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)
	if err != nil {
		C.JSON(400, gin.H{
			"error": "failed paypal attempt 1 ",
		})
		return
	}
	fmt.Println(paymentResult)
}

func pay(ctx context.Context, unit int, token, value string) error {
	appContext := &paypal.ApplicationContext{
		ReturnURL: "http://localhost:8080/success",
		CancelURL: "http://localhost:8080/cancel",
	}
	request := []paypal.PurchaseUnitRequest{
		{
			Amount: &paypal.PurchaseUnitAmount{
				Value:    value,
				Currency: "USD",
				Breakdown: &paypal.PurchaseUnitAmountBreakdown{
					ItemTotal: &paypal.Money{
						Value:    value,
						Currency: "USD",
					},
				},
			},
		},
	}
	orderPayer := &paypal.CreateOrderPayer{}
	clientID := os.Getenv("clientID")
	secretID := os.Getenv("Paypal_Secret")
	c, err := paypal.NewClient(clientID, secretID, paypal.APIBaseSandBox)
	if err != nil {
		fmt.Println("paypal.NewClient", err)
	}
	c.SetLog(os.Stdout) // Set log to terminal stdout

	_, err = c.GetAccessToken()
	if err != nil {
		fmt.Println("c.GetAccessToken", err)
	}
	c.SetLog(os.Stdout) // Set log to terminal stdout

	order, err := c.CreateOrder(paypal.OrderIntentCapture, request, orderPayer, appContext)
	if err != nil {
		return err
	}
	for {
		time.Sleep(20 * time.Second)
		orderStatus, err := c.GetOrder(order.ID)
		if err != nil {
			return err
		}
		if orderStatus.Status == "completed" {
			break
		}
	}
	capRequest := paypal.CaptureOrderRequest{}
	capOrder, err := c.CaptureOrder(order.ID, capRequest)
	if err != nil {
		return err
	}
	fmt.Println(capOrder)
	return nil
}

func Orders(ctx context.Context) gin.HandlerFunc {
	fnc := func(c *gin.Context) {
		//reqBody := make(map[string]string)
		//c.Bind(&reqBody)
		if err := pay(ctx, 1, "foo", strconv.Itoa(1000)); err != nil {
			fmt.Println("amb1s1 error: ", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		} else {

			c.JSON(http.StatusOK, gin.H{})
		}

		fmt.Println("amb1s1 id: ", 1000)
	}
	return fnc
}
