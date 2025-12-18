package handler

import "github.com/gin-gonic/gin"

// NewRouter はおみくじハンドラを紐づけた gin.Engine を返す。
func NewRouter(drawHandler *DrawHandler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/draws/random", drawHandler.GetRandomDraw)

	return router
}
