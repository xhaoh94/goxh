package mods

import (
	"github.com/xhaoh94/goxh/app"
	"github.com/xhaoh94/goxh/engine/module"
	"github.com/xhaoh94/goxh/game/mods/gate"
	"github.com/xhaoh94/goxh/game/mods/login"
)

const (
	//Gate 网关服务
	Gate app.ServiceType = "gate"
	//Login 登录服务
	Login app.ServiceType = "login"
)

type (
	//MainModule 主模块
	MainModule struct {
		module.Module
	}
)

//OnInit 初始化
func (mm *MainModule) OnInit() {
	switch app.SType {
	case Gate:
		mm.Put(&gate.GateModule{})
		break
	case Login:
		mm.Put(&login.LoginModule{})
		break
	}
}
