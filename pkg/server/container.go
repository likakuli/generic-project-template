package server

import (
	"net/http"
	"sync"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"

	"github.com/likakuli/generic-project-template/pkg/config"
	"github.com/likakuli/generic-project-template/pkg/controllers"
	"github.com/likakuli/generic-project-template/pkg/repositories"
	"github.com/likakuli/generic-project-template/pkg/services"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	once sync.Once
	c    *container

	// todo: update the default value as you expected
	defaultSubsystem = "generic_project_template"
)

func open(connString string, maxOpenConn int, maxIdleConn int) *gorm.DB {
	db, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := db.DB()
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxOpenConns(maxOpenConn)
	sqlDB.SetMaxIdleConns(maxIdleConn)

	return db
}

func defaultControllerContainer(cfg *config.DBConfig) *container {
	if c == nil {
		once.Do(func() {
			c = &container{
				db: open(cfg.ConnectionString, cfg.MaxOpenConn, cfg.MaxIdleConn),
			}
		})
	}
	return c
}

func (c *container) Dispose() {
	if c.db != nil {
		sqlDB, _ := c.db.DB()
		sqlDB.Close()
	}
}

// IContainer provider the method to get all needed controllers
type IContainer interface {
	// Dispose all resources when server stopped
	Dispose()
	GetPlayerController() *controllers.PlayerController
	// todo: add your providers here
}

type container struct {
	db *gorm.DB
}

func configHealthzRouter(router gin.IRouter, container IContainer) {
	router.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})
}

func configPProfRouter(router gin.IRouter, container IContainer) {
	pprof.Register(router.(*gin.Engine))
}

func configMetricsRouter(router gin.IRouter, container IContainer) {
	// todo: add your custom prometheus metrics here if exists
	controllers.NewPrometheusRouteRegister(defaultSubsystem).ConfigRoute(router.(*gin.Engine))
}

func (c *container) GetPlayerController() *controllers.PlayerController {
	return controllers.NewPlayerController(
		services.NewPlayerService(
			repositories.NewPlayerRepository(c.db)))
}
