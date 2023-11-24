package handler

import (
	"fmt"
	"github.com/alif-github/task-management/config"
	delivery_helper "github.com/alif-github/task-management/delivery/util"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type UsersHandler struct {
	userUsecase domain.UserUsecase
}

func NewUsersHandler(u domain.UserUsecase) UsersHandler {
	return UsersHandler{u}
}

func (input UsersHandler) LogoutHandler(c *gin.Context) {
	cookies := c.Request.Cookies()

	var (
		fileName  = "user_handler.go"
		funcName  = "LogoutHandler"
		jwtToken  *http.Cookie
		requestID string
	)

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)

	defer delivery_helper.WriteLogStdout(nil, &util.ErrorModel{}, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	for _, itemCookie := range cookies {
		if itemCookie.Name == "fit_token" {
			jwtToken = itemCookie
			break
		}
	}

	if jwtToken != nil {
		http.SetCookie(c.Writer, &http.Cookie{
			Name:     "fit_token",
			Value:    "",
			HttpOnly: true,
			Path:     "/",
			MaxAge:   int(time.Now().Add(-1 * time.Hour).Unix()),
		})
		c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Logout Success!!", nil))
	} else {
		c.JSON(http.StatusNotFound, delivery_helper.WriteErrorResponse(requestID, "Cookies Not Found/Unauthorized"))
	}
}

func (input UsersHandler) LoginHandler(c *gin.Context) {
	var (
		fileName  = "user_handler.go"
		funcName  = "LoginHandler"
		login     domain.Login
		errs      error
		err       util.ErrorModel
		requestID string
	)

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	if errs = c.ShouldBindJSON(&login); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	defer delivery_helper.WriteLogStdout(errs, &util.ErrorModel{}, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	m := make(map[string]string)
	m[login.Username] = login.Password

	res, err := input.userUsecase.GetLogin(c, &domain.User{Username: login.Username})
	if err.Err != nil {
		return
	}

	if util.CheckPassword(login.Password, res.Password) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = res.ID
		claims["username"] = login.Username
		claims["permission"] = res.Role
		claims["exp"] = time.Now().Add(30 * time.Minute).Unix()

		var (
			tokenStr  string
			signature = config.ApplicationConfiguration.GetServer().SignatureKey
		)

		tokenStr, errs = token.SignedString([]byte(signature))
		if errs != nil {
			c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Failed to create token!"))
			return
		}

		delivery_helper.SetJWTTokenCookie(c, tokenStr, claims)
		c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Token has generated!", nil))
	} else {
		c.JSON(http.StatusUnauthorized, delivery_helper.WriteErrorResponse(requestID, "Invalid credential!"))
	}
}

func (input UsersHandler) RegisterHandler(c *gin.Context) {
	var (
		fileName  = "user_handler.go"
		funcName  = "RegisterHandler"
		user      domain.User
		requestID string
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	if errs = c.ShouldBindJSON(&user); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	err = input.userUsecase.Store(c, &user)
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Register!", nil))
}

func (input UsersHandler) UpdateUserHandler(c *gin.Context) {
	var (
		fileName  = "user_handler.go"
		funcName  = "UpdateUserHandler"
		requestID string
		users     domain.User
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	_, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Claims not found"))
		return
	}

	idParam := c.Param("id")
	id, errs := strconv.Atoi(idParam)
	if errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	if errs = c.ShouldBindJSON(&users); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	users.ID = int64(id)
	err = input.userUsecase.Update(c, &users)
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Update!", nil))
}

func (input UsersHandler) DeleteUserHandler(c *gin.Context) {
	var (
		fileName  = "user_handler.go"
		funcName  = "DeleteUserHandler"
		requestID string
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	_, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Claims not found"))
		return
	}

	idParam := c.Param("id")
	id, errs := strconv.Atoi(idParam)
	if errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	err = input.userUsecase.Delete(c, &domain.User{ID: int64(id)})
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Delete!", nil))
}

func (input UsersHandler) FetchUserHandler(c *gin.Context) {
	var (
		fileName   = "user_handler.go"
		funcName   = "FetchUserHandler"
		user       []domain.User
		pagination domain.Pagination
		requestID  string
		errs       error
		err        util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	_, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Claims not found"))
		return
	}

	if errs = c.ShouldBindQuery(&pagination); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	user, err = input.userUsecase.Fetch(c, pagination.Page, pagination.Limit)
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Get Detail!", user))
}

func (input UsersHandler) ViewUserHandler(c *gin.Context) {
	var (
		fileName  = "user_handler.go"
		funcName  = "LoginHandler"
		user      domain.User
		requestID string
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	_, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Claims not found"))
		return
	}

	idParam := c.Param("id")
	id, errs := strconv.Atoi(idParam)
	if errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	user, err = input.userUsecase.GetByID(c, int64(id))
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Get Detail!", user))
}
