package app

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ladiesman2127/birthdays/api/server/router"
	"github.com/ladiesman2127/birthdays/internal/app/controllers"
	"github.com/ladiesman2127/birthdays/internal/database"
)

func Start() {
	DB_NAME := os.Getenv("DB_NAME")
	app := gin.Default()

	db := database.New(&DB_NAME)

	userController := controllers.NewUsersController(db)

	router.InitRoutes(userController, app)

	if err := app.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
