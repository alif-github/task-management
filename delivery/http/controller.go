package controller

import (
	"fmt"
	"github.com/alif-github/task-management/app/serverconfig"
	"github.com/alif-github/task-management/config"
	"github.com/alif-github/task-management/delivery/http/handler"
	"github.com/alif-github/task-management/delivery/http/middleware"
	task_repository "github.com/alif-github/task-management/repository/task_repository/pgsql"
	user_repository "github.com/alif-github/task-management/repository/user_repository/pgsql"
	"github.com/alif-github/task-management/usecase/task_usecase"
	"github.com/alif-github/task-management/usecase/user_usecase"
	"github.com/gin-gonic/gin"
	"time"
)

func Controller() {
	prefixPath := config.ApplicationConfiguration.GetServer().PrefixPath
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.BasicMiddleware)

	// For set init service
	userHandler := handler.NewUsersHandler(user_usecase.NewUsersUsecase(user_repository.NewPgsqlUserRepository(serverconfig.ServerAttribute.DBConnection), 10*time.Second))
	taskHandler := handler.NewTasksHandler(task_usecase.NewTaskUsecase(task_repository.NewPgsqlTaskRepository(serverconfig.ServerAttribute.DBConnection), 10*time.Second))

	// WhiteListAPI no authorization here (public)
	whiteListAPI := r.Group(prefixPath + "/oauth")
	whiteListAPI.POST("/login", userHandler.LoginHandler)
	whiteListAPI.POST("/register", userHandler.RegisterHandler)
	whiteListAPI.POST("/logout", userHandler.LogoutHandler)

	// PrivateAPIUser authorization needed (private)
	privateAPIUser := r.Group(prefixPath + "/users")
	privateAPIUser.Use(middleware.AuthMiddleware)
	privateAPIUser.GET("", userHandler.FetchUserHandler)
	privateAPIUser.GET("/:id", userHandler.ViewUserHandler)
	privateAPIUser.PUT("/:id", userHandler.UpdateUserHandler)
	privateAPIUser.DELETE("/:id", userHandler.DeleteUserHandler)

	// PrivateAPITask authorization needed (private)
	privateAPITask := r.Group(prefixPath + "/tasks")
	privateAPITask.Use(middleware.AuthMiddleware)
	privateAPITask.POST("", taskHandler.StoreTaskHandler)
	privateAPITask.GET("", taskHandler.FetchTaskHandler)
	privateAPITask.GET("/:id", taskHandler.ViewTaskHandler)
	privateAPITask.PUT("/:id", taskHandler.UpdateTaskHandler)
	privateAPITask.DELETE("/:id", taskHandler.DeleteTaskHandler)

	_ = r.Run(fmt.Sprintf(`:%s`, config.ApplicationConfiguration.GetServer().Port))
}
