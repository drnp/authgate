/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file redis.go
 * @package runtime
 * @author Dr.NP <np@herewe.tech>
 * @since 11/07/2023
 */

package runtime

import (
	"net"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/redis/v3"
)

var Storage fiber.Storage

func InitStorage() error {
	host, port, _ := net.SplitHostPort(Config.Redis.Addr)
	portNum, _ := strconv.Atoi(port)
	store := redis.New(redis.Config{
		Host:     host,
		Port:     portNum,
		Password: Config.Redis.Password,
		Database: Config.Redis.DB,
	})

	Storage = store

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
