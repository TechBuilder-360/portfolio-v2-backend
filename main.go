package main

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/config"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/database/redis"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/middleware"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/router"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/gin-swagger/example/basic/docs"
	"os"
	"time"
)

func Init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)

	config.Load()
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

	configureSwagger()
	Init()

	docs.SwaggerInfo_swagger.Host = config.Instance.Host

	dbConnection := database.ConnectDB(config.Instance.DbURL)
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

	redis.NewClient(config.Instance.RedisURL, config.Instance.RedisPassword, config.Instance.Namespace, nil)

	r := gin.New()
	r.Use(gin.Recovery())
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
	docs.SwaggerInfo_swagger.Title = "ORIS API"
	docs.SwaggerInfo_swagger.Description = "This documentation contains url path description for Oris APIs"
	docs.SwaggerInfo_swagger.Version = "1.0"
	docs.SwaggerInfo_swagger.Host = config.Instance.Host
	docs.SwaggerInfo_swagger.BasePath = "/v1"
	docs.SwaggerInfo_swagger.Schemes = []string{"http", "https"}
}
