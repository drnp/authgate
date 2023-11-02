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
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Misc struct {
}

func InitMisc() *Misc {
	h := new(Misc)

	runtime.Server.Any("/", h.index)
	runtime.Server.Any("/docs/*", echoSwagger.WrapHandler)
	runtime.Server.GET("/routers", h.routers)

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
func (h *Misc) index(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, utils.WrapResponse(nil))
}

// routers
func (h *Misc) routers(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, utils.WrapResponse(runtime.Server.Routes()))
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
