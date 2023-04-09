package controllers

import (
	"eleventh-learn/database"
	"eleventh-learn/helpers"
	"eleventh-learn/models"
	"math"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)

	product := models.Product{}
	userID := uint(userData["id"].(float64))

	if contentType == appJSON {
		c.ShouldBindJSON(&product)
	} else {
		c.ShouldBind(&product)
	}

	product.UserID = userID

	err := db.Debug().Create(&product).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "bad request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, product)

}

func UpdateProduct(c *gin.Context) {
	db := database.GetDB()
	contentType := helpers.GetContentType(c)
	userData := c.MustGet("userData").(jwt.MapClaims)

	productId, _ := strconv.Atoi(c.Param("productId"))
	userID := uint(userData["id"].(float64))

	product := models.Product{}

	if contentType == appJSON {
		c.ShouldBindJSON(&product)
	} else {
		c.ShouldBind(&product)
	}

	product.UpdatedBy = userID

	err := db.Model(&product).Where("id = ?", productId).Updates(models.Product{
		Title:       product.Title,
		Description: product.Description,
	}).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "bad request",
			"message": err.Error(),
		})

		return
	}

	// Query the updated product by its ID
	updatedProduct := models.Product{}
	err = db.First(&updatedProduct, productId).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "bad request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, updatedProduct)
}

func DeleteProduct(c *gin.Context) {
	db := database.GetDB()

	productId, _ := strconv.Atoi(c.Param("productId"))

	product := models.Product{}

	err := db.Model(&product).Where("id = ?", productId).Delete(&product).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad request",
			"message": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

func GetProductByID(c *gin.Context) {
	db := database.GetDB()

	productId, _ := strconv.Atoi(c.Param("productId"))

	product := models.Product{}
	err := db.Preload("User").Where("id = ?", productId).First(&product).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"err":     "bad request",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func GetAllProducts(c *gin.Context) {
	db := database.GetDB()

	var products []models.Product
	var count int64

	// get query parameters
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}

	// calculate offset
	offset := (page - 1) * limit

	// get total count of products
	db.Model(&models.Product{}).Count(&count)

	// get products with pagination
	db.Preload("User").Offset(offset).Limit(limit).Find(&products)

	// prepare response data
	responseData := gin.H{
		"products": products,
		"meta": gin.H{
			"total_data":   count,
			"current_page": page,
			"limit":        limit,
			"total_page":   int(math.Ceil(float64(count) / float64(limit))),
		},
	}

	c.JSON(http.StatusOK, responseData)
}
