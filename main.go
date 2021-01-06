package main

import (
	"github.com/FlexHC/MinecraftStore/handler"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"log"
)


func main() {
	db, err := sqlx.Connect("mysql", "root:passwd@tcp(localhost:3306)/minecraft_store?parseTime=true")
	if err != nil {
		log.Fatal(err)
		return
	}

	productHandlers := handler.ProductHandlers{DB: db}
	paymentHandlers := handler.PaymentHandlers{DB: db}
	router := gin.Default()
	router.GET("/products", productHandlers.GetProducts)
	router.POST("/payment", paymentHandlers.NewPayment)
	router.POST("/payment/callback", paymentHandlers.PaymentCallback)

	_ = router.Run(":8081")
}
