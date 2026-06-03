package main

import (
	"log"
	"net/http"
	"time"

	"go-sample-http/config"
	"go-sample-http/internal/db"
	"go-sample-http/internal/item"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)

	database, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatal("db init failed:", err)
	}
	defer database.Close()

	// 레이어 조립: repository → usecase → handler
	itemRepo := item.NewRepository(database)
	itemUseCase := item.NewUseCase(itemRepo)
	itemHandler := item.NewHandler(itemUseCase)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().UTC().Format(time.RFC3339),
		})
	})

	itemHandler.RegisterRoutes(r)

	r.Run(":" + cfg.Port)
}
