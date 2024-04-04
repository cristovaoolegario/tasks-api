package main

import (
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/config"
	"github.com/cristovaoolegario/tasks-api/internal/controller"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	mysql "github.com/cristovaoolegario/tasks-api/internal/infra/db/mysql"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	dbConnect := mysql.InitDB(cfg.DbConnection)
	db, err := dbConnect.DB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	userRepo := mysql.NewUserRepository(dbConnect)
	taskRepo := mysql.NewTaskRepository(dbConnect)

	userService := service.NewUserService(userRepo)
	taskService := service.NewTaskService(taskRepo)
	authService := auth.NewAuthService(cfg.AuthSecret, userService)

	loginController := controller.NewLoginController(authService)
	userController := controller.NewUserController(userService)
	taskController := controller.NewTaskController(taskService, authService)

	CreateDefaultUserIfNotExists(userRepo)

	router := gin.Default()

	router.GET("/login", loginController.LoginHandler)

	userRoutes := router.Group("/api/users",
		auth.TokenAuthMiddleware(cfg.AuthSecret),
		auth.RoleManagerMiddleware())
	{
		userRoutes.POST("", userController.CreateUser)
		userRoutes.GET("/:username", userController.GetUser)
	}

	taskRoutes := router.Group("/api/tasks",
		auth.TokenAuthMiddleware(cfg.AuthSecret))
	{
		taskRoutes.GET("", taskController.FindByUserID)
		taskRoutes.POST("", taskController.CreateTaskHandler)
		taskRoutes.PUT("/:id", taskController.UpdateTaskHandler)

		taskRoutes.GET("/:id", auth.RoleManagerMiddleware(), taskController.FindByID)
		taskRoutes.DELETE("/:id", auth.RoleManagerMiddleware(), taskController.DeleteTaskHandler)
	}

	health := healthcheck.NewHandler()
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 1*time.Second))
	go http.ListenAndServe(":8086", health)

	router.Run(cfg.AppPort)
}

func CreateDefaultUserIfNotExists(userRepo *mysql.UserRepository) {
	user, err := userRepo.FindByUsername("newuser")
	if user == nil {
		newUser := &model.User{Username: "newuser", Role: "manager", Password: "securepassword"}
		err = userRepo.Create(newUser)
		if err != nil {
			panic("failed to create user")
		}
	}
}
