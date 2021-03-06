package goxh

import (
	"os"
	"os/signal"

	"github.com/xhaoh94/goxh/app"
	"github.com/xhaoh94/goxh/engine/codec"
	"github.com/xhaoh94/goxh/engine/config"
	"github.com/xhaoh94/goxh/engine/module"
	"github.com/xhaoh94/goxh/engine/network"
	"github.com/xhaoh94/goxh/engine/network/actor"
	"github.com/xhaoh94/goxh/engine/network/types"
	"github.com/xhaoh94/goxh/engine/xlog"
)

//Start 启动
func Start(appConfPath string) {
	config.LoadAppConfig(appConfPath)
	xlog.Init()
	if app.SID == "" {
		xlog.Error("It needs to be done first Init ")
		return
	}
	xlog.Info("server start. sid -> [%s]", app.SID)
	xlog.Info("server type -> [%s]", app.ServiceType)
	xlog.Info("server version -> [%s]", app.Version)
	network.Start()
	actor.Start()
	module.Start()
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	<-sigChan
}

//Shutdown 关闭
func Shutdown() {
	module.Stop()
	actor.Stop()
	network.Stop()
	xlog.Info("server exit. sid -> [%s]", app.SID)
	xlog.Destroy()
	os.Exit(1)
}

//SetModule 设置主模块
func SetModule(m module.IModule) {
	module.SetModule(m)
}

//SetCodec 设置解编码器
func SetCodec(c codec.ICodec) {
	codec.SetCodec(c)
}

//SetOutsideService 设置外部服务类型
func SetOutsideService(ser types.IService, addr string) {
	network.SetOutsideService(ser, addr)
}

//SetInteriorService 设置内部服务类型
func SetInteriorService(ser types.IService, addr string) {
	network.SetInteriorService(ser, addr)
}

//SetGrpcAddr 设置grpc服务
func SetGrpcAddr(addr string) {
	network.SetGrpcAddr(addr)
}
