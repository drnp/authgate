/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file Realm.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 10/27/2023
 */

package response

import "time"

/* {{{ [Response codes && messages] */
const (
	CodeListRealmFailed   = 30500001
	CodeGetRealmFailed    = 30500002
	CodeCreateRealmFailed = 30500003
	CodeUpdateRealmFailed = 30500004
	CodeDeleteRealmFailed = 30500005
)

const (
	MsgListRealmFailed   = "List Realm failed"
	MsgGetRealmFailed    = "Get Realm failed"
	MsgCreateRealmFailed = "Create Realm failed"
	MsgUpdateRealmFailed = "Update Realm failed"
	MsgDeleteRealmFailed = "Delete Realm failed"
)

/* }}} */

type RealmGet struct {
	ID        string    `json:"id" xml:"id"`
	Name      string    `json:"name" xml:"name"`
	Status    int       `json:"status" xml:"status"`
	CreatedAt time.Time `json:"created_at" xml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at"`
}

type RealmPost struct {
	ID     string `json:"id" xml:"id"`
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
