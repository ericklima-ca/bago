package controllers

import (
	"net/http"
	"strconv"

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

	oc.DB.Scopes(paginate(c)).Preload("PurchaseOrders.Status").
		Preload("PurchaseOrders.Sells.Product").
		Order("id desc").
		Find(&orders)
	
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
		"body": gin.H{
			"orders": orders,
		},
	})
}
// to be verified
func (oc *OrderController) Delete(c *gin.Context) {
	if value, _ := c.Get("user"); value.(models.User).Role != "admin" {

	}
	orderID, _ := strconv.Atoi(c.Param("id"))
	var order models.Order
	oc.DB.Where("id = ?", orderID).Delete(&order)
	c.JSON(http.StatusAccepted, gin.H{
		"ok": true,
		"body": gin.H {
			"msg": "order deleted",
		},
	})
}

func (oc *OrderController) Update(c *gin.Context) {

}

func paginate(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	return func (db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(c.Query("page"))
		if page == 0 {
		page = 1
		}

		pageSize, _ := strconv.Atoi(c.Query("page_size"))
		if pageSize == 0 {
		page = 20
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}