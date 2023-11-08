/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file Realm.go
 * @package response
 * @author Dr.NP <np@herewe.tech>
 * @since 11/08/2023
 */

package response

import "time"

type PostToken struct {
	ClientID              string    `json:"client_id" xml:"client_id"`
	AccessToken           string    `json:"access_token" xml:"access_token"`
	AccessTokenExpiresAt  time.Time `json:"access_token_expires_at" xml:"access_token_expires_at"`
	RefreshToken          string    `json:"refresh_token,omitempty" xml:"refresh_token,omitempty"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at,omitempty" xml:"refresh_token_expires_at,omitempty"`
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
