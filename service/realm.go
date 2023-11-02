/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file realm.go
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

type Realm struct {
}

type RealmSvcOptions struct {
	ID   string
	Name string
}

func (s *Realm) List(ctx context.Context, opt *RealmSvcOptions) ([]*model.Realm, error) {
	m := &model.Realm{}

	return m.List(ctx)
}

func (s *Realm) Get(ctx context.Context, opt *RealmSvcOptions) (*model.Realm, error) {
	m := &model.Realm{
		ID:   opt.ID,
		Name: opt.Name,
	}

	err := m.Get(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Realm) Create(ctx context.Context, realm *model.Realm) error {
	if realm == nil {
		return errors.New("null realm instance")
	}

	return realm.Create(ctx)
}

func (s *Realm) Update(ctx context.Context, realm *model.Realm) error {
	if realm == nil {
		return errors.New("null realm instance")
	}

	return realm.Update(ctx)
}

func (s *Realm) Delete(ctx context.Context, opt *RealmSvcOptions) error {
	m := &model.Realm{
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
