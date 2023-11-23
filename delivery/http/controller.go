package controller

import (
	"fmt"
	"github.com/alif-github/clean-arch/app/serverconfig"
	"github.com/alif-github/clean-arch/config"
	"github.com/alif-github/clean-arch/delivery/http"
	"github.com/alif-github/clean-arch/delivery/http/middleware"
	"github.com/alif-github/clean-arch/repository/pgsql"
	"github.com/alif-github/clean-arch/usecase"
	"github.com/gin-gonic/gin"
	"time"
)

func Controller() {
	prefixPath := config.ApplicationConfiguration.GetServer().PrefixPath
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.BasicMiddleware)

	// For set init service
	userHandler := http.NewUsersHandler(usecase.NewUsersUsecase(pgsql.NewPsqlUsersRepository(serverconfig.ServerAttribute.DBConnection), 10*time.Second))
	taskHandler := http.NewTasksHandler(usecase.NewTaskUsecase(pgsql.NewPsqlTasksRepository(serverconfig.ServerAttribute.DBConnection), 10*time.Second))

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
