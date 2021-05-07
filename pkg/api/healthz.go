package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func configHealthzRouter(router gin.IRouter) {
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}
