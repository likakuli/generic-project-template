package server

import (
	"context"
	"time"

	"github.com/likakuli/generic-project-template/pkg/api"

	"github.com/gin-gonic/gin"
)

type Server struct {
	port   string
	engine *gin.Engine
}

func NewServer(port string, mode string, qps int) *Server {
	engine := gin.Default()

	gin.SetMode(mode)

	// middlewares
	engine.Use(gin.Recovery())
	engine.Use(leakyBucketRateLimiter(qps))

	// routers
	api.ConfigRouter(engine)

	return &Server{
		port:   port,
		engine: engine,
	}
}

func (s *Server) Run(ctx context.Context) {
	go s.engine.Run(":" + s.port)

	<-ctx.Done()
	time.Sleep(1 * time.Second)
}
