/*
OUR BEAUTIFUL PROJET WAS LEAD By :

- FASSINOU Noudehouene Régis Jaurès

ACCOMPANIED BY BRILLIANT DEVELOPERS:

- Kuba STANISLAWSKI

- Mohamed Makhtar MBEMGUE

- Arafat FEICAL IDRISSA

- Max Corval KAKPODJO AISSI
*/

package main

import (
	"fmt"

	"estiam_golang_api_course_finalproject/config"
	"estiam_golang_api_course_finalproject/handlers"
	"estiam_golang_api_course_finalproject/repos"
	"estiam_golang_api_course_finalproject/services"

	"github.com/labstack/echo/v4"
)

func main() {
	server := echo.New()

	// load config
	config := config.Load()
	userRepo := repos.NewUserRepository(config.DbConn)
	userService := services.NewUserService(userRepo, 14)

	healthHandler := handlers.NewHealthHandler()
	server.GET("/live", healthHandler.IsAlive)

	// REMOVE THAT ENDPOINT
	// userHandler := handlers.NewUserHandler(userService)
	// server.GET("/users/:id", userHandler.Get)

	// TODO: Register a new endpoint for POST user

	// Register the endpoints
	userHandler := handlers.NewUserHandler(userService)
	server.POST("/users", userHandler.Post)  // Endpoint for creating a new user
	server.POST("/login", userHandler.Login) // Endpoint for user login and JWT generation
	server.GET("/users", userHandler.GetAllUsers)

	if err := server.Start(":8080"); err != nil {
		fmt.Println(err)
	}
}
