package util

import (
	"fmt"
	"github.com/alif-github/task-management/config"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
)

func WriteSuccessResponse(reqID, msg string, data interface{}) domain.Response {
	return domain.Response{
		RequestID: reqID,
		Status:    true,
		Message:   msg,
		Data:      data,
	}
}

func WriteErrorResponse(reqID, msg string) domain.Response {
	return domain.Response{
		RequestID: reqID,
		Status:    false,
		Message:   msg,
		Data:      nil,
	}
}

func WriteLogStdout(errs error, err *util.ErrorModel, requestID, handlerName string) {
	if errs != nil {
		logModel := util.LoggerModel{
			RequestID:   requestID,
			Class:       handlerName,
			Application: config.ApplicationConfiguration.GetServer().Application,
			Version:     config.ApplicationConfiguration.GetServer().Version,
			Code:        http.StatusInternalServerError,
			Message:     errs.Error(),
		}
		util.LogError(logModel.LoggerZapFieldObject())
		return
	}

	if err.Err != nil {
		logModel := util.LoggerModel{
			RequestID:   requestID,
			Class:       fmt.Sprintf(`[%s, %s]`, err.FileName, err.FuncName),
			Application: config.ApplicationConfiguration.GetServer().Application,
			Version:     config.ApplicationConfiguration.GetServer().Version,
			Code:        err.Code,
			Message:     err.Err.Error(),
		}
		util.LogError(logModel.LoggerZapFieldObject())
		return
	}

	logModel := util.LoggerModel{
		RequestID:   requestID,
		Class:       handlerName,
		Application: config.ApplicationConfiguration.GetServer().Application,
		Version:     config.ApplicationConfiguration.GetServer().Version,
	}
	util.LogInfo(logModel.LoggerZapFieldObject())
}

func SetJWTTokenCookie(c *gin.Context, token string, claims jwt.MapClaims) {
	cookie := &http.Cookie{
		Name:     "fit_token",
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		MaxAge:   int(claims["exp"].(int64)),
	}
	http.SetCookie(c.Writer, cookie)
}
