package api

import (
	"github.com/gin-gonic/gin"
)

func ConfigRouter(router gin.IRouter) {
	configHealthzRouter(router)
	configMetricsRouter(router)
	configPProfRouter(router)

	configDemoRouter(router)
}
