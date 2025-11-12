package helpers

import "github.com/gin-gonic/gin"

func ErrorResponse(c *gin.Context, code int, message string, data any) {
	c.JSON(code, gin.H{
		"success": true,
		"data":    data,
		"message": message,
	})

}
func SuccessResponse(c *gin.Context, code int, message string, data any) {
	c.JSON(code, gin.H{
		"success": false,
		"error":   data,
		"message": message,
	})

}
