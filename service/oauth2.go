/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file oauth2.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 11/01/2023
 */

package service

import (
	"authgate/runtime"
	"context"
	"net/http"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
)

type OAuth2 struct {
	oauth2Manager *manage.Manager
	oauth2Server  *server.Server
}

func NewOAuth2Service() *OAuth2 {
	svc := new(OAuth2)
	clientStore := store.NewClientStore()
	clientStore.Set("my-client", &models.Client{
		ID:     "my-client",
		Secret: "abcdefg",
	})
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	manager.MapAccessGenerate(generates.NewAccessGenerate())
	manager.MapAuthorizeGenerate(generates.NewAuthorizeGenerate())
	manager.MapClientStorage(clientStore)
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	srv := server.NewDefaultServer(manager)
	srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		runtime.Logger.Errorf("internal : %s", err)

		return
	})
	srv.SetResponseErrorHandler(func(re *errors.Response) {
		runtime.Logger.Errorf("response : %s / %s", re.Error, re.Description)
	})
	srv.SetPasswordAuthorizationHandler(func(ctx context.Context, clientID, username, password string) (string, error) {
		return "", nil
	})
	srv.SetUserAuthorizationHandler(func(wr http.ResponseWriter, req *http.Request) (string, error) {
		return "user", nil
	})

	svc.oauth2Manager = manager
	svc.oauth2Server = srv

	return svc
}

func (s *OAuth2) Authorize(ctx context.Context, wr http.ResponseWriter, req *http.Request) error {
	err := s.oauth2Server.HandleAuthorizeRequest(wr, req)

	return err
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
