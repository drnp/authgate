/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file account.go
 * @package request
 * @author Dr.NP <np@herewe.tech>
 * @since 10/28/2023
 */

package request

type AccountPost struct {
	RealmID  string `json:"realm_id" xml:"realm_id"`
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Mobile   string `json:"mobile" xml:"mobile"`
	Password string `json:"password" xml:"password"`
}

type AccountPut struct {
	Username string `json:"username" xml:"username"`
	Email    string `json:"email" xml:"email"`
	Mobile   string `json:"mobile" xml:"mobile"`
	Password string `json:"password" xml:"password"`
	Status   int    `json:"status" xml:"status"`
}

type AccountAuth struct {
	RealmID  string `json:"realm_id" xml:"realm_id"`
	Username string `json:"username,omitempty" xml:"username,omitempty"`
	Email    string `json:"email,omitempty" xml:"email,omitempty"`
	Mobile   string `json:"mobile,omitempty" xml:"mobile,omitempty"`
	Password string `json:"password" xml:"password"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
