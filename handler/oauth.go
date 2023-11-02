/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file oauth.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package handler

import (
	"authgate/runtime"
	"authgate/service"
	"fmt"

	"github.com/labstack/echo/v4"
)

type OAuth struct {
	svcOAuth2 *service.OAuth2
}

func InitOAuth() *OAuth {
	h := new(OAuth)
	h.svcOAuth2 = service.NewOAuth2Service()

	og := runtime.Server.Group("/oauth")
	og.GET("/login", h.loginPage).Name = "OAuthGetLogin"
	og.GET("/register", h.registerPage).Name = "OAuthGetRegister"
	og.POST("/authorize", h.authorize).Name = "OAuthPostAuthorize"
	og.POST("/token", h.token).Name = "OAuthPostToken"
	og.POST("/revoke", h.revoke).Name = "OAuthPostRevoke"
	og.POST("/introspect", h.introspect).Name = "OAuthPostIntrospect"

	return h
}

func (h *OAuth) loginPage(ctx echo.Context) error {
	return nil
}

func (h *OAuth) registerPage(ctx echo.Context) error {
	return nil
}

func (h *OAuth) authorize(ctx echo.Context) error {
	//err := h.svcOAuth.Authorize(ctx.Request().Context(), ctx.Response().Writer, ctx.Request())
	err := h.svcOAuth2.Authorize(ctx.Request().Context(), ctx.Response().Writer, ctx.Request())
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (h *OAuth) token(ctx echo.Context) error {
	return nil
}

func (h *OAuth) revoke(ctx echo.Context) error {
	return nil
}

func (h *OAuth) introspect(ctx echo.Context) error {
	return nil
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
