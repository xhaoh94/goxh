package main

import (
	"flag"

	"github.com/xhaoh94/goxh"
	"github.com/xhaoh94/goxh/app"
	"github.com/xhaoh94/goxh/engine/codec"
	"github.com/xhaoh94/goxh/engine/network/service/tcp"
	"github.com/xhaoh94/goxh/engine/network/service/websocket"
	"github.com/xhaoh94/goxh/game/mods"
	"github.com/xhaoh94/goxh/util"
)

func main() {
	flag.StringVar(&app.SID, "id", util.GetUUID(), "uuid")
	flag.StringVar((*string)(&app.SType), "type", "all", "服务类型")
	flag.StringVar(&app.InteriorAddr, "interiorAddr", "127.0.0.1:10001", "服务地址")
	flag.StringVar(&app.OutsideAddr, "outsideAddr", "127.0.0.1:10002", "服务地址")
	flag.StringVar(&app.RPCAddr, "rpcaddr", "127.0.0.1:10003", "rpc服务地址")
	flag.Parse()
	goxh.SetCodec(new(codec.ProtobufCodec))
	goxh.SetInteriorService(new(tcp.TService))
	goxh.SetOutsideService(new(websocket.WService))
	goxh.SetModule(new(mods.MainModule))
	goxh.Start("xhgo.ini")
	goxh.Shutdown()
}
