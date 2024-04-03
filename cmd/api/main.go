package main

import (
	"fmt"
	"github.com/cristovaoolegario/tasks-api/internal/auth"
	"github.com/cristovaoolegario/tasks-api/internal/config"
	"github.com/cristovaoolegario/tasks-api/internal/controller"
	mysql "github.com/cristovaoolegario/tasks-api/internal/infra/db/mysql"
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

	repo := mysql.NewUserRepository(dbConnect)
	authService := auth.NewAuthService(cfg.AuthSecret, repo)
	loginController := controller.NewLoginController(repo, authService)

	//newUser := &model.User{Username: "newuser", Role: "manager", Password: "securepassword"}
	//err = repo.Create(newUser)
	//if err != nil {
	//	panic("failed to create user")
	//}

	health := healthcheck.NewHandler()
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db, 1*time.Second))

	fmt.Println("Ready! Listing on port", cfg.AppPort)

	http.HandleFunc("/login", loginController.LoginHandler)
	go http.ListenAndServe(":8086", health)

	http.ListenAndServe(cfg.AppPort, nil)
}
