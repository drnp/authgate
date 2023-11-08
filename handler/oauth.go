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
	"authgate/handler/request"
	"authgate/handler/response"
	"authgate/runtime"
	"authgate/service"
	"authgate/utils"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type OAuth struct {
	// svcOAuthFosite *service.OAuthFosite
	// svcOAuthOAuth2 *service.OAuthOAuth2
	// svcOAuthRemote *service.OAuthRemote
	svcZZAuth *service.ZZAuth
	store     *session.Store
}

func InitOAuth() *OAuth {
	h := new(OAuth)
	// h.svcOAuthFosite = service.NewOAuthFositeService()
	// h.svcOAuthOAuth2 = service.NewOAuthOAuth2Service()
	// h.svcOAuthRemote = service.NewOAuthRemoteService()
	h.svcZZAuth = service.NewZZAuth()
	h.store = session.New(session.Config{
		Storage: runtime.Storage,
	})

	og := runtime.Server.Group("/oauth")
	og.Get("/login", h.loginPage).Name("OAuthLoginPage")
	og.Post("/login", h.login).Name("OAuthPostLogin")
	og.Get("/logout", h.logout).Name("OAuthGetLogout")
	og.Get("/register", h.registerPage).Name("OAuthRegisterPage")
	og.Post("/register", h.register).Name("OAuthPostRegister")
	og.Get("/confirm", h.confirmPage).Name("OAuthConfirmPage")
	og.Post("/confirm", h.confirm).Name("OAuthPostConfirm")
	og.Get("/authorize", h.authorize).Name("OAuthGetAuthorize")
	og.Post("/token", h.token).Name("OAuthPostToken")
	og.Post("/revoke", h.revoke).Name("OAuthPostRevoke")
	og.Post("/introspect", h.introspect).Name("OAuthPostIntrospect")

	return h
}

func (h *OAuth) loginPage(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) login(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) registerPage(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) logout(c *fiber.Ctx) error {
	e := utils.WrapResponse(nil)
	sess, err := h.store.Get(c)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeStorageFailed
		e.Message = response.MsgStorageFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	err = sess.Destroy()
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeStorageFailed
		e.Message = response.MsgStorageFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	// Redirect
	err = c.RedirectToRoute("OAuthLoginPage", nil)
	if err != nil {
		e.Status = fiber.StatusBadRequest
		e.Code = response.CodeGeneralHTTPError
		e.Message = response.MsgGeneralHTTPError
		e.Data = err.Error()

		return c.Status(fiber.StatusBadRequest).Format(e)
	}

	return nil
}

func (h *OAuth) register(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) confirmPage(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) confirm(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) authorize(c *fiber.Ctx) error {
	e := utils.WrapResponse(nil)
	sess, err := h.store.Get(c)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeStorageFailed
		e.Message = response.MsgStorageFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	// Check login
	user, ok := sess.Get("user").(utils.SessionUser)
	if !ok {
		// Not online
		return c.RedirectToRoute("OAuthLoginPage", nil)
	}

	req := &request.GetAuthorize{}
	err = c.QueryParser(req)
	if err == nil {
		err = req.Validation()
	}

	if err != nil {
		e.Status = fiber.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return c.Status(fiber.StatusBadRequest).Format(e)
	}

	client, err := h.svcZZAuth.ValidClient(c.Context(), req.ClientID)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeGetClientFailed
		e.Message = response.MsgGetClientFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	// Check client visible
	// All pass here

	// Generate code
	sc, err := h.svcZZAuth.GenerateToken(c.Context(), client.ClientID, client.SecretKey, user)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeAuthInternal
		e.Message = response.MsgAuthInternal
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	// Redirect
	u, _ := url.Parse(req.RedirectURI)
	q := u.Query()
	q.Add("code", sc.Code)
	q.Add("state", req.State)
	u.RawQuery = q.Encode()

	return c.Redirect(u.String())
}

func (h *OAuth) token(c *fiber.Ctx) error {
	e := utils.WrapResponse(nil)
	req := new(request.PostToken)
	err := c.BodyParser(req)
	if err != nil ||
		req.Code == "" ||
		req.ClientID == "" ||
		req.ClientSecret == "" {
		e.Status = fiber.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		if err != nil {
			e.Data = err.Error()
		} else {
			e.Data = "params not enough"
		}

		return c.Status(fiber.StatusBadRequest).Format(e)
	}

	sc, err := h.svcZZAuth.GetToken(c.Context(), req.Code)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeAuthInternal
		e.Message = response.MsgAuthInternal
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	if sc == nil {
		// No token here
		e.Status = fiber.StatusNotFound
		e.Code = response.CodeTargetNotFound
		e.Message = response.MsgTargetNotFound
		e.Data = "token not found via given code"

		return c.Status(fiber.StatusNotFound).Format(e)
	}

	if req.ClientID != sc.ClientID || req.ClientSecret != sc.ClientSecret {
		// Check client failed
		e.Status = fiber.StatusForbidden
		e.Code = response.CodeAuthFailed
		e.Message = response.MsgAuthFailed
		e.Data = "client authorize failed"

		return c.Status(fiber.StatusForbidden).Format(e)
	}

	resp := &response.PostToken{
		ClientID:              req.ClientID,
		AccessToken:           sc.AccessToken,
		AccessTokenExpiresAt:  sc.AccessTokenExpiresAt,
		RefreshToken:          sc.RefreshToken,
		RefreshTokenExpiresAt: sc.RefreshTokenExpiresAt,
	}
	e.Data = resp

	return c.Format(e)
}

func (h *OAuth) revoke(c *fiber.Ctx) error {
	return nil
}

func (h *OAuth) introspect(c *fiber.Ctx) error {
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
