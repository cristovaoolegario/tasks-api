// @title Tasks API
// @version 1.0
// @description This is a sample api for managing technicians tasks

// @contact.name API Support
// @contact.email cristovaoolegario@gmail.com

// @securityDefinitions.apiKey token
// @in header
// @name Authorization
// @Security JWT
package main

import (
	"database/sql"
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/config"
	"github.com/cristovaoolegario/tasks-api/internal/controller"
	"github.com/cristovaoolegario/tasks-api/internal/domain/model"
	"github.com/cristovaoolegario/tasks-api/internal/domain/service"
	"github.com/cristovaoolegario/tasks-api/internal/infra/db/mysql"
	"github.com/cristovaoolegario/tasks-api/internal/infra/kafka"
	"github.com/gin-gonic/gin"
	"github.com/heptiolabs/healthcheck"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"time"

	_ "github.com/cristovaoolegario/tasks-api/docs"
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
	producer := kafka.NewProducerServiceImp(cfg.BrokerHost)
	managerNotificationService := service.NewManagerNotificationService(cfg.ManagerNotificationTopic, producer)

	loginController := controller.NewLoginController(authService)
	userController := controller.NewUserController(userService)
	taskController := controller.NewTaskController(taskService, authService, managerNotificationService)

	CreateDefaultUserIfNotExists(userRepo)

	router := setupRoutes(loginController, cfg, userController, taskController, db)

	router.Run(cfg.AppPort)
}

func setupRoutes(loginController *controller.LoginController,
	cfg *config.Config,
	userController *controller.UserController,
	taskController *controller.TaskController,
	db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.POST("/login", loginController.LoginHandler)

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
	go http.ListenAndServe(cfg.HealthCheckPort, health)

	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return router
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
