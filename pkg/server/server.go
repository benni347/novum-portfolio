// Package server is the part of the service that serves it.
package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
)

var serverAddr string

const appName = "novum-portfolio-server"

// ErrorResponse is a struct for error responses
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// NewServer creates the new server and sets up the listener
func NewServer(app *fiber.App) *fiber.App {
	app.Get("/", func(c *fiber.Ctx) error {
		return Render(c, IndexComponent("test"))
	})

	app.Use(healthcheck.New())
	app.Use(NotFoundMiddleware)
	return app
}

func NotFoundMiddleware(c *fiber.Ctx) error {
	return Render(c, NotFound(), templ.WithStatus(http.StatusNotFound))
}

func Render(c *fiber.Ctx, component templ.Component, options ...func(*templ.ComponentHandler)) error {
	componentHandler := templ.Handler(component)
	for _, o := range options {
		o(componentHandler)
	}
	return adaptor.HTTPHandler(componentHandler)(c)
}
