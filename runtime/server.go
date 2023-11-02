/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file server.go
 * @package runtime
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package runtime

import (
	"os"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

var Server *echo.Echo

func InitServer() error {
	// app := fiber.New(fiber.Config{
	// 	ServerHeader:          AppName,
	// 	DisableKeepalive:      false,
	// 	AppName:               AppName,
	// 	Prefork:               Config.HTTP.Prefork,
	// 	DisableStartupMessage: true,
	// })
	// app.Use(fiberzap.New(fiberzap.Config{
	// 	Logger: LoggerRaw,
	// }))
	// app.Use(recover.New())
	// app.Use(cors.New())
	// app.Use(requestid.New())

	// Server = app
	app := echo.New()
	app.HideBanner = true
	app.Use(echozap.ZapLogger(LoggerRaw))
	app.Use(middleware.Recover())
	app.Use(middleware.CORS())
	app.Use(middleware.RequestID())

	Server = app

	return nil
}

func Serve() error {
	if Server == nil {
		// Not initialized
		return errors.New("Server not initialized")
	}

	Logger.Infof("starting HTTP server on [%s]", Config.HTTP.ListenAddr)

	//return Server.Listen(Config.HTTP.ListenAddr)
	return Server.Start(Config.HTTP.ListenAddr)
}

func Exit() {
	// TODO: Pure runtime
	os.Exit(-1)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
