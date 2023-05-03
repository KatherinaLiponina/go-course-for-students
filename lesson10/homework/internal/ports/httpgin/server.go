package httpgin

import (
	"net/http"

	"homework10/internal/app"

	"github.com/gin-gonic/gin"
)

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	api := handler.Group("/api/v1")
	AppRouter(api, a)
	s := &http.Server{Addr: port, Handler: handler}

	return s
}
