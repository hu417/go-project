package flag

import "gvb_server/global"

func Version() {
	global.Logger.Infof("系统版本: %s", global.Config.System.Version)
}
