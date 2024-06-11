package main

import (
	"log"

	"solution1/session-2-latihan-crud-user-gin/router"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	router.SetupRouter(r)

	log.Println("Running server on port 8080")
	r.Run(":8080")
}

// import (
// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	r := gin.Default()

// 	r.GET("/", func(c *gin.Context) {
// 		c.JSON(200, gin.H{
// 			"message": "Halo dari Gin!",
// 		})
// 	})

// 	r.POST("/post", func(c *gin.Context) {
// 		var json struct {
// 			Message string `json:"message"`
// 		}
// 		if err := c.ShouldBindJSON(&json); err == nil {
// 			c.JSON(200, gin.H{"Message": json.Message})
// 		} else {
// 			c.JSON(400, gin.H{"error": err.Error()})
// 		}
// 	})

// 	r.Run(":8080")
// }
