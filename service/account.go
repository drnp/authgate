/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file account.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 10/27/2023
 */

package service

import (
	"authgate/model"
	"authgate/utils"
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
}

type AccountSvcOptions struct {
	ID       string
	RealmID  string
	Username string
	Email    string
	Mobile   string
	Password string
}

func (s *Account) List(ctx context.Context, opt *AccountSvcOptions) ([]*model.Account, error) {
	m := &model.Account{
		RealmID: opt.RealmID,
	}

	return m.List(ctx)
}

func (s *Account) Get(ctx context.Context, opt *AccountSvcOptions) (*model.Account, error) {
	m := &model.Account{
		ID:      opt.ID,
		RealmID: opt.RealmID,
	}

	err := m.Get(ctx)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *Account) Create(ctx context.Context, account *model.Account) error {
	if account == nil {
		return errors.New("null account instance")
	}

	if account.RealmID == "" {
		return errors.New("empty realm_id")
	}

	if account.Password == "" {
		return errors.New("no password provided")
	}

	salt := utils.RandomString(model.SaltLength)
	pwd := account.Password + salt
	salted, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	account.Salt = salt
	account.Password = string(salted)

	return account.Create(ctx)
}

func (s *Account) Update(ctx context.Context, account *model.Account) error {
	if account == nil {
		return errors.New("null account instance")
	}

	if account.Password != "" {
		salt := utils.RandomString(model.SaltLength)
		pwd := account.Password + salt
		salted, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		account.Salt = salt
		account.Password = string(salted)
	}

	return account.Update(ctx)
}

func (s *Account) Delete(ctx context.Context, opt *AccountSvcOptions) error {
	m := &model.Account{
		ID: opt.ID,
	}

	return m.Delete(ctx)
}

func (s *Account) Auth(ctx context.Context, opt *AccountSvcOptions) (bool, error) {
	if opt.Password == "" {
		return false, errors.New("no password ")
	}

	if opt.RealmID == "" {
		return false, errors.New("empty realm_id")
	}

	m := &model.Account{
		RealmID:  opt.RealmID,
		Username: opt.Username,
		Email:    opt.Email,
		Mobile:   opt.Mobile,
		Password: opt.Password,
	}
	err := m.Get(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		// Account does not exists
		return false, errors.New("Account does not exists")
	}

	if err != nil {
		return false, err
	}

	pwd := opt.Password + m.Salt
	checked := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(pwd))
	if checked == nil {
		// Status OK
		return true, nil
	}

	return false, nil
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
