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
	"authgate/handler/request"
	"authgate/handler/response"
	"authgate/runtime"
	"authgate/service"
	"authgate/utils"
	"encoding/base64"
	"os"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

type Misc struct {
	svcZZAuth *service.ZZAuth
	store     *session.Store
}

func InitMisc() *Misc {
	h := new(Misc)
	h.svcZZAuth = service.NewZZAuth()
	h.store = session.New(session.Config{
		Storage: runtime.Storage,
	})

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

	runtime.Server.Get("/login", h.loginPage).Name("LoginPage")
	runtime.Server.Post("/login", h.login).Name("PostLogin")
	runtime.Server.Get("/logout", h.logout).Name("GetLogout")
	runtime.Server.Get("/register", h.registerPage).Name("RegisterPage")
	runtime.Server.Post("/register", h.register).Name("PostRegister")
	runtime.Server.Get("/confirm", h.confirmPage).Name("ConfirmPage")
	runtime.Server.Post("/confirm", h.confirm).Name("PostConfirm")

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

// @Tags Misc
// @Summary Show login page
// @Description 常规登录页面。如果用户已登录，会显示欢迎页面。
// @ID LoginPage
// @Produce html
// @Success 200 302 {object} nil
// @Router /login [get]
func (h *Misc) loginPage(c *fiber.Ctx) error {
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
	_, ok := sess.Get("user").([]byte)
	if !ok {
		// Not online
		return c.SendFile("./static/login.html")
	}

	// Welcome
	return c.SendFile("./static/welcome.html")
}

// @Tags Misc
// @Summary Process login request
// @Description 处理登录请求，并生成平台session.登录成功后，如果url中参数 r 不为空（base64），将跳转至目标地址
// @ID PostLogin
// @Accept json
// @Produce json
// @Param _ body request.LoginForm true "登录信息"
// @Success 200 {object} nil
// @Success 302 {object} nil
// @Failure 500 {object} utils.Envelope
// @Failure 400 {object} utils.Envelope
// @Failure 401 {object} utils.Envelope
// @Router /login [post]
func (h *Misc) login(c *fiber.Ctx) error {
	e := utils.WrapResponse(nil)
	sess, err := h.store.Get(c)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeStorageFailed
		e.Message = response.MsgStorageFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	req := new(request.LoginForm)
	err = c.BodyParser(req)
	if err != nil || req.Account == "" || req.Password == "" {
		e.Status = fiber.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return c.Status(fiber.StatusBadRequest).Format(e)
	}

	callback := c.Context().Referer()
	r := c.Query("r")
	if r != "" {
		// Redirect back
		callback, _ = base64.StdEncoding.DecodeString(r)
	}

	user, err := h.svcZZAuth.ValidUser(c.Context(), req.Account, req.Password)
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeGetAccountFailed
		e.Message = response.MsgGetAccountFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	if user == nil {
		// Authenticate failed
		e.Status = fiber.StatusUnauthorized
		e.Code = response.CodeAuthFailed
		e.Message = response.MsgAuthFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusUnauthorized).Format(e)
	}

	su := utils.SessionUser{
		ID:          user.ID,
		Name:        user.Name,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Account:     user.Account,
		MobilePhone: user.MobilePhone,
	}
	sess.Set("user", su.Serialize())
	err = sess.Save()
	if err != nil {
		e.Status = fiber.StatusInternalServerError
		e.Code = response.CodeStorageFailed
		e.Message = response.MsgStorageFailed
		e.Data = err.Error()

		return c.Status(fiber.StatusInternalServerError).Format(e)
	}

	return c.Redirect(string(callback))
}

// @Tags Misc
// @Summary Process logout request
// Description 处理登出，成功会跳转回登录页面。
// @ID GetLogout
// @Success 302 {object} nil
// @Failure 500 {object} utils.Envelope
// @Failure 400 {object} utils.Envelope
// @Router /logout [get]
func (h *Misc) logout(c *fiber.Ctx) error {
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
	err = c.Redirect("/login")
	if err != nil {
		e.Status = fiber.StatusBadRequest
		e.Code = response.CodeGeneralHTTPError
		e.Message = response.MsgGeneralHTTPError
		e.Data = err.Error()

		return c.Status(fiber.StatusBadRequest).Format(e)
	}

	return nil
}

// @Tags Misc
// @Summary Show register page
// @Description 常规注册页面，如果用户已登录，会显示欢迎页面。
// @ID RegisterPage
// @Produce html
// @Success 200 302 {object} nil
// @Router /register [get]
func (h *Misc) registerPage(c *fiber.Ctx) error {
	return nil
}

// @Tags Misc
// @Summary Process register request
// @Description 处理注册请求，不自动登录。成功后跳转到登录页面。
// @ID PostRegister
// @Accept json
// Produce json
// @Param _ body request.RegisterForm true "注册信息"
// @Success 302 {object} nil
// @Failure 500 {object} utils.Envelope
// @Failure 400 {object} utils.Envelope
func (h *Misc) register(c *fiber.Ctx) error {
	return nil
}

func (h *Misc) confirmPage(c *fiber.Ctx) error {
	return nil
}

func (h *Misc) confirm(c *fiber.Ctx) error {
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
