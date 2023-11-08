/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file oauth.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 11/05/2023
 */

package request

import (
	"authgate/utils"
	"errors"
	"net/url"
)

var (
	ErrInvalidRequest      = errors.New("invalid request")
	ErrInvalidRedirectURI  = errors.New("invalid redirect_uri")
	ErrInvalidResponseType = errors.New("invalid response type")
)

type GetAuthorize struct {
	ClientID     string `query:"client_id"`
	RedirectURI  string `query:"redirect_uri"`
	ResponseType string `query:"response_type"`
	Scope        string `query:"scope"`
	State        string `query:"state"`
	Nonce        string `query:"nonce"`
}

func (r *GetAuthorize) Validation() error {
	if r.ClientID == "" || r.Scope == "" || r.State == "" {
		return ErrInvalidRequest
	}

	u, err := url.ParseRequestURI(r.RedirectURI)
	if err != nil {
		return ErrInvalidRedirectURI
	}

	r.RedirectURI = u.String()
	if r.ResponseType != utils.ResponseTypeCode &&
		r.ResponseType != utils.ResponseTypeToken &&
		r.ResponseType != utils.ResponseTypePassword &&
		r.ResponseType != utils.ResponseTypeClientCredentials {
		return ErrInvalidResponseType
	}

	return nil
}

type PostToken struct {
	GrantType    string `json:"grant_type" form:"grant_type"`
	Code         string `json:"code" form:"code"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
