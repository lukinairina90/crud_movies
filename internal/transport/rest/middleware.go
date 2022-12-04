package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

const AuthorizationHeaderName = "Authorization"

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fields := logrus.Fields{
			"method":          c.Request.Method,
			"uri":             c.Request.RequestURI,
			"request-in-time": t.Format(time.RFC3339),
		}

		c.Next()

		dur := time.Since(t)
		fields["request-handling-duration"] = dur.Milliseconds()

		logrus.WithFields(fields).Info()
	}
}

func (a *Auth) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getTokenFromRequest(c)
		if err != nil {
			logError("authMiddleware", err)
			c.JSON(http.StatusUnauthorized, nil)
			return
		}

		uid, err := a.userService.ParseToken(c, token)
		if err != nil {
			logError("authMiddleware", err)
			c.JSON(http.StatusUnauthorized, NewUnauthorizedErr("bad Authorization token"))
			return
		}

		c.Set(fmt.Sprintf("%d", ctxUserID), uid)

		c.Next()
	}
}

func getTokenFromRequest(c *gin.Context) (string, error) {
	header := c.GetHeader(AuthorizationHeaderName)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
