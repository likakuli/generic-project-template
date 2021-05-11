package server

import (
	"github.com/gin-gonic/gin"

	"github.com/likakuli/generic-project-template/pkg/server/middlewares"
)

type RouteRegister func(router gin.IRouter, container IContainer)

var (
	routers []RouteRegister
)

func init() {
	routers = []RouteRegister{
		configHealthzRouter,
		configPProfRouter,
		configMetricsRouter,
		// todo: add your biz routers
		configPlayerRouter,
	}
}

// register all expected routers
func configRouter(router gin.IRouter, container IContainer) {
	for _, register := range routers {
		register(router, container)
	}
}

func configPlayerRouter(router gin.IRouter, container IContainer) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/score/:player1/vs/:player2",
			middlewares.Trace("player"),
			container.GetPlayerController().GetPlayerScore)
	}
}
