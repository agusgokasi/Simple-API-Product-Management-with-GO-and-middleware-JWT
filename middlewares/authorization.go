package middlewares

import (
	"eleventh-learn/database"
	"eleventh-learn/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func ProductAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDB()
		productId, err := strconv.Atoi(c.Param("productId"))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Bad Request",
				"message": err.Error(),
			})

			return
		}

		userData := c.MustGet("userData").(jwt.MapClaims)
		userID := uint(userData["id"].(float64))
		userIsAdmin := userData["is_admin"].(bool) // added line to get user's admin status
		product := models.Product{}

		err = db.Select("user_id").First(&product, uint(productId)).Error
		if err != nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"error":   "not found",
				"message": err.Error(),
			})

			return
		}

		if product.UserID != userID && !userIsAdmin { // check if user is admin or the owner of the product
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "you're not allowed to access this",
			})

			return
		}

		c.Next()
	}
}
