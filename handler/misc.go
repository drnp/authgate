/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file misc.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package handler

import (
	_ "authgate/docs"
	"authgate/runtime"
	"authgate/utils"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

type Misc struct {
}

func InitMisc() *Misc {
	h := new(Misc)

	// runtime.Server.Any("/", h.index)
	// runtime.Server.Any("/docs/*", echoSwagger.WrapHandler)
	// runtime.Server.GET("/routers", h.routers)
	runtime.Server.All("/", h.index).Name("Index")
	runtime.Server.Get("/routers", h.routers).Name("GetRouters")
	runtime.Server.Get("/swagger.json", h.swaggerJson).Name("GetSwaggerJson")
	runtime.Server.All("/docs/*", swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
	}))

	return h
}

// index

// @Tags Misc
// @Summary Just an empty portal
// @Description 一个不返回任何有效数据的路由，用于展示JSON envelope
// @ID Index
// @Produce json
// @Success 200 {object} nil
// @Router / [get]
func (h *Misc) index(c *fiber.Ctx) error {
	return c.Format(utils.WrapResponse(nil))
}

// func (h *Misc) index(ctx echo.Context) error {
// 	return ctx.JSON(http.StatusOK, utils.WrapResponse(nil))
// }

// routers

// routers

// @Tags Misc
// @Summary Get HTTP routers
// @Description 返回路由列表，由go-fiber自动生成
// @ID GetRouters
// @Produce json
// @Success 200 {object} nil
// @Router /routers [get]
func (h *Misc) routers(c *fiber.Ctx) error {
	return c.Format(utils.WrapResponse(runtime.Server.Stack()))
}

// swaggerJson : Helper for swag docs
func (h *Misc) swaggerJson(c *fiber.Ctx) error {
	content, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		c.WriteString("{}")

		return err
	}

	c.Write(content)

	return nil
}

// func (h *Misc) routers(ctx echo.Context) error {
// 	return ctx.JSON(http.StatusOK, utils.WrapResponse(runtime.Server.Routes()))
// }

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
