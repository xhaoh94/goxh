package login

import (
	"context"

	"github.com/xhaoh94/goxh/engine/module"
	"github.com/xhaoh94/goxh/engine/network"
	"github.com/xhaoh94/goxh/engine/network/types"
	"github.com/xhaoh94/goxh/engine/xlog"
	"github.com/xhaoh94/goxh/game/netpack"
)

type (
	//LoginModule 主模块
	LoginModule struct {
		module.Module
	}
)

//OnInit 初始化
func (mm *LoginModule) OnInit() {
	network.Register(111, mm.test)
	network.RegisterRPC(mm.test2)
}

func (m *LoginModule) test(ctx context.Context, session types.ISession, req *netpack.ReqTest) {
	xlog.Info("接受的消息%v", req)
	session.Send(2222, &netpack.RspTest{Token: "test"})
}

func (m *LoginModule) test2(ctx context.Context, req *netpack.ReqTest) *netpack.RspTest {
	xlog.Info("接受的消息2%v", req)
	return &netpack.RspTest{Token: "500"}
}
