/*
 * Copyright (C) HereweTech, Inc - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

/**
 * @file oidc.go
 * @package handler
 * @author Dr.NP <np@herewe.tech>
 * @since 10/26/2023
 */

package handler

type OIDC struct {
}

func InitOIDC() *OIDC {
	h := new(OIDC)

	return h
}

/*
 * Local variables:
 * tab-width: 4
 * c-basic-offset: 4
 * End:
 * vim600: sw=4 ts=4 fdm=marker
 * vim<600: sw=4 ts=4
 */
