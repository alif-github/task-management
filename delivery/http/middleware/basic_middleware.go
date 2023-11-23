package middleware

import (
	"errors"
	"github.com/alif-github/clean-arch/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BasicMiddleware(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	c.Header("Content-Type", "application/json")
	defer func() {
		if errs := recover(); errs != nil {
			r := errors.New(errs.(string))
			util.LogError(util.DefaultGenerateLogModel(500, r.Error()).LoggerZapFieldObject())
		}
	}()

	c.Set("request_id", uuid.New().String())
	c.Next()
}
