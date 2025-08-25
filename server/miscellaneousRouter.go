package server

import "github.com/gofiber/fiber/v2"

func (s *fiberServer) initMiscellaneousRoutes() {
	router := s.app.Group("/v1/misc")

	router.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy---v1----test CI/CD",
		})
	})
}
