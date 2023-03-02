package router

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/config"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/controller"
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(
	router *gin.Engine,

) {
	var (
		generalController = controller.NewGeneralController()
	)

	// Setup routes
	generalController.RegisterRoutes(router)

	//v1 := router.Group("/v1")

	// *****************
	// XXXXXXXXX  ROUTES
	// *****************

	// *****************
	// TEST ROUTES {NOT FOR PRODUCTION}
	// *****************
	if config.Instance.GetEnv() != config.ProductionEnv {
		//testController.RegisterRoutes(router)
		//v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

}
