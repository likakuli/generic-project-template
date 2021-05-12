package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leopoldxx/go-utils/httputils"
	"github.com/leopoldxx/go-utils/trace"

	"github.com/likakuli/generic-project-template/pkg/interfaces"
	"github.com/likakuli/generic-project-template/pkg/viewmodels"
)

type PlayerControllerOption func(controller *PlayerController)

type PlayerController struct {
	service    interfaces.IPlayerService
	restCli    *httputils.RestCli
	debugLevel httputils.DebugLevel
	// todo: add more
}

func (controller *PlayerController) GetPlayerScore(c *gin.Context) {
	tracer := trace.GetTraceFromRequest(c.Request)

	player1, _ := parseString(c.Params, "player1")
	player2, _ := parseString(c.Params, "player2")

	scores, err := controller.service.GetScores(c.Request.Context(), player1, player2)
	if err != nil {
		tracer.Infof("GetScores failed with err: %s", err.Error())
		returnError(c, http.StatusInternalServerError, err)
	}
	tracer.Info("GetScores succeed!")

	c.JSON(http.StatusOK, viewmodels.ScoreVM{Score: scores})
}

func NewPlayerController(service interfaces.IPlayerService, options ...PlayerControllerOption) *PlayerController {
	controller := &PlayerController{
		service: service,
	}

	for _, option := range options {
		option(controller)
	}

	return controller
}
