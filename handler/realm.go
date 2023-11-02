/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file realm.go
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

type Realm struct {
	svcRealm *service.Realm
}

func InitRealm() *Realm {
	h := new(Realm)
	h.svcRealm = new(service.Realm)

	runtime.Server.GET("/realms", h.list).Name = "RealmGetList"
	runtime.Server.GET("/realm/:id", h.get).Name = "RealmGet"
	runtime.Server.POST("realm", h.post).Name = "RealmPost"
	runtime.Server.PUT("/realm/:id", h.put).Name = "RealmPut"
	runtime.Server.DELETE("/realm/:id", h.delete).Name = "RealmDelete"

	return h
}

func (h *Realm) list(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	list, err := h.svcRealm.List(ctx.Request().Context(), nil)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeListRealmFailed
		e.Message = response.MsgListRealmFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	var resp []*response.RealmGet
	for _, info := range list {
		resp = append(resp, &response.RealmGet{
			ID:        info.ID,
			Name:      info.Name,
			Status:    info.Status,
			CreatedAt: info.CreatedAt,
			UpdatedAt: info.UpdatedAt,
		})
	}

	e.Data = resp

	return ctx.JSON(http.StatusOK, e)
}

func (h *Realm) get(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	info, err := h.svcRealm.Get(ctx.Request().Context(), &service.RealmSvcOptions{
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
		e.Code = response.CodeGetRealmFailed
		e.Message = response.MsgGetRealmFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Data = &response.RealmGet{
		ID:        info.ID,
		Name:      info.Name,
		Status:    info.Status,
		CreatedAt: info.CreatedAt,
		UpdatedAt: info.UpdatedAt,
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Realm) post(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	req := new(request.RealmPost)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	realm := &model.Realm{
		Name:   req.Name,
		Status: model.RealmStatusValid,
	}
	err = h.svcRealm.Create(ctx.Request().Context(), realm)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeCreateRealmFailed
		e.Message = response.MsgCreateRealmFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	e.Status = http.StatusCreated
	e.Data = &response.RealmPost{
		ID:     realm.ID,
		Name:   realm.Name,
		Status: realm.Status,
	}

	return ctx.JSON(http.StatusCreated, e)
}

func (h *Realm) put(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	req := new(request.RealmPut)
	err := ctx.Bind(req)
	if err != nil {
		e.Status = http.StatusBadRequest
		e.Code = response.CodeInvalidParameter
		e.Message = response.MsgInvalidParameter
		e.Data = err.Error()

		return ctx.JSON(http.StatusBadRequest, e)
	}

	realm := &model.Realm{
		ID:     id,
		Name:   req.Name,
		Status: req.Status,
	}
	err = h.svcRealm.Update(ctx.Request().Context(), realm)
	if err != nil {
		e.Status = http.StatusInternalServerError
		e.Code = response.CodeUpdateRealmFailed
		e.Message = response.MsgUpdateRealmFailed
		e.Data = err.Error()

		return ctx.JSON(http.StatusInternalServerError, e)
	}

	return ctx.JSON(http.StatusOK, e)
}

func (h *Realm) delete(ctx echo.Context) error {
	e := utils.WrapResponse(nil)
	id := ctx.Param("id")
	err := h.svcRealm.Delete(ctx.Request().Context(), &service.RealmSvcOptions{
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
		e.Code = response.CodeDeleteRealmFailed
		e.Message = response.MsgDeleteRealmFailed
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
