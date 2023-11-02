/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 10/28/2023
 */

package request

type ClientPost struct {
	RealmID     string `json:"realm_id" xml:"realm_id"`
	Name        string `json:"name" xml:"name"`
	RedirectURL string `json:"redirect_url" xml:"redirect_url"`
}

type ClientPut struct {
	Name         string `json:"name" xml:"name"`
	AccessSecret string `json:"access_secret" xml:"access_secret"`
	RedirectURL  string `json:"redirect_url" xml:"redirect_url"`
	Status       int    `json:"status" xml:"status"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
