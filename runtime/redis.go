/*
 * Copyright (C) LiangYu, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file redis.go
 * @package runtime
 * @author Dr.NP <zhanghao@liangyu.ltd>
 * @since 10/26/2023
 */

package runtime

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var Redis *redis.Client

func InitRedis() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Addr,
		Password: Config.Redis.Password,
		DB:       Config.Redis.DB,
	})
	err := rdb.Ping(context.TODO()).Err()
	if err != nil {
		Logger.Fatalf("redis connection failed: %s", err)
	}

	Redis = rdb

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
