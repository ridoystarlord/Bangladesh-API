package Routes

import (
	"bangladesh-api/Controllers"

	"github.com/gofiber/fiber/v2"
)

func ZilaDetailsRoute(route fiber.Router) {
	route.Post("/new", Controllers.CreateNewZila)
	route.Post("/get-all", Controllers.GetAllZila)
}
func ZilaNamesRoute(route fiber.Router) {
	route.Post("/names", Controllers.GetAllZilaNames)
	route.Post("/names/:divisionName", Controllers.GetAllZilaNamesByDivision)
}
