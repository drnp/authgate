/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file client.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 10/28/2023
 */

package response

import "time"

/* {{{ [Response codes && messages] */
const (
	CodeListClientFailed   = 40500001
	CodeGetClientFailed    = 40500002
	CodeCreateClientFailed = 40500003
	CodeUpdateClientFailed = 40500004
	CodeDeleteClientFailed = 40500005
)

const (
	MsgListClientFailed   = "List client failed"
	MsgGetClientFailed    = "Get client failed"
	MsgCreateClientFailed = "Create client failed"
	MsgUpdateClientFailed = "Update client failed"
	MsgDeleteClientFailed = "Delete client failed"
)

/* }}} */

type ClientGet struct {
	ID           string    `json:"id" xml:"id"`
	RealmID      string    `json:"realm_id" xml:"realm_id"`
	Name         string    `json:"name" xml:"name"`
	AccessKey    string    `json:"access_key" xml:"access_key"`
	AccessSecret string    `json:"-" xml:"-"`
	Status       int       `json:"status" xml:"status"`
	CreatedAt    time.Time `json:"created_at" xml:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" xml:"updated_at"`
}

type ClientPost struct {
	ID           string `json:"id" xml:"id"`
	RealmID      string `json:"realm_id" xml:"realm_id"`
	Name         string `json:"name" xml:"name"`
	AccessKey    string `json:"access_key" xml:"access_key"`
	AccessSecret string `json:"access_secret" xml:"access_secret"`
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
