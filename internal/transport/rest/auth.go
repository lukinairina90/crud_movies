package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lukinairina90/crud_movies/internal/domain"
	"github.com/sirupsen/logrus"
)

type UserService interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Auth struct {
	userService UserService
}

func NewAuth(userService UserService) *Auth {
	return &Auth{userService: userService}
}

func (a *Auth) InjectRoutes(r *gin.Engine, middlewares ...gin.HandlerFunc) {
	auth := r.Group("/auth").Use(middlewares...)
	{
		auth.POST("/sign-up", a.signUp)
		auth.POST("/sign-in", a.signIn)
		auth.GET("/refresh", a.refresh)
	}
}

// @Summary SignUp
// @Tags auth
// @Description create account
// @ID create-account
// @Accept  json
// @Produce  json
// @Param input body domain.SignUpInput true "account info"
// @Success 200 {object} domain.SignUpInput
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} BadRequestErr
// @Failure default {object} BadRequestErr
// @Router /auth/sign-up [post]
func (a *Auth) signUp(ctx *gin.Context) {
	var inp domain.SignUpInput

	if err := ctx.BindJSON(&inp); err != nil {
		logError("signUp", err)
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("cannot parse body user", nil))
		return
	}

	if err := inp.Validate(); err != nil {
		vErrs := err.(validator.ValidationErrors)
		errs := make(map[string]string)
		for _, fErr := range vErrs {
			errs[fErr.ActualTag()] = fErr.Error()
		}

		logError("signUp", err)
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("validation error", errs))
		return
	}

	err := a.userService.SignUp(ctx, inp)
	if err != nil {
		logError("signUp", err)
		ctx.JSON(http.StatusInternalServerError, NewInternalServerErr("userService.SignUp error"))
		return
	}

	ctx.JSON(http.StatusOK, inp)
}

// @Summary SignIn
// @Tags auth
// @Description login
// @ID login
// @Accept  json
// @Produce  json
// @Param input body domain.SignInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} BadRequestErr
// @Failure default {object} BadRequestErr
// @Router /auth/sign-in [post]
func (a *Auth) signIn(ctx *gin.Context) {
	var inp domain.SignInInput
	if err := ctx.BindJSON(&inp); err != nil {
		logError("signIn", err)
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("connot parse body user", nil))
	}

	if err := inp.Validate(); err != nil {
		vErrs := err.(validator.ValidationErrors)
		errs := make(map[string]string)
		for _, fErr := range vErrs {
			errs[fErr.ActualTag()] = fErr.Error()
		}

		logError("signIn", err)
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("validation error", errs))
	}

	accessToken, refreshToken, err := a.userService.SignIn(ctx, inp)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			HandleNotFoundError(ctx, err)
			return
		}

		logError("signIn", err)
		ctx.JSON(http.StatusInternalServerError, NewBadRequestErr("userService.SignUp error", nil))
		return
	}

	ctx.SetCookie("refresh-token", refreshToken, 3600, "/auth", "localhost", false, true)

	ctx.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}

// @Summary refresh
// @Tags auth
// @Description returns accessToken and sets in cookies refresh-token
// @ID refresh
// @Accept  json
// @Produce  json
// @Header 200 {string} Token "token"
// @Success 200 {string} string
// @Failure 400,404 {object} BadRequestErr
// @Failure 500 {object} BadRequestErr
// @Failure default {object} BadRequestErr
// @Router /auth/refresh [get]
func (a *Auth) refresh(ctx *gin.Context) {
	cookie, err := ctx.Cookie("refresh-token")
	if err != nil {
		logError("refresh", err)
		ctx.JSON(http.StatusBadRequest, NewBadRequestErr("get cookie from request error", nil))
		return
	}
	logrus.Info("%s", cookie)

	accessToken, refreshToken, err := a.userService.RefreshTokens(ctx, cookie)
	if err != nil {
		logError("SignIn | refresh", err)
		ctx.JSON(http.StatusInternalServerError, NewBadRequestErr("refresh token error", nil))
		return
	}

	ctx.SetCookie("refresh-token", refreshToken, 3600, "/auth", "localhost", false, true)

	ctx.JSON(http.StatusOK, map[string]string{
		"token": accessToken,
	})
}
