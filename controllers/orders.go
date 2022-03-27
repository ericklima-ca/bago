package controllers

import (
	"net/http"

	"github.com/ericklima-ca/bago/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func (oc *OrderController) CreateMany(c *gin.Context) {
	var orders []models.Order
	if err := c.ShouldBindJSON(&orders); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"ok": false,
			"body": gin.H{
				"error": "data does not match",
			},
		})
		return
	}
	if r := oc.DB.Create(orders); r.Error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"ok": false,
			"body": gin.H{
				"error": r.Error,
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"body": gin.H{
			"msg":    "orders created successfully",
			"orders": orders,
		},
	})
	return
}

func (oc *OrderController) GetAll(c *gin.Context) {
	var orders []models.Order
	oc.DB.Preload("PurchaseOrders.Status").Preload("PurchaseOrders.Sells.Product").Find(&orders)
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"body": gin.H{
			"orders": orders,
		},
	})
}
