package middleware

import (
	"fmt"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/response"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func Http500Handler(c *gin.Context, recovered interface{}) {
	if err, ok := recovered.(string); ok {
		log.Error("error: %s", err)
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError,
		response.ErrorResponse("internal server error"),
	)
}

func Http404Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.AbortWithStatusJSON(http.StatusNotFound,
			response.ErrorResponse(fmt.Sprintf("Path %s not found", c.Request.URL.Path)),
		)
	}
}

func Http405Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.AbortWithStatusJSON(http.StatusMethodNotAllowed,
			response.ErrorResponse(fmt.Sprintf("Method Not Allowed for path %s", c.Request.URL.Path)),
		)
	}
}
