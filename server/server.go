package server

import (
	"fmt"
	"hewhew-backend/config"
	"hewhew-backend/database"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type fiberServer struct {
	app  *fiber.App
	db   database.Database
	conf *config.Config
}

var (
	once   sync.Once
	server *fiberServer
)

func NewFiberServer(conf *config.Config, db database.Database) *fiberServer {
	fiberApp := fiber.New()

	once.Do(func() {
		server = &fiberServer{
			app:  fiberApp,
			db:   db,
			conf: conf,
		}
	})
	return server
}
func getCORSMiddleware(allowOrigins []string) fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: strings.Join(allowOrigins, ","),
		AllowHeaders: "Origin, Content-Type, Accept",
	})
}

func getLoggerMiddleware() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${ip} - ${method} ${path} - ${status} - ${latency}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Asia/Bangkok",
	})
}
func (s *fiberServer) Start() {
	CORSMiddleware := getCORSMiddleware(s.conf.Server.AllowedOrigins)
	LoggerMiddleware := getLoggerMiddleware()
	s.app.Use(CORSMiddleware)
	s.app.Use(LoggerMiddleware)

	// Initialize routes
	s.initMiscellaneousRoutes()
	s.initUserRouter()

	s.Listen()

}

func (s *fiberServer) Listen() error {
	serverURL := fmt.Sprintf(":%d", s.conf.Server.Port)
	return s.app.Listen(serverURL)
}
