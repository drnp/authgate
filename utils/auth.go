/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file oauth2.go
 * @package utils
 * @author Dr.NP <np@herewe.tech>
 * @since 11/05/2023
 */

package utils

import (
	"authgate/runtime"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	ResponseTypeCode              = "code"
	ResponseTypeToken             = "token"
	ResponseTypePassword          = "password"
	ResponseTypeClientCredentials = "client_credentials"
)

type SessionUser struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	Email       string `json:"email"`
	Account     string `json:"account"`
	MobilePhone string `json:"mobile_phone"`
}

func (su SessionUser) Serialize() []byte {
	b, _ := json.Marshal(su)

	return b
}

type SessionCode struct {
	Code                  string    `json:"code"`
	ClientID              string    `json:"client_id"`
	ClientSecret          string    `json:"client_secret"`
	AccessToken           string    `json:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func (sc SessionCode) Serialize() []byte {
	b, _ := json.Marshal(sc)

	return b
}

func (sc *SessionCode) Unserialize(b []byte) {
	json.Unmarshal(b, sc)
}

type Sign struct {
	Issuer    string
	Sub       string
	Name      string
	Type      string
	ExpiresIn time.Duration
}

type JWT struct {
	Token  string
	Expiry time.Time
}

// Sign JWT token
func JWTSign(sign *Sign) (*JWT, error) {
	now := time.Now()
	exp := now.Add(sign.ExpiresIn)
	if sign.Issuer == "" {
		sign.Issuer = runtime.EnvPrefix + "::" + runtime.AppName
	}

	claims := jwt.MapClaims{
		"issuer": sign.Issuer,
		"sub":    sign.Sub,
		"name":   sign.Name,
		"iat":    now.Unix(),
		"exp":    exp.Unix(),
		"type":   sign.Type,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ts, err := token.SignedString([]byte(runtime.Config.Auth.JWTAccessSecret))
	if err != nil {
		runtime.Logger.Errorf("sign JWT token failed : %s", err)

		return nil, err
	}

	return &JWT{
		Token:  ts,
		Expiry: exp,
	}, nil
}

// Valid JWT token
func JWTValid(ts string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			// Error
			return nil, fmt.Errorf("unexpected signing method: %v", token.Method.Alg())
		}

		return []byte(runtime.Config.Auth.JWTRefreshSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid claims format")
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
