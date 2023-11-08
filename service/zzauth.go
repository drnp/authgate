/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file zzauth.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 11/05/2023
 */

package service

import (
	"authgate/runtime"
	"authgate/utils"
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

const (
	ZZAuthSecret      = "yaZpfeS*Vgahl6J4wR0Vrli%gt@8h!xrDUQ&nxCfzr!FFttX$P4c3Lk9fgKK3RJ&"
	ZZAuthSalt        = "pBA2@P5dfq0#OXS63kdmVyuedfZGjYCu"
	ZZClientValidPath = "/api/v1/oauth/client/valid"
	ZZUserValidPath   = "/api/v1/oauth/user/valid"

	AccessCodeLength = 40
)

type ZZClientValidRequest struct {
	ClientID  string `json:"client_id"`
	Timestamp int    `json:"timestamp"`
	Expire    int    `json:"expire"`
	Token     string `json:"token"`
}

type ZZClient struct {
	ID          int    `json:"id"`
	ClientID    string `json:"client_id"`
	ClientName  string `json:"client_name"`
	ClientLogo  string `json:"client_logo"`
	ClientDesc  string `json:"client_desc"`
	SecretKey   string `json:"secret_key"`
	RedirectURL string `json:"redirect_url"`
	TenantName  string `json:"tenant_name"`
}

type ZZClientValidResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *ZZClient `json:"data"`
}

type ZZUserValidRequest struct {
	Account   string `json:"account"`
	Password  string `json:"password"`
	Timestamp int    `json:"timestamp"`
	Expire    int    `json:"expire"`
	Token     string `json:"token"`
}

type ZZUser struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Account     string `json:"account"`
	MobilePhone string `json:"mobile_phone"`
}

type ZZUserValidResponse struct {
	Code int     `json:"code"`
	Msg  string  `json:"msg"`
	Data *ZZUser `json:"data"`
}

type ZZAuth struct{}

func NewZZAuth() *ZZAuth {
	svc := new(ZZAuth)

	return svc
}

func (s *ZZAuth) ValidClient(ctx context.Context, clientID string) (*ZZClient, error) {
	now := time.Now()
	req := &ZZClientValidRequest{
		ClientID:  clientID,
		Timestamp: int(now.Unix()),
		Expire:    int(now.Unix()) + 120,
	}
	resp := new(ZZClientValidResponse)
	req.Token = utils.MD5String(
		strconv.Itoa(req.Timestamp) +
			ZZAuthSecret +
			req.ClientID +
			ZZAuthSalt +
			strconv.Itoa(req.Expire))
	c := fiber.Post(runtime.Config.ZZAuth.BaseURL + ZZClientValidPath).JSON(req)
	status, body, _ := c.Bytes()
	if status != fiber.StatusOK {
		return nil, errors.New("fetch client failed")
	}

	err := json.Unmarshal([]byte(body), resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (s *ZZAuth) ValidUser(ctx context.Context, account, password string) (*ZZUser, error) {
	now := time.Now()
	req := &ZZUserValidRequest{
		Account:   account,
		Password:  password,
		Timestamp: int(now.Unix()),
		Expire:    int(now.Unix()) + 120,
	}
	resp := new(ZZUserValidResponse)
	req.Token = utils.MD5String(
		strconv.Itoa(req.Timestamp) +
			ZZAuthSecret +
			req.Account +
			ZZAuthSalt +
			req.Password +
			strconv.Itoa(req.Expire))
	c := fiber.Post(runtime.Config.ZZAuth.BaseURL + ZZUserValidPath).JSON(req)
	status, body, _ := c.Bytes()
	if status != fiber.StatusOK {
		return nil, errors.New("fetch user failed")
	}

	err := json.Unmarshal([]byte(body), resp)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (s *ZZAuth) GenerateToken(ctx context.Context, clientID, secretKey string, user *utils.SessionUser) (*utils.SessionCode, error) {
	jwtAccess, err := utils.JWTSign(&utils.Sign{
		Sub:       strconv.Itoa(user.ID),
		Name:      user.Account,
		Type:      "access",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTAccessExpiry) * time.Second,
		Key:       []byte(secretKey),
	})
	if err != nil {
		return nil, err
	}

	jwtRefresh, err := utils.JWTSign(&utils.Sign{
		Sub:       strconv.Itoa(user.ID),
		Name:      user.Account,
		Type:      "refresh",
		ExpiresIn: time.Duration(runtime.Config.Auth.JWTRefreshExpiry) * time.Second,
		Key:       []byte(secretKey),
	})
	if err != nil {
		return nil, err
	}

	code := utils.RandomString(AccessCodeLength)
	sc := utils.SessionCode{
		Code:                  code,
		ClientID:              clientID,
		ClientSecret:          secretKey,
		AccessToken:           jwtAccess.Token,
		AccessTokenExpiresAt:  jwtAccess.Expiry,
		RefreshToken:          jwtRefresh.Token,
		RefreshTokenExpiresAt: jwtRefresh.Expiry,
	}

	err = runtime.Storage.Set(code, sc.Serialize(), time.Duration(runtime.Config.Auth.AuthorizeCodeExpiry)*time.Second)
	if err != nil {
		return nil, err
	}

	return &sc, nil
}

func (s *ZZAuth) GetToken(ctx context.Context, code string) (*utils.SessionCode, error) {
	b, err := runtime.Storage.Get(code)
	if err != nil {
		return nil, err
	}

	if b == nil {
		return nil, nil
	}

	sc := new(utils.SessionCode)
	sc.Unserialize(b)

	return sc, runtime.Storage.Delete(code)
}

func (s *ZZAuth) RefreshToken(ctx context.Context, refreshToken, clientSecret string) (*utils.SessionCode, error) {
	claims, err := utils.JWTValid(refreshToken, clientSecret)
	if err != nil {
		return nil, err
	}

	sign := new(utils.Sign)
	sign.Sub, _ = claims["sub"].(string)
	sign.Name, _ = claims["name"].(string)
	sign.Type = "access"
	sign.ExpiresIn = time.Duration(runtime.Config.Auth.JWTAccessExpiry) * time.Second
	sign.Key = []byte(clientSecret)
	jwtAccess, err := utils.JWTSign(sign)
	if err != nil {
		return nil, err
	}

	sc := &utils.SessionCode{
		AccessToken:          jwtAccess.Token,
		AccessTokenExpiresAt: jwtAccess.Expiry,
	}

	return sc, nil
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
