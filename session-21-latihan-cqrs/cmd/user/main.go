package main

import (
	"log"
	"net/http"

	"solution1/session-21-latihan-cqrs/entity" // Replace "your-module-name" with the actual module name

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// CreateUser handles POST /users requests
func CreateUser(c *gin.Context, db *gorm.DB) {
	var user entity.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// GetUser handles GET /users/:id requests
func GetUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user entity.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser handles PUT /users/:id requests
func UpdateUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user entity.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user.ID = 0 // Prevent overriding the ID in the database
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /users/:id requests
func DeleteUser(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	if err := db.Delete(&entity.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// ListUsers handles GET /users requests
func ListUsers(c *gin.Context, db *gorm.DB) {
	var users []entity.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func main() {
	// Database connection
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "kafka_practice.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto migrate the User model
	db.AutoMigrate(&entity.User{})

	// Initialize Gin router
	r := gin.Default()

	// Define routes
	r.POST("/users", func(c *gin.Context) {
		CreateUser(c, db)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		GetUser(c, db)
	})
	r.PUT("/users/:id", func(c *gin.Context) {
		UpdateUser(c, db)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		DeleteUser(c, db)
	})
	r.GET("/users", func(c *gin.Context) {
		ListUsers(c, db)
	})

	// Start the server
	r.Run(":8083")
}
