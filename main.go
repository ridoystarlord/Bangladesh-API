package main

import (
	"bangladesh-api/DBManager"
	"bangladesh-api/Routes"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func SetupRoutes(app *fiber.App) {
	Routes.DivisionDetailsRoute(app.Group("/division"))
	Routes.DivisionNamesRoute(app.Group("/division/only"))
	Routes.ZilaDetailsRoute(app.Group("/zila"))
	Routes.ZilaNamesRoute(app.Group("/zila/only"))
}

func main() {
    fmt.Println(("Hello Bangladesh Api"))
	fmt.Print("Initializing Database Connection ... ")
	initState := DBManager.InitCollections()

	if initState {
		fmt.Println("[OK]")
	} else {
		fmt.Println("[FAILED]")
		return
	}

	fmt.Print("Initializing the server ... ")
	app := fiber.New()
	app.Use(cors.New())
	app.Use(pprof.New())
	SetupRoutes(app)
	app.Static("/Public", "./Public")
	fmt.Println("[OK]")
	app.Listen(":8080")
}