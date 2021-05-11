package server

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"

	"github.com/likakuli/generic-project-template/pkg/config"
	"github.com/likakuli/generic-project-template/pkg/server/middlewares"
)

type Server struct {
	port      string
	engine    *gin.Engine
	container IContainer
}

func NewServer(cfg *config.Config) *Server {
	engine := gin.Default()

	gin.SetMode(cfg.Server.Mode)

	// use middlewares
	engine.Use(gin.Recovery())
	engine.Use(middlewares.LeakyBucketRateLimiter(cfg.Server.API_QPS))

	// init container
	container := defaultControllerContainer(cfg.DB)
	// register routers
	configRouter(engine, container)

	return &Server{
		port:      cfg.Server.Port,
		engine:    engine,
		container: container,
	}
}

func (s *Server) Run(ctx context.Context) {
	go s.engine.Run(":" + s.port)
	glog.Info("server started..")
	defer s.container.Dispose()

	<-ctx.Done()
	glog.Infof("server stopped..")
	time.Sleep(1 * time.Second)
}
