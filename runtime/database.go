/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file database.go
 * @package runtime
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package runtime

import (
	"database/sql"
	"time"

	"github.com/alexlast/bunzap"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func InitDB() error {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(Config.Database.DSN)))
	DB = bun.NewDB(sqldb, pgdialect.New())
	if Config.Debug {
		DB.AddQueryHook(bunzap.NewQueryHook(bunzap.QueryHookOptions{
			Logger:       LoggerRaw,
			SlowDuration: 200 * time.Millisecond,
		}))
	}

	err := DB.Ping()
	if err != nil {
		Logger.Fatalf("database connection failed : %s", err)
	}

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
