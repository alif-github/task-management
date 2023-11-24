package handler

import (
	"fmt"
	delivery_helper "github.com/alif-github/task-management/delivery/util"
	"github.com/alif-github/task-management/domain"
	"github.com/alif-github/task-management/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TasksHandler struct {
	taskUsecase domain.TaskUsecase
}

func NewTasksHandler(t domain.TaskUsecase) TasksHandler {
	return TasksHandler{t}
}

func (input TasksHandler) StoreTaskHandler(c *gin.Context) {
	var (
		fileName  = "tasks_handler.go"
		funcName  = "StoreTaskHandler"
		task      domain.TaskRequest
		requestID string
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Claims not found"))
		return
	}

	if errs = c.ShouldBindJSON(&task); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	var userID int64
	myClaims, ok := claims.(jwt.MapClaims)
	if ok {
		userIDTemp := myClaims["user_id"]
		userID = int64(userIDTemp.(float64))
	}

	err = input.taskUsecase.Add(c, domain.ContextModel{UserLoginID: userID}, &task)
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Store Data!", nil))
}

func (input TasksHandler) UpdateTaskHandler(c *gin.Context) {
	var (
		fileName  = "tasks_handler.go"
		funcName  = "UpdateTaskHandler"
		requestID string
		task      domain.TaskUpdateRequest
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	claims, exists := c.Get("claims")
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

	if errs = c.ShouldBindJSON(&task); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	var (
		userID     int64
		permission string
		limitedID  int64
	)

	myClaims, ok := claims.(jwt.MapClaims)
	if ok {
		userIDTemp := myClaims["user_id"]
		userID = int64(userIDTemp.(float64))
		permissionTemp := myClaims["permission"]
		permission = permissionTemp.(string)
	}

	if permission == "own" {
		limitedID = userID
	}

	task.ID = int64(id)
	err = input.taskUsecase.Update(c, domain.ContextModel{
		UserLoginID: userID,
		LimitedID:   limitedID,
	}, &task)
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Update!", nil))
}

func (input TasksHandler) DeleteTaskHandler(c *gin.Context) {
	var (
		fileName  = "tasks_handler.go"
		funcName  = "DeleteTaskHandler"
		requestID string
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	claims, exists := c.Get("claims")
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

	var (
		userID     int64
		permission string
		limitedID  int64
	)

	myClaims, ok := claims.(jwt.MapClaims)
	if ok {
		userIDTemp := myClaims["user_id"]
		userID = int64(userIDTemp.(float64))
		permissionTemp := myClaims["permission"]
		permission = permissionTemp.(string)
	}

	if permission == "own" {
		limitedID = userID
	}

	err = input.taskUsecase.Delete(c, domain.ContextModel{
		UserLoginID: userID,
		LimitedID:   limitedID,
	}, &domain.TaskRequest{ID: int64(id)})
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Delete!", nil))
}

func (input TasksHandler) FetchTaskHandler(c *gin.Context) {
	var (
		fileName   = "tasks_handler.go"
		funcName   = "FetchTaskHandler"
		task       []domain.ListTaskResponse
		pagination domain.Pagination
		requestID  string
		errs       error
		err        util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusInternalServerError, delivery_helper.WriteErrorResponse(requestID, "Claims not found"))
		return
	}

	if errs = c.ShouldBindQuery(&pagination); errs != nil {
		c.JSON(http.StatusBadRequest, delivery_helper.WriteErrorResponse(requestID, errs.Error()))
		return
	}

	var (
		userID     int64
		permission string
		limitedID  int64
	)

	myClaims, ok := claims.(jwt.MapClaims)
	if ok {
		userIDTemp := myClaims["user_id"]
		userID = int64(userIDTemp.(float64))
		permissionTemp := myClaims["permission"]
		permission = permissionTemp.(string)
	}

	if permission == "own" {
		limitedID = userID
	}

	task, err = input.taskUsecase.Fetch(c, domain.ContextModel{
		UserLoginID: userID,
		LimitedID:   limitedID,
	}, &domain.TaskRequest{GetListParameter: domain.GetListParameter{
		Page:   int64(pagination.Page),
		Limit:  int64(pagination.Limit),
		Filter: pagination.Filter,
		Order:  pagination.Order,
	}})
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Get Detail!", task))
}

func (input TasksHandler) ViewTaskHandler(c *gin.Context) {
	var (
		fileName  = "tasks_handler.go"
		funcName  = "ViewTaskHandler"
		task      domain.ViewTaskResponse
		requestID string
		errs      error
		err       util.ErrorModel
	)

	defer delivery_helper.WriteLogStdout(errs, &err, requestID, fmt.Sprintf(`[%s, %s]`, fileName, funcName))

	requestIDTemp, _ := c.Get("request_id")
	requestID = requestIDTemp.(string)
	claims, exists := c.Get("claims")
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

	var (
		userID     int64
		permission string
		limitedID  int64
	)

	myClaims, ok := claims.(jwt.MapClaims)
	if ok {
		userIDTemp := myClaims["user_id"]
		userID = int64(userIDTemp.(float64))
		permissionTemp := myClaims["permission"]
		permission = permissionTemp.(string)
	}

	if permission == "own" {
		limitedID = userID
	}

	task, err = input.taskUsecase.GetByID(c, domain.ContextModel{
		UserLoginID: userID,
		LimitedID:   limitedID,
	}, int64(id))
	if err.Err != nil {
		c.JSON(err.Code, delivery_helper.WriteErrorResponse(requestID, err.Err.Error()))
		return
	}

	c.JSON(http.StatusOK, delivery_helper.WriteSuccessResponse(requestID, "Success Get Detail!", task))
}
