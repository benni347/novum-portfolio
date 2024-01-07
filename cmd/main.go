package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	v1 "github.com/benni347/novum-portfolio/pkg/server"
	"github.com/benni347/novum-portfolio/pkg/utils"
	"github.com/labstack/echo/v4"

	_ "net/http/pprof"

	// "github.com/pkg/profile"
	log "github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

const appName = "novum-portfolio"

var version = "develop"

func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Usage = "My new and improved portfolio"
	app.Version = version
	app.Authors = []*cli.Author{
		{
			Name:  "Cédric Skwar",
			Email: "cdrc@5y5.one",
		},
	}
	// Define flags
	flags := []cli.Flag{
		&cli.StringFlag{Name: "config"},
		altsrc.NewIntFlag(&cli.IntFlag{
			Name:    "port",
			Value:   8080,
			Usage:   "Port number",
			EnvVars: []string{"EXPOSE_PORT"},
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "log.level",
			Value: log.InfoLevel.String(),
			Usage: "Log level",
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:  "address",
			Value: "0.0.0.0",
			Usage: "address",
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:  "profile",
			Value: false,
			Usage: "Activates Profiling support via 'net/http/pprof'",
		}),
	}
	app.Suggest = true
	app.Compiled = time.Now()

	app.Before = func(c *cli.Context) error {
		err := altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))(c)
		if err != nil {
			return err
		}
		return nil
	}

	// Define action to be executed when the app is run
	app.Action = func(c *cli.Context) error {
		initApp(c)
		runApp(c)
		return nil
	}

	app.Flags = flags
	// Run the app
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func initApp(c *cli.Context) {
	// logging
	initLogging(c)
}

// initLogging
func initLogging(c *cli.Context) {
	logLevel, _ := log.ParseLevel(c.String("log.level"))
	utils.NewLogger(logLevel)
	utils.Logger.WithField("log-level", logLevel).Debug("logger setup")
}

func runApp(c *cli.Context) {
	utils.Logger.WithFields(log.Fields{
		"Port":      c.Int("port"),
		"Address":   c.String("address"),
		"Profiling": c.Bool("profile"),
	}).Debug("Programm Arguments")

	utils.Logger.WithFields(log.Fields{
		"appStart": appName,
		"listenOn": fmt.Sprintf("%s:%d", c.String("address"), c.Int("port")),
	}).Debug("started")
	ln, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", c.String("address"), c.Int("port")))
	if err != nil {
		utils.LogFatalDefaultFormat(appName, "runApp", err, "net.listen")
	}
	defer func() {
		if closeErr := ln.Close(); closeErr != nil {
			utils.LogErrorDefaultFormat(appName, "runApp", closeErr, "ln.Close()")
		}
	}()

	router := echo.New()
	if c.Bool("profile") {
		go func() {
			utils.Logger.Debug(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	router = v1.NewServer(router)
	// defer profile.Start().Stop()
	utils.LogFatalDefaultFormat(appName, "runApp", http.Serve(ln, router), "http.Serve")
}
