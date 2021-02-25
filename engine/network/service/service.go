package service

import (
	"context"
	"sync"

	"github.com/xhaoh94/goxh/engine/network"
	"github.com/xhaoh94/goxh/engine/network/service/servicebase"
)

type (
	//Service 服务器
	Service struct {
		addr          string
		acceptFunc    servicebase.AcceptFn
		Wg            sync.WaitGroup
		IsRun         bool
		Ctx           context.Context
		CtxCancelFunc context.CancelFunc
	}
)

var (
	sessionPool sync.Pool
)

//NewSession 创建session
func (ser *Service) NewSession(channel servicebase.IChannel, t int, delFn func(network.ISession)) network.ISession {
	s := sessionPool.Get().(*Session)
	s.init(ser.Ctx, channel, t, delFn)
	return s
}

//Init 服务初始化
func (ser *Service) Init(addr string, accept servicebase.AcceptFn) {
	ser.Ctx, ser.CtxCancelFunc = context.WithCancel(context.TODO())
	ser.addr = addr
	ser.acceptFunc = accept
	sessionPool = sync.Pool{
		New: func() interface{} {
			return &Session{}
		},
	}
}

//GetAddr 获取地址
func (ser *Service) GetAddr() string {
	return ser.addr
}

//OnAccept 新链接回调
func (ser *Service) OnAccept(channel servicebase.IChannel) {
	if ser.acceptFunc != nil {
		ser.acceptFunc(channel)
	}
}
