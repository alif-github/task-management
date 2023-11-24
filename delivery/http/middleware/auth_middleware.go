package middleware

import (
	"github.com/alif-github/task-management/config"
	delivery_helper "github.com/alif-github/task-management/delivery/util"
	"github.com/alif-github/task-management/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AuthMiddleware(c *gin.Context) {
	requestID, _ := c.Get("request_id")

	var tokenStr string
	cookies := c.Request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "fit_token" {
			tokenStr = cookie.Value
		}
	}

	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, delivery_helper.WriteErrorResponse(requestID.(string), "Token is missing"))
		c.Abort()
		return
	}

	token, errs := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.ApplicationConfiguration.GetServer().SignatureKey), nil
	})

	if errs != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, delivery_helper.WriteErrorResponse(requestID.(string), "Unauthorized"))
		c.Abort()
		return
	}

	defer func() {
		if errs != nil {
			util.LogError(util.DefaultGenerateLogModel(500, errs.Error()).LoggerZapFieldObject())
		}
	}()

	claims, _ := token.Claims.(jwt.MapClaims)
	c.Set("claims", claims)

	c.Next()
}
