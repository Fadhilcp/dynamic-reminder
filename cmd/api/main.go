package main

import (
	"log"
	"net/http"
	"time"
	"dynamic-reminder/internal/database"
	"dynamic-reminder/internal/handlers"
	"dynamic-reminder/internal/repository"
	"dynamic-reminder/internal/scheduler"
	"github.com/gin-gonic/gin"
)

func main() {

	dbPath := "./reminder.db"

	db := database.InitialiseDB(dbPath)
	defer db.Close()

	database.SeedTasks(db)

	ruleRepo := repository.NewRuleRepository(db)
	ruleHandler := handlers.NewRuleHandler(ruleRepo)

	scheduler.Run(db, 30*time.Second)

	router := gin.Default()

	router.GET("/health", func (c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"status": "server is up and running.",
		})
	})

	rules := router.Group("/rules")
	{
		rules.POST("", ruleHandler.CreateRule)
		rules.GET("", ruleHandler.GetRules)
		rules.PUT("/:ruleId", ruleHandler.UpdateRule)
		rules.PATCH("/:ruleId/toggle", ruleHandler.ToggleRule)
		rules.DELETE("/:ruleId", ruleHandler.DeleteRule)
	}

	port := ":3000"
	log.Printf("Server starting on port %s...", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
