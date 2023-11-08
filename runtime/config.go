/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file config.go
 * @package runtime
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package runtime

import (
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
)

const (
	AppName   = "authgate"
	EnvPrefix = "zzauth"
)

type mainConfig struct {
	HTTP struct {
		ListenAddr         string `json:"listen_addr" mapstructure:"listen_addr"`
		Prefork            bool   `json:"prefork" mapstructure:"prefork"`
		LongPollingTimeout int64  `json:"long_polling_timeout" mapstructure:"long_polling_timeout"` // In second
	} `json:"http" mapstructure:"http"`
	Database struct {
		DSN string `json:"dsn" mapstructure:"dsn"`
	} `json:"database" mapstructure:"database"`
	Nats struct {
		URL string `json:"url" mapstructure:"url"`
	} `json:"nats" mapstructure:"nats"`
	Redis struct {
		Addr     string `json:"addr" mapstructure:"addr"`
		Password string `json:"password" mapstructure:"password"`
		DB       int    `json:"db" mapstructure:"db"`
	}
	Auth struct {
		JWTAccessSecret     string `json:"jwt_access_secret" mapstructure:"jwt_access_secret"`
		JWTRefreshSecret    string `json:"jwt_refresh_secret" mapstructure:"jwt_refresh_secret"`
		JWTAccessExpiry     int64  `json:"jwt_access_expiry" mapstructure:"jwt_access_expiry"`         // In second
		JWTRefreshExpiry    int64  `json:"jwt_refresh_expiry" mapstructure:"jwt_refresh_expiry"`       // In second
		AuthorizeCodeExpiry int64  `json:"authorize_code_expiry" mapstructure:"authorize_code_expiry"` // In second
	} `json:"auth" mapstructure:"auth"`
	Debug bool `json:"debug" mapstructure:"debug"`

	// Additional
	ZZAuth struct {
		BaseURL string `json:"base_url" mapstructure:"base_url"`
	} `json:"zzauth" mapstructure:"zzauth"`
}

var Config mainConfig

var defaultConfigs = map[string]interface{}{
	"http.listen_addr":           ":9900",
	"http.prefork":               false,
	"http.long_polling_timeout":  30,
	"database.dsn":               "postgres://postgres@localhost:5432/postgres?sslmode=disable",
	"nats.url":                   nats.DefaultURL,
	"redis.addr":                 "localhost:6379",
	"auth.jwt_access_secret":     "access_secret",
	"auth.jwt_refresh_secret":    "refresh_secret",
	"auth.jwt_access_expiry":     2 * 60 * 60,
	"auth.jwt_refresh_expiry":    30 * 24 * 60 * 60,
	"auth.authorize_code_expiry": 5 * 60,
	"debug":                      false,

	"zzauth.base_url": "http://zzauth.herewe.tech",
}

func LoadConfig() error {
	for cfgKey, cfgVal := range defaultConfigs {
		viper.SetDefault(cfgKey, cfgVal)
	}

	viper.SetConfigName(AppName)
	viper.SetConfigType("json")
	viper.AddConfigPath("/etc/" + AppName)
	viper.AddConfigPath("$HOME/." + AppName)
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		// Logging
		Logger.Warnf("reading config file error : %s", err)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix(EnvPrefix)
	viper.AutomaticEnv()
	err = viper.Unmarshal(&Config)
	if err != nil {
		LoggerRaw.Fatal(err.Error())
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
