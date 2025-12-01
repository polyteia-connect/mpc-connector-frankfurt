package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/vaccinated-ids", func(c *gin.Context) {
		time.Sleep(3 * time.Second)
		c.JSON(200, gin.H{
			"fileStateIds": []string{"6fe386ea-e15c-4e10-8d9e-046db466cdac"},
		})
	})
	r.Run(":3000")
}
