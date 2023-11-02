/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 10/27/2023
 */

package service

import (
	"authgate/model"
	"context"
	"errors"
)

type Client struct {
}

type ClientSvcOptions struct {
	ID      string
	RealmID string
	Name    string
}

func (s *Client) List(ctx context.Context, opt *ClientSvcOptions) ([]*model.Client, error) {
	m := &model.Client{
		RealmID: opt.RealmID,
	}

	return m.List(ctx)
}

func (s *Client) Get(ctx context.Context, opt *ClientSvcOptions) (*model.Client, error) {
	m := &model.Client{
		ID:      opt.ID,
		RealmID: opt.RealmID,
		Name:    opt.Name,
	}

	err := m.Get(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Client) Create(ctx context.Context, client *model.Client) error {
	if client == nil {
		return errors.New("null client instance")
	}

	if client.RealmID == "" {
		return errors.New("empty realm_id")
	}

	return client.Create(ctx)
}

func (s *Client) Update(ctx context.Context, client *model.Client) error {
	if client == nil {
		return errors.New("null client instance")
	}

	return client.Update(ctx)
}

func (s *Client) Delete(ctx context.Context, opt *ClientSvcOptions) error {
	m := &model.Client{
		ID: opt.ID,
	}

	return m.Delete(ctx)
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
