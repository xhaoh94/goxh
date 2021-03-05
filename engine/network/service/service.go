package service

import (
	"context"
	"sync"

	"github.com/xhaoh94/goxh/engine/network/types"
)

type (
	//Service 服务器
	Service struct {
		addr          string
		acceptFunc    func(types.IChannel)
		Wg            sync.WaitGroup
		IsRun         bool
		Ctx           context.Context
		CtxCancelFunc context.CancelFunc
	}
)

var (
	sessionPool *sync.Pool
)

func init() {
	sessionPool = &sync.Pool{
		New: func() interface{} {
			return &Session{}
		},
	}
}

func (ser *Service) GetSessionPool() *sync.Pool {
	return sessionPool
}

//NewSession 创建session
func (ser *Service) NewSession(channel types.IChannel, tag int, delFn func(types.ISession)) types.ISession {
	s := sessionPool.Get().(*Session)
	s.init(ser.Ctx, channel, tag, delFn)
	return s
}

//Init 服务初始化
func (ser *Service) Init(addr string, accept func(types.IChannel)) {
	ser.Ctx, ser.CtxCancelFunc = context.WithCancel(context.TODO())
	ser.addr = addr
	ser.acceptFunc = accept
}

//GetAddr 获取地址
func (ser *Service) GetAddr() string {
	return ser.addr
}

//OnAccept 新链接回调
func (ser *Service) OnAccept(channel types.IChannel) {
	if ser.acceptFunc != nil {
		ser.acceptFunc(channel)
	}
}
