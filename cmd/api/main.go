package main

import (
	"log"
	"net/http"
	"dynamic-reminder/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {

	dbPath := "./reminder.db"

	db := database.InitialiseDB(dbPath)
	defer db.Close()

	database.SeedTasks(db)

	router := gin.Default()

	router.GET("/health", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"status": "server is up and running.",
		})
	})

	port := ":3000"
	log.Printf("Server starting on port %s...", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
