/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file account.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 10/28/2023
 */

package response

import "time"

/* {{{ [Response codes && messages] */
const (
	CodeListAccountFailed   = 50500001
	CodeGetAccountFailed    = 50500002
	CodeCreateAccountFailed = 50500003
	CodeUpdateAccountFailed = 50500004
	CodeDeleteAccountFailed = 50500005
)

const (
	MsgListAccountFailed   = "List account failed"
	MsgGetAccountFailed    = "Get account failed"
	MsgCreateAccountFailed = "Create account failed"
	MsgUpdateAccountFailed = "Update account failed"
	MsgDeleteAccountFailed = "Delete account failed"
)

type AccountGet struct {
	ID        string    `json:"id" xml:"id"`
	RealmID   string    `json:"realm_id" xml:"realm_id"`
	Username  string    `json:"username" xml:"username"`
	Email     string    `json:"email" xml:"email"`
	Mobile    string    `json:"mobile" xml:"mobile"`
	Status    int       `json:"status" xml:"status"`
	CreatedAt time.Time `json:"created_at" xml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" xml:"updated_at"`
}

type AccountPost struct {
	ID       string `json:"id" xml:"id"`
	RealmID  string `json:"realm_id" xml:"realm_id"`
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Mobile   string `json:"mobile" xml:"mobile"`
	Status   int    `json:"status" xml:"status"`
}

/* }}} */

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
