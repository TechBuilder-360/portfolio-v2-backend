package controller

import (
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/constant"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/log"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/types"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/common/util"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/response"
	"github.com/TechBuilder-360/portfolio-v2-backend/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IAuthController interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
	RegisterRoutes(router *gin.Engine)
}

type authController struct {
	as service.IAuthService
}

// NewAuthController instantiates Auth Controller
func NewAuthController() IAuthController {
	return &authController{
		as: service.NewAuthService(),
	}
}

func (ctl *authController) RegisterRoutes(router *gin.Engine) {

	v1 := router.Group("/v1")
	auth := v1.Group("/auth")

	auth.POST("/register", ctl.Register)
	auth.POST("/login", ctl.Login)
}

// Register godoc
// @Summary      registration
// @Description  Register User
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data   body    types.Authentication  true  "register"
// @Success      200  {object}  response.SuccessResp{Data=types.RegisterResponse}
// @Router       /auth/register [post]
func (ctl *authController) Register(ctx *gin.Context) {
	var body types.Authentication

	logger := log.WithFields(log.FromContext(ctx).Fields)
	requestIdentifier := util.GenerateUUID()
	//logger = logger.WithField(util.UserIdentifier, util.ExtractUserIdContext(ctx)).WithField(util.RequestIdentifier, requestIdentifier)
	logger.Info("Register request")

	ctx.Header(constant.RequestIdentifier, requestIdentifier)

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("error while parsing request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationErrorResponse(util.CustomErrorResponse(err)))
		return
	}

	res, err := ctl.as.RegisterUser(ctx, body, logger)

	if err != nil {
		logger.Error("error while registering: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.DataResponse("Registration successful", res))
}

// Login godoc
// @Summary      Login
// @Description  Login User
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        data   body    types.Authentication  true  "login"
// @Success      200  {object}  response.SuccessResp{Data=types.LoginResponse}
// @Router       /auth/login [post]
func (ctl *authController) Login(ctx *gin.Context) {
	var body types.Authentication

	logger := log.WithFields(log.FromContext(ctx).Fields)
	requestIdentifier := util.GenerateUUID()
	logger.Info("Register request")

	ctx.Header(constant.RequestIdentifier, requestIdentifier)

	if err := ctx.ShouldBindJSON(&body); err != nil {
		logger.Error("error while parsing request body: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ValidationErrorResponse(util.CustomErrorResponse(err)))
		return
	}

	res, err := ctl.as.Login(ctx, body, logger)

	if err != nil {
		logger.Error("error while login: %v", err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, response.ErrorResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, response.DataResponse("Successful", res))
}
