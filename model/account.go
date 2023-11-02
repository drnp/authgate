/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file account.go
 * @package model
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package model

import (
	"authgate/runtime"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	AccountStatusValid   = 0
	AccountStatusInvalid = 255
)

const (
	SaltLength = 32
)

type Account struct {
	bun.BaseModel `bun:"table:accounts"`

	ID       string `bun:"id,pk,type:uuid" json:"id"`
	RealmID  string `bun:"realm_id,type:uuid" json:"realm_id"`
	Salt     string `bun:"salt" json:"salt"`
	Username string `bun:"username" json:"username"`
	Password string `bun:"password" json:"password"`
	Status   int    `bun:"status" json:"status"`

	// Identities
	Email  string `bun:"email" json:"email"`
	Mobile string `bun:"mobile" json:"mobile"`

	CreatedAt time.Time    `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time    `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt sql.NullTime `bun:"deleted_at,nullzero" json:"-"`
}

func (m *Account) List(ctx context.Context) ([]*Account, error) {
	var accounts []*Account
	sq := runtime.DB.NewSelect().Model(&accounts)
	if m.RealmID != "" {
		sq = sq.Where("realm_id = ?", m.RealmID)
	}

	err := sq.Scan(ctx, &accounts)
	if err != nil {
		runtime.Logger.Errorf("list accounts failed : %s", err)
	}

	return accounts, err
}

func (m *Account) Get(ctx context.Context) error {
	sq := runtime.DB.NewSelect().Model(m).Limit(1)
	if m.ID != "" {
		sq = sq.Where("id = ?", m.ID)
	}

	if m.RealmID != "" {
		sq = sq.Where("realm_id = ?", m.RealmID)
	}

	err := sq.Scan(ctx, m)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("query non-exists account <%s>", m.ID)
		} else {
			runtime.Logger.Errorf("query account failed : %s", err)
		}
	}

	return err
}

func (m *Account) Create(ctx context.Context) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	if m.Status != AccountStatusValid {
		m.Status = AccountStatusInvalid
	}

	iq := runtime.DB.NewInsert().Model(m).Returning("*")
	_, err := iq.Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("insert account failed : %s", err)
	}

	return err
}

func (m *Account) Update(ctx context.Context) error {
	uq := runtime.DB.NewUpdate().Model(m).Where("id = ?", m.ID)
	if m.Username != "" {
		uq = uq.Set("username = ?", m.Username)
	}

	if m.Email != "" {
		uq = uq.Set("email = ?", m.Email)
	}

	if m.Mobile != "" {
		uq = uq.Set("mobile = ?", m.Mobile)
	}

	if m.Salt != "" {
		uq = uq.Set("salt = ?", m.Salt)
	}

	if m.Password != "" {
		uq = uq.Set("password = ?", m.Password)
	}

	if m.Status != AccountStatusValid {
		m.Status = AccountStatusInvalid
	}

	_, err := uq.Set("status = ?", m.Status).Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("update account failed : %s", err)
	}

	return err
}

func (m *Account) Delete(ctx context.Context) error {
	dq := runtime.DB.NewDelete().Model(m).Where("id = ?", m.ID)
	_, err := dq.Exec(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("delete non-exists account <%s>", m.ID)

			return nil
		}

		runtime.Logger.Errorf("delete account failed : %s", err)
	}

	return err
}

func (m *Account) Init(ctx context.Context) error {
	_, err := runtime.DB.NewCreateTable().Model(m).IfNotExists().Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("Create table <accounts> failed : %s", err)

		return err
	}

	runtime.DB.NewCreateIndex().Model(m).Index("idx_accounts_created_at").Column("created_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_accounts_updated_at").Column("updated_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_accounts_deleted_at").Column("deleted_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_accounts_realm_id").Column("realm_id").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Unique().Index("uq_accounts_username").Column("realm_id", "username").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Unique().Index("uq_accounts_email").Column("realm_id", "email").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Unique().Index("uq_accounts_mobile").Column("realm_id", "mobile").Exec(ctx)

	return nil
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
