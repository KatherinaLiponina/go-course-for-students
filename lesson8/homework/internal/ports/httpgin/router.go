package httpgin

import (
	"homework8/internal/app"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a app.App) {

	r.Use(gin.Recovery())
	r.Use(CustomMiddleware)

	r.GET("/ads/:id", GetAdByID(a))
	r.POST("/ads", CreateAd(a))
	r.PUT("/ads/:id/status", ChangeAdStatus(a))	
	r.PUT("/ads/:id", UpdateAd(a))
	r.GET("/ads", Select(a))
	r.GET("/ads/title", FindAdByTitle(a))

	r.POST("/users", CreateUser(a))
	r.POST("/users/:id", UpdateUser(a))

}

func CustomMiddleware(c * gin.Context) {
	t := time.Now().UTC()
	c.Next()
	latency := time.Since(t)
	status := c.Writer.Status()

	log.Println("latency", latency, "method", c.Request.Method, "path", c.Request.URL.Path, "status", status)
}