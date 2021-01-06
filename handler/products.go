package handler

import (
	"github.com/FlexHC/MinecraftStore/model"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"
)

type ProductHandlers struct {
	DB *sqlx.DB
}

func (p *ProductHandlers) GetProducts(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	products := make([]model.Product, 0)
	err := p.DB.Select(&products, "SELECT * FROM `products`")
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "database error",
		})
		return
	}
	c.JSON(http.StatusOK, products)
	return
}
