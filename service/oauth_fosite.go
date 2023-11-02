/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file oauth_forsite.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 10/31/2023
 */

package service

import (
	"authgate/runtime"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/openid"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/jwt"
)

type OAuthFosite struct {
	oauth2Provider fosite.OAuth2Provider
}

func NewOAuthFositeService() *OAuthFosite {
	svc := new(OAuthFosite)
	key, _ := rsa.GenerateKey(rand.Reader, 2048)
	svc.oauth2Provider = compose.ComposeAllEnabled(
		&fosite.Config{
			AccessTokenLifespan: time.Minute * time.Duration(runtime.Config.Auth.JWTAccessExpiry),
		},
		storage.NewExampleStore(),
		key,
	)

	return svc
}

func (s *OAuthFosite) Authorize(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	ar, err := s.oauth2Provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", ar.GetGrantedScopes())
	req.ParseForm()
	for _, scope := range req.PostForm["scopes"] {
		fmt.Println(scope, " granted")
		ar.GrantScope(scope)
	}
	token := newToken()
	resp, err := s.oauth2Provider.NewAuthorizeResponse(ctx, ar, token)
	if err != nil {
		//s.oauth2Provider.WriteAuthorizeError(ctx, writer, ar, err)

		return err
	}

	s.oauth2Provider.WriteAuthorizeResponse(ctx, writer, ar, resp)

	return nil
}

func newToken() *openid.DefaultSession {
	return &openid.DefaultSession{
		Claims: &jwt.IDTokenClaims{},
		Headers: &jwt.Headers{
			Extra: make(map[string]interface{}),
		},
	}
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
