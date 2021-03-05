package types

import (
	"github.com/xhaoh94/goxh/engine/network/rpc"
)

type (
	//IService 服务器接口
	IService interface {
		Init(string, func(IChannel))
		Start()
		Stop()
		GetAddr() string
		OnAccept(IChannel)
		NewSession(IChannel, int, func(ISession)) ISession
		ConnectChannel(string) IChannel
	}
	//IChannel 信道接口
	IChannel interface {
		Start()
		Stop()
		Send(data []byte)
		RemoteAddr() string
		LocalAddr() string
		SetCallBackFn(func([]byte), func())
		GetService() IService
	}
	//ISession 会话接口
	ISession interface {
		Start()
		Stop()
		UID() string
		RemoteAddr() string
		LocalAddr() string
		Send(uint32, interface{})
		Call(interface{}, interface{}) rpc.IDefaultRPC
		Reply(interface{}, uint32)
		Actor(uint32, uint32, interface{})
		SendData([]byte)
	}
)
