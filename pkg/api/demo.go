package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leopoldxx/go-utils/trace"
)

func configDemoRouter(router gin.IRouter) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/demo",
			Trace("demo"),
			GetDemo)
	}
}

func GetDemo(c *gin.Context) {
	tracer := trace.GetTraceFromRequest(c.Request)

	time.Sleep(10 * time.Millisecond)
	tracer.Info("get demo api called")

}
