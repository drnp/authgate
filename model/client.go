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

package model

import (
	"authgate/runtime"
	"authgate/utils"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const (
	ClientStatusValid   = 0
	ClientStatusInvalid = 255
)

const (
	AccessKeyLength    = 32
	AccessSecretLength = 40
)

type Client struct {
	bun.BaseModel `bun:"table:clients"`

	ID           string `bun:"id,pk,type:uuid" json:"id"`
	RealmID      string `bun:"realm_id,type:uuid" json:"realm_id"`
	Name         string `bun:"name" json:"name"`
	AccessKey    string `bun:"access_key" json:"access_key"`
	AccessSecret string `bun:"access_secret" json:"access_secret"`
	RedirectURL  string `bun:"redirect_url" json:"redirect_url"`
	Status       int    `bun:"status" json:"status"`

	CreatedAt time.Time    `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time    `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt sql.NullTime `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

func (m *Client) List(ctx context.Context) ([]*Client, error) {
	var clients []*Client
	sq := runtime.DB.NewSelect().Model(&clients)
	if m.RealmID != "" {
		sq = sq.Where("realm_id = ?", m.RealmID)
	}

	err := sq.Scan(ctx, &clients)
	if err != nil {
		runtime.Logger.Errorf("list clients failed : %s", err)
	}

	return clients, err
}

func (m *Client) Get(ctx context.Context) error {
	sq := runtime.DB.NewSelect().Model(m).Limit(1)
	if m.ID != "" {
		sq = sq.Where("id = ?", m.ID)
	}

	if m.RealmID != "" {
		sq = sq.Where("realm_id = ?", m.RealmID)
	}

	if m.Name != "" {
		sq = sq.Where("name = ?", m.Name)
	}

	err := sq.Scan(ctx, m)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("query non-exists client <%s>", m.ID)
		} else {
			runtime.Logger.Errorf("query client failed : %s", err)
		}
	}

	return err
}

func (m *Client) Create(ctx context.Context) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	if m.Status != ClientStatusValid {
		m.Status = ClientStatusInvalid
	}

	if m.AccessKey == "" {
		m.AccessKey = utils.RandomString(AccessKeyLength)
	}

	m.AccessSecret = utils.RandomString(AccessSecretLength)
	iq := runtime.DB.NewInsert().Model(m).Returning("*")
	_, err := iq.Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("insert client failed : %s", err)
	}

	return err
}

func (m *Client) Update(ctx context.Context) error {
	uq := runtime.DB.NewUpdate().Model(m).Where("id = ?", m.ID)
	if m.Name != "" {
		uq = uq.Set("name = ?", m.Name)
	}

	if m.AccessSecret != "" {
		uq = uq.Set("access_secret = ?", m.AccessSecret)
	}

	if m.RedirectURL != "" {
		uq = uq.Set("redirect_url = ?", m.RedirectURL)
	}

	if m.Status != ClientStatusValid {
		m.Status = ClientStatusInvalid
	}

	uq = uq.Set("status = ?", m.Status).Set("updated_at = CURRENT_TIMESTAMP")
	_, err := uq.Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("update client failed : %s", err)
	}

	return err
}

func (m *Client) Delete(ctx context.Context) error {
	dq := runtime.DB.NewDelete().Model(m).Where("id = ?", m.ID)
	_, err := dq.Exec(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("delete non-exists client <%s>", m.ID)

			return nil
		}

		runtime.Logger.Errorf("delete client failed : %s", err)
	}

	return err
}

func (m *Client) Init(ctx context.Context) error {
	_, err := runtime.DB.NewCreateTable().Model(m).IfNotExists().Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("Create table <clients> failed : %s", err)

		return err
	}

	runtime.DB.NewCreateIndex().Model(m).Index("idx_clients_created_at").Column("created_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_clients_updated_at").Column("updated_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_clients_deleted_at").Column("deleted_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_clients_realm_id").Column("realm_id").Exec(ctx)

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
