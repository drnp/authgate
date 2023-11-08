/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file account.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package handler

import (
	"authgate/handler/request"
	"authgate/handler/response"
	"authgate/model"
	"authgate/service"
	"authgate/utils"
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Account struct {
	svcAccount *service.Account
}

func InitAccount() *Account {
	h := new(Account)
	h.svcAccount = new(service.Account)

	// runtime.Server.GET("/accounts", h.list).Name = "AccountGetList"
	// runtime.Server.GET("/account/:id", h.get).Name = "AccountGet"
	// runtime.Server.POST("/account", h.post).Name = "AccountPost"
	// runtime.Server.PUT("/account/:id", h.put).Name = "AccountPut"
	// runtime.Server.DELETE("/account/:id", h.delete).Name = "AccountDelete"
	// runtime.Server.POST("/auth", h.auth).Name = "AccountAuth"

	return h
}

func (h *Account) list(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	realmID := ctx.QueryParam("realm_id")
	list, err := h.svcAccount.List(ctx.Request().Context(), &service.AccountSvcOptions{
		RealmID: realmID,
	})
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeListAccountFailed
		e.Message = response.MsgListAccountFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	var resp []*response.AccountGet
	for _, info := range list {
		resp = append(resp, &response.AccountGet{
			ID:        info.ID,
			RealmID:   info.RealmID,
			Username:  info.Username,
			Email:     info.Email,
			Mobile:    info.Mobile,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
		})
	}

	e.Data = resp

	return ctx.JSON(http.StatusOK, e)
}

func (h *Account) get(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	info, err := h.svcAccount.Get(ctx.Request().Context(), &service.AccountSvcOptions{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// No matching
			e.Status = http.StatusNotFound
			e.Code = response.CodeTargetNotFound
			e.Message = response.MsgTargetNotFound

			return ctx.JSON(http.StatusNotFound, e)
		}

		e.Status = http.StatusInternalServerError
		e.Code = response.CodeGetAccountFailed
		e.Message = response.MsgGetAccountFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Data = &response.AccountGet{
		ID:        info.ID,
		RealmID:   info.RealmID,
		Username:  info.Username,
		Email:     info.Email,
		Mobile:    info.Mobile,
		Status:    info.Status,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Account) post(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	req := new(request.AccountPost)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	account := &model.Account{
		RealmID:  req.RealmID,
		Username: req.Username,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		Status:   model.AccountStatusValid,
	}
	err = h.svcAccount.Create(ctx.Request().Context(), account)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeCreateAccountFailed
		e.Message = response.MsgCreateAccountFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Status = http.StatusCreated
	e.Data = &response.AccountPost{
		ID:       account.ID,
		RealmID:  account.RealmID,
		Username: account.Username,
		Email:    account.Email,
		Mobile:   account.Mobile,
		Status:   account.Status,
	}

	return ctx.JSON(http.StatusCreated, e)
}

func (h *Account) put(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	req := new(request.AccountPut)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	account := &model.Account{
		ID:       id,
		Username: req.Username,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
		Status:   req.Status,
	}
	err = h.svcAccount.Update(ctx.Request().Context(), account)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeUpdateAccountFailed
		e.Message = response.MsgUpdateAccountFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Account) delete(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	err := h.svcAccount.Delete(ctx.Request().Context(), &service.AccountSvcOptions{
		ID: id,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			e.Status = http.StatusNotFound
			e.Code = response.CodeTargetNotFound
			e.Message = response.MsgTargetNotFound

			return ctx.JSON(http.StatusNotFound, e)
		}

		e.Status = http.StatusInternalServerError
		e.Code = response.CodeDeleteAccountFailed
		e.Message = response.MsgDeleteAccountFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Account) auth(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	req := new(request.AccountAuth)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	checked, err := h.svcAccount.Auth(ctx.Request().Context(), &service.AccountSvcOptions{
		RealmID:  req.RealmID,
		Username: req.Username,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Data = checked

	return ctx.JSON(http.StatusOK, e)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
