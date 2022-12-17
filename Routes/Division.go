package Routes

import (
	"bangladesh-api/Controllers"

	"github.com/gofiber/fiber/v2"
)

func DivisionDetailsRoute(route fiber.Router) {
	route.Post("/new", Controllers.CreateNewDivision)
	route.Post("/get-all", Controllers.GetAllDivision)
	route.Post("/get-all-populated", Controllers.GetAllDivisionByPopulated)
}
func DivisionNamesRoute(route fiber.Router) {
	route.Post("/names", Controllers.GetAllDivisionNames)
}
