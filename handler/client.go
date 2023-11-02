/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package handler

import (
	"authgate/handler/request"
	"authgate/handler/response"
	"authgate/model"
	"authgate/runtime"
	"authgate/service"
	"authgate/utils"
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Client struct {
	svcClient *service.Client
}

func InitClient() *Client {
	h := new(Client)
	h.svcClient = new(service.Client)

	runtime.Server.GET("/clients", h.list).Name = "ClientGetList"
	runtime.Server.GET("/client/:id", h.get).Name = "ClientGet"
	runtime.Server.POST("/client", h.post).Name = "ClientPost"
	runtime.Server.PUT("/client/:id", h.put).Name = "ClientPut"
	runtime.Server.DELETE("/client/:id", h.delete).Name = "ClientDelete"

	return h
}

func (h *Client) list(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	realmID := ctx.QueryParam("realm_id")
	list, err := h.svcClient.List(ctx.Request().Context(), &service.ClientSvcOptions{
		RealmID: realmID,
	})
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeListClientFailed
		e.Message = response.MsgListClientFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	var resp []*response.ClientGet
	for _, info := range list {
		resp = append(resp, &response.ClientGet{
			ID:        info.ID,
			RealmID:   info.RealmID,
			Name:      info.Name,
			AccessKey: info.AccessKey,
			Status:    info.Status,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
		})
	}

	e.Data = resp

	return ctx.JSON(http.StatusOK, e)
}

func (h *Client) get(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	info, err := h.svcClient.Get(ctx.Request().Context(), &service.ClientSvcOptions{
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
		e.Code = response.CodeGetClientFailed
		e.Message = response.MsgGetClientFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Data = &response.ClientGet{
		ID:        info.ID,
		RealmID:   info.RealmID,
		Name:      info.Name,
		AccessKey: info.AccessKey,
		Status:    info.Status,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Client) post(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	req := new(request.ClientPost)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	client := &model.Client{
		RealmID: req.RealmID,
		Name:    req.Name,
		Status:  model.ClientStatusValid,
	}
	err = h.svcClient.Create(ctx.Request().Context(), client)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeCreateClientFailed
		e.Message = response.MsgCreateClientFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Status = http.StatusCreated
	e.Data = &response.ClientPost{
		ID:           client.ID,
		RealmID:      client.RealmID,
		Name:         client.Name,
		AccessKey:    client.AccessKey,
		AccessSecret: client.AccessSecret,
		Status:       client.Status,
	}

	return ctx.JSON(http.StatusCreated, e)
}

func (h *Client) put(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	req := new(request.ClientPut)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	client := &model.Client{
		ID:           id,
		Name:         req.Name,
		AccessSecret: req.AccessSecret,
		Status:       req.Status,
	}
	err = h.svcClient.Update(ctx.Request().Context(), client)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeUpdateClientFailed
		e.Message = response.MsgUpdateClientFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Client) delete(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	err := h.svcClient.Delete(ctx.Request().Context(), &service.ClientSvcOptions{
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
		e.Code = response.CodeDeleteClientFailed
		e.Message = response.MsgDeleteClientFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

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
