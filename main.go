package main

import (
	docs "github.com/TechBuilder-360/portfolio-v2-backend/docs"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/config"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database/redis"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/middleware"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/router"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

// @title           ORIS API
// @version         1.0
// @description     This is oris server api.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Oris API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  techbuilder.circle@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	config.Load()
	configureSwagger()
	Init()

	if config.Instance.GetEnv() != config.SandboxEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	dbConnection := database.ConnectDB(config.Instance.DbURL)
	err := database.MigrateAll(dbConnection)
	if err != nil {
		log.Error(err.Error())
	}
	//  set connection pool
	sqlDB, _ := dbConnection.DB()

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour)

	defer func() {
		err := sqlDB.Close()
		if err != nil {
			log.Error("error: %v", err)
			return
		}
	}()

	redis.NewClient()

	r := gin.New()
	setupSystemRouteHandler(r)
	router.SetUpRoutes(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func setupSystemRouteHandler(router *gin.Engine) {
	router.Use(middleware.CORS())
	router.NoMethod(middleware.Http405Handler())
	router.NoRoute(middleware.Http404Handler())
	router.Use(gin.CustomRecovery(middleware.Http500Handler))
}

func configureSwagger() {
	docs.SwaggerInfo.Title = "ORIS API"
	docs.SwaggerInfo.Description = "This documentation contains url path description for Oris APIs"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = config.Instance.Host
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
