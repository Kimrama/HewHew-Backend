package server

import "github.com/gofiber/fiber/v2"

type User struct {
	ID   string `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

var userList = []User{
	{ID: "1", Name: "kimrama"},
	{ID: "2", Name: "Boat"},
}

func (s *fiberServer) initMiscellaneousRoutes() {
	router := s.app.Group("/v1/misc")

	router.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy---v1----test CI/CD",
		})
	})

	router.Get("/users", func(c *fiber.Ctx) error {
		return c.JSON(userList)
	})

	router.Get("/kimrama", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"name": "kimrama",
		})
	})

	router.Post("/users", func(c *fiber.Ctx) error {
		var newUser User
		if err := c.BodyParser(&newUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		userList = append(userList, newUser)
		return c.Status(fiber.StatusCreated).JSON(newUser)
	})

	router.Get("/kimrama/:number", func(c *fiber.Ctx) error {
		number := c.Params("number")
		return c.JSON(fiber.Map{
			"name":   "kimrama",
			"number": number,
		})
	})

	router.Get("/Boat", func(c *fiber.Ctx) error {
		return c.JSON((fiber.Map{
			"name": "Boat",
		}))
	})

}
