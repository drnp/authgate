/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file realm.go
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
	RealmStatusValid   = 0
	RealmStatusInvalid = 255
)

type Realm struct {
	bun.BaseModel `bun:"table:realms"`

	ID     string `bun:"id,pk,type:uuid" json:"id"`
	Name   string `bun:"name" json:"name"`
	Status int    `bun:"status" json:"status"`

	CreatedAt time.Time    `bun:"created_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time    `bun:"updated_at,nullzero,notnull,default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt sql.NullTime `bun:"deleted_at,soft_delete,nullzero" json:"-"`
}

func (m *Realm) List(ctx context.Context) ([]*Realm, error) {
	var realms []*Realm
	sq := runtime.DB.NewSelect().Model(&realms)
	err := sq.Scan(ctx, &realms)
	if err != nil {
		runtime.Logger.Errorf("list realms failed : %s", err)
	}

	return realms, err
}

func (m *Realm) Get(ctx context.Context) error {
	sq := runtime.DB.NewSelect().Model(m).Limit(1)
	if m.ID != "" {
		sq = sq.Where("id = ?", m.ID)
	}

	if m.Name != "" {
		sq = sq.Where("name = ?", m.Name)
	}

	err := sq.Scan(ctx, m)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("query non-exists realm <%s>", m.ID)
		} else {
			runtime.Logger.Errorf("query realm failed : %s", err)
		}
	}

	return err
}

func (m *Realm) Create(ctx context.Context) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	if m.Status != RealmStatusValid {
		m.Status = RealmStatusInvalid
	}

	iq := runtime.DB.NewInsert().Model(m).Returning("*")
	_, err := iq.Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("insert realm failed : %s", err)
	}

	return err
}

func (m *Realm) Update(ctx context.Context) error {
	uq := runtime.DB.NewUpdate().Model(m).Where("id = ?", m.ID)
	if m.Name != "" {
		uq = uq.Set("name = ?", m.Name)
	}

	if m.Status != RealmStatusValid {
		m.Status = RealmStatusInvalid
	}

	uq = uq.Set("status = ?", m.Status).Set("updated_at = CURRENT_TIMESTAMP")
	_, err := uq.Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("update realm failed : %s", err)
	}

	return err
}

func (m *Realm) Delete(ctx context.Context) error {
	dq := runtime.DB.NewDelete().Model(m).Where("id = ?", m.ID)
	_, err := dq.Exec(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			runtime.Logger.Warnf("delete non-exists realm <%s>", m.ID)

			return nil
		}

		runtime.Logger.Errorf("delete realm failed : %s", err)
	}

	return err
}

func (m *Realm) Init(ctx context.Context) error {
	_, err := runtime.DB.NewCreateTable().Model(m).IfNotExists().Exec(ctx)
	if err != nil {
		runtime.Logger.Errorf("Create table <realms> failed : %s", err)

		return err
	}

	runtime.DB.NewCreateIndex().Model(m).Index("idx_realms_created_at").Column("created_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_realms_updated_at").Column("updated_at").Exec(ctx)
	runtime.DB.NewCreateIndex().Model(m).Index("idx_realms_deleted_at").Column("deleted_at").Exec(ctx)

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
