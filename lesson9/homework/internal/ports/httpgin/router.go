package httpgin

import (
	"homework9/internal/app"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

func AppRouter(r *gin.RouterGroup, a app.App) {

	r.Use(gin.CustomRecovery(CustomPanicRecover))
	r.Use(CustomLogger)

	r.GET("/ads/:id", GetAdByID(a))
	r.POST("/ads", CreateAd(a))
	r.PUT("/ads/:id/status", ChangeAdStatus(a))
	r.PUT("/ads/:id", UpdateAd(a))
	r.GET("/ads", Select(a))
	r.GET("/ads/title", FindAdByTitle(a))

	r.POST("/users", CreateUser(a))
	r.POST("/users/:id", UpdateUser(a))

}

func CustomLogger(c *gin.Context) {
	t := time.Now().UTC()
	c.Next()
	latency := time.Since(t)
	status := c.Writer.Status()

	log.Println("latency", latency, "method", c.Request.Method, "path", c.Request.URL.Path, "status", status)
}

func CustomPanicRecover(c *gin.Context, err any) {
	log.Println("panic: " + err.(error).Error())
	buf := make([]byte, 2048)
	n := runtime.Stack(buf, false)
	log.Println(string(buf[:n]))
	c.AbortWithStatusJSON(http.StatusInternalServerError, err.(error).Error())
}
