package api

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func configPProfRouter(router gin.IRouter) {
	pprof.Register(router.(*gin.Engine))
}
