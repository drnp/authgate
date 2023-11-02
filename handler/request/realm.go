/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file realm.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package request

type RealmPost struct {
	Name string `json:"name" xml:"name"`
}

type RealmPut struct {
	Name   string `json:"name" xml:"name"`
	Status int    `json:"status" xml:"status"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
