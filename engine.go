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
	"github.com/xhaoh94/goxh/engine/xlog"
	"github.com/xhaoh94/goxh/util"
)

//Init 初始化 sid:服务id mType:服务类型 interior:内部服务地址(不可空) outside:外部服务地址(可空) grpc:rpc服务地址(可空)
func Init(sid string, mType app.ServiceType, interiorAddr string, outsideAddr string, grpcAddr string) {
	if sid == "" {
		app.SID = util.GetUUID()
	} else {
		app.SID = sid
	}

	app.SType = mType
	app.InteriorAddr = interiorAddr
	app.OutsideAddr = outsideAddr
	app.RPCAddr = grpcAddr
}

//Start 启动
func Start(appConfPath string) {
	xlog.Init()
	if app.SID == "" {
		xlog.Error("It needs to be done first Init ")
		return
	}
	config.LoadAppConfig(appConfPath)

	xlog.Info("server start. sid -> [%s]", app.SID)
	xlog.Info("server type -> [%s]", app.SType)
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
func SetOutsideService(ser network.IService) {
	network.SetOutsideService(ser)
}

//SetInteriorService 设置内部服务类型
func SetInteriorService(ser network.IService) {
	network.SetInteriorService(ser)
}
