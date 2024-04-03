package main

import (
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/config"
	"github.com/cristovaoolegario/tasks-api/internal/controller"
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
	userService := service.NewUserService(userRepo)
	authService := auth.NewAuthService(cfg.AuthSecret, userService)
	loginController := controller.NewLoginController(authService)
	userController := controller.NewUserController(userService)

	router := gin.Default()

	userRoutes := router.Group("/api/users",
		auth.TokenAuthMiddleware(cfg.AuthSecret),
		auth.RoleManagerMiddleware())
	{
		userRoutes.POST("", userController.CreateUser)
		userRoutes.GET("/:username", userController.GetUser)
	}

	router.GET("/login", loginController.LoginHandler)

	health := healthcheck.NewHandler()
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 1*time.Second))
	go http.ListenAndServe(":8086", health)

	router.Run(cfg.AppPort)
}
