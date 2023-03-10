package router

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/config"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/controller"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes(
	router *gin.Engine,

) {
	var (
		generalController = controller.NewGeneralController()
		authController    = controller.NewAuthController()
	)

	// Setup routes
	generalController.RegisterRoutes(router)

	v1 := router.Group("/v1")

	// *****************
	// Auth  ROUTES
	// *****************
	authController.RegisterRoutes(router)

	// *****************
	// TEST ROUTES {NOT FOR PRODUCTION}
	// *****************
	if config.Instance.GetEnv() != config.ProductionEnv {
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

}
