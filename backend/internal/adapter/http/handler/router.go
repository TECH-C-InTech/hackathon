package handler

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// NewRouter は HTTP ハンドラーを紐づけた gin.Engine を返す。
func NewRouter(drawHandler *DrawHandler, postHandler *PostHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// CORS設定
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	corsOrigins := os.Getenv("CORS_ALLOW_ORIGINS")
	if corsOrigins == "" {
		config.AllowOrigins = []string{"http://localhost:3000"}
	} else {
		origins := strings.Split(corsOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}
		config.AllowOrigins = origins
	}

	router.Use(cors.New(config))

	router.GET("/draws/random", drawHandler.GetRandomDraw)
	router.POST("/posts", postHandler.CreatePost)

	return router
}
