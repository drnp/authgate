/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file main.go
 * @package main
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package main

import (
	"authgate/handler"
	"authgate/model"
	"authgate/runtime"
	"context"
	"os"

	"github.com/urfave/cli/v2"
)

func actionServe(c *cli.Context) error {
	handler.InitMisc()
	// handler.InitAccount()
	// handler.InitClient()
	// handler.InitRealm()
	handler.InitOAuth()
	handler.InitOIDC()

	return runtime.Serve()
}

func actionInitdb(c *cli.Context) error {
	var err error

	ctx := context.TODO()
	mAccount := new(model.Account)
	mClient := new(model.Client)
	mRealm := new(model.Realm)

	err = mAccount.Init(ctx)
	if err != nil {
		return err
	}

	runtime.Logger.Info("Table <accounts> created")

	err = mClient.Init(ctx)
	if err != nil {
		return err
	}

	runtime.Logger.Info("Table <clients> created")

	err = mRealm.Init(ctx)
	if err != nil {
		return err
	}

	runtime.Logger.Info("Table <realms> created")

	return nil
}

// Portal

// @title ZZAuth::Authgate API
// @version 0.0.1
// @description Authgate API
// @contact.name HereweTech CO.LTD
// @contact.url https://herewe.tech
// @contact.email support@herewetech.com

// @host authgate.d.herewe.tech
// @BasePath /
func main() {
	runtime.LoadConfig()
	runtime.InitLogger()
	runtime.InitServer()
	runtime.InitNats()
	//runtime.InitRedis()
	runtime.InitStorage()
	runtime.InitDB()

	app := &cli.App{
		Name: runtime.AppName,
		Commands: []*cli.Command{
			{
				Name:   "serve",
				Usage:  "Run service",
				Action: actionServe,
			},
			{
				Name:   "initdb",
				Usage:  "Initialize database tables",
				Action: actionInitdb,
			},
		},
		DefaultCommand: "serve",
	}

	if err := app.Run(os.Args); err != nil {
		runtime.Logger.Fatal(err)
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
