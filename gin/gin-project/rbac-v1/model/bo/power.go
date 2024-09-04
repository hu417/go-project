package bo

import "rbac-v1/model/po"

type Power struct {
	*po.Power
	Operations []*po.Operation `json:"operations"`
}
