package rest

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrMovieNotFound       = errors.New("book not found")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type NotFoundErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewNotFoundErr(message string) NotFoundErr {
	return NotFoundErr{Code: http.StatusNotFound, Message: message}
}

type BadRequestErr struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Fields  map[string]string `json:"fields,omitempty"`
}

func NewBadRequestErr(mgs string, fields map[string]string) BadRequestErr {
	return BadRequestErr{Code: http.StatusBadRequest, Message: mgs, Fields: fields}
}

type InternalServerErr struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func NewInternalServerErr(message string) InternalServerErr {
	return InternalServerErr{Code: http.StatusInternalServerError, Message: message}
}

type UnauthorizedErr struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
}

func NewUnauthorizedErr(message string) *UnauthorizedErr {
	return &UnauthorizedErr{
		Message: message,
		Code:    http.StatusUnauthorized,
	}
}

func HandleNotFoundError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, map[string]string{
		"error": err.Error(),
	})
}
