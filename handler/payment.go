package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/FlexHC/MinecraftStore/model"
	"github.com/FlexHC/MinecraftStore/payment/hotpay"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
	"os"
)

type PaymentHandlers struct {
	DB *sqlx.DB
}

type paymentRequest struct {
	ProductId    uint8  `json:"productId"`
	Nickname     string `json:"nickname"`
	PersonalData string `json:"personalData"`
	Email        string `json:"email"`
}

func (p *PaymentHandlers) NewPayment(c *gin.Context) {
	request := paymentRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	var product model.Product
	err = p.DB.Get(&product, "SELECT * FROM `products` WHERE `id`=? LIMIT 1", request.ProductId)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "database error",
		})
		return
	}

	log.Printf("%#v\n", product)
	transaction := hotpay.NewTransaction(product, "http://2.tcp.eu.ngrok.io:15273/", request.Email, request.Nickname, request.PersonalData)
	err = transaction.CreateDatabaseEntry(p.DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error creating payment entry in db",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"url": transaction.GetURL().String(),
	})

}

func (p *PaymentHandlers) PaymentCallback(c *gin.Context) {
	amount := c.PostForm("KWOTA")
	paymentID := c.PostForm("ID_PLATNOSCI")
	orderID := c.PostForm("ID_ZAMOWIENIA")
	status := c.PostForm("STATUS")
	secret := c.PostForm("SEKRET")
	clientHash := c.PostForm("HASH")

	stringToHash := fmt.Sprintf("%s;%s;%s;%s;%s;%s",
		os.Getenv("HOTPAY_HASH"),
		amount,
		paymentID,
		orderID,
		status,
		secret,
	)
	hash := sha256.New()
	hash.Write([]byte(stringToHash))

	if clientHash != hex.EncodeToString(hash.Sum(nil)) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "hash mismatch"})
		return
	}

	_, err := p.DB.Exec("UPDATE `transactions` SET `status` = ? WHERE `id` = ?", status, paymentID)
	if err != nil {
		log.Printf("Error while updating status of transaction %s. New status should be %s. Error: %#v\n", paymentID, status, err)
	}
	if status == hotpay.Done {
		var product model.Product
		err = p.DB.Get(&product, "SELECT `command` FROM `products` WHERE `id` = (SELECT `product_id` FROM `purchases` WHERE `transaction_id` = ?)", orderID)
		if err != nil {
			log.Printf("Error while executing command for transaction %s\n", orderID)
			return
		}
		log.Printf("Executing command %s\n", product.Command)

	}
}
