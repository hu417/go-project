package bo

import "rbac-v1/model/po"

type Role struct {
	*po.Role
	Powers []*po.Power `json:"powers"`
}
