package controller

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/config"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type IGeneralController interface {
	Health(c *gin.Context)
	Welcome(c *gin.Context)
	RegisterRoutes(router *gin.Engine)
}

type generalController struct {
}

func NewGeneralController() IGeneralController {
	return &generalController{}
}

type WelcomeResponse struct {
	Service string    `json:"service"`
	Version string    `json:"version"`
	Env     string    `json:"env"`
	Date    time.Time `json:"date"`
}

func (ctl *generalController) RegisterRoutes(router *gin.Engine) {

	router.GET("/", ctl.Welcome)
	router.GET("/ping", ctl.Health)

	//v1 := router.Group("/v1")
}

func (ctl *generalController) Health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "I'm all the way up 🚀",
	})
}

func (ctl *generalController) Welcome(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.SuccessResponse("Welcome to Oris API",
		&WelcomeResponse{Service: "oris-api", Env: config.Instance.Env, Version: "1.0.0", Date: time.Now()},
		nil),
	)
}
