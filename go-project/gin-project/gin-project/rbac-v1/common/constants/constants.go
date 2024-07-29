package constants

import "time"

const (
	JWT_SALT = "addevops"
	JWT_TOKEN_EXP = 24 * time.Hour

	TOKEN_HEADER_KEY = "Authorization"

	CTX_USER_ID = "user_id"
	CTX_USERNAME = "username"

	RBAC_MODEL = `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
	`
)