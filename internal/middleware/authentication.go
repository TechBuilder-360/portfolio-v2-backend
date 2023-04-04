package middleware

import (
	"errors"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/constant"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/util"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/model"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/repository"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/response"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		token := extractJWT(auth)
		if token == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		as := service.NewAuthService()
		tk, err := as.ValidateToken(util.AddrToString(token))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var claims jwt.MapClaims

		if tk.Valid {
			claims = tk.Claims.(jwt.MapClaims)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims[constant.VerifiedEmail] == false {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse("email not verified"))
			return
		}

		if claims[constant.AccountStatus] == constant.Disabled {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse("account is not active"))
			return
		}
		c.Set(constant.AccountID, claims[constant.AccountID].(string))
		c.Next()
	}
}

// PartialAuth Middleware to authenticate when account is disabled or not verified
func PartialAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		token := extractJWT(auth)
		if token == nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		as := service.NewAuthService()
		tk, err := as.ValidateToken(util.AddrToString(token))
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var claims jwt.MapClaims

		if tk.Valid {
			claims = tk.Claims.(jwt.MapClaims)
		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		if claims[constant.VerifiedEmail] == true || claims[constant.AccountStatus] == constant.Active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("authentication failed"))
			return
		}
		c.Set(constant.AccountID, claims[constant.AccountID].(string))
		c.Next()
	}
}

func ExtractAccount(c *gin.Context) (account *model.Account, err error) {
	// get Account model
	if id, exist := c.Get(constant.AccountID); exist != true {
		return nil, errors.New("account not found")
	} else {
		account, err = repository.NewAccountRepository().GetById(c, id.(string))
		if err != nil {
			return nil, errors.New("account not found")
		}
	}

	return account, nil
}

func extractJWT(auth string) *string {
	if strings.HasPrefix(auth, "Bearer ") {
		token := auth[7:]
		return &token
	}

	return nil
}
