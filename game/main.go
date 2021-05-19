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
	flag.StringVar(&app.SID, "sid", util.GetUUID(), "uuid")
	flag.StringVar(&app.ServiceType, "type", "all", "服务类型")
	iAddr := flag.String("iAddr", "127.0.0.1:10001", "服务地址")
	oAddr := flag.String("oAddr", "127.0.0.1:10002", "服务地址")
	rAddr := flag.String("grpcAddr", "127.0.0.1:10003", "grpc服务地址")
	flag.Parse()
	goxh.SetCodec(new(codec.ProtobufCodec))
	goxh.SetInteriorService(new(tcp.TService), *iAddr)
	goxh.SetOutsideService(new(websocket.WService), *oAddr)
	goxh.SetGrpcAddr(*rAddr)
	goxh.SetModule(new(mods.MainModule))
	goxh.Start("xhgo.ini")
	goxh.Shutdown()
}
