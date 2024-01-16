// Package server is the part of the service that serves it.
package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/benni347/novum-portfolio/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
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

// WriteJSONError writes error as JSON to ResponseWriter
func WriteJSONError(w http.ResponseWriter, error string, message string, statusCode int) (int, error) {
	response := ErrorResponse{Error: error, Message: message}
	jsonData, err := json.Marshal(response)
	if err != nil {
		utils.LogErrorDefaultFormat(appName, "WriteJSONError", err, "Marshalling JSON Data")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	n, err := w.Write(jsonData)
	if err != nil {
		utils.LogErrorDefaultFormat(appName, "WriteJSONError", err, "Writing JSON Data")
	}
	return n, err
}

func logRequests(requestersIP string,
	timeStamp time.Time,
	method string,
	httpVersion string,
	returnValue int,
	returnSize int,
	userAgent string,
) {
	// Define your time layout
	layout := "02/Jan/2006:15:04:05 -0700"

	// Use the Format method on timeStamp to get it as a string
	formattedTime := timeStamp.Format(layout)

	utils.Logger.WithFields(log.Fields{
		"requestersIp":            requestersIP,
		"Time":                    formattedTime,
		"HTTP Method":             method,
		"HTTP Version":            httpVersion,
		"HTTP return Status Code": returnValue,
		"Returned Data Size":      returnSize,
		"User Agent":              userAgent,
	}).Info("Log Requests")
}

func health(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
