package network

import (
	"github.com/xhaoh94/goxh/engine/network/service/servicebase"
)

type (

	//IService 服务器接口
	IService interface {
		Init(string, servicebase.AcceptFn)
		Start()
		GetAddr() string
		OnAccept(servicebase.IChannel)
		Stop()
		NewSession(servicebase.IChannel, int, func(ISession)) ISession
		ConnectChannel(string) servicebase.IChannel
	}
	//ISession 会话接口
	ISession interface {
		Start()
		Stop()
		UID() string
		RemoteAddr() string
		LocalAddr() string
		Send(uint32, interface{})
		Call(interface{}, interface{}) servicebase.IDefaultRPC
		Reply(interface{}, uint32)
		Actor(uint32, uint32, interface{})
		SendData([]byte)
	}
)
