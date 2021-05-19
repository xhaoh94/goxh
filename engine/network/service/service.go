package service

import (
	"context"
	"sync"

	"github.com/xhaoh94/goxh/consts"
	"github.com/xhaoh94/goxh/engine/network/types"
	"github.com/xhaoh94/goxh/engine/xlog"
)

type (
	//Service 服务器
	Service struct {
		addr          string
		idToSession   map[string]*Session //Accept Map
		idMutex       sync.Mutex
		addrToSession map[string]*Session //Connect Map
		addrMutex     sync.Mutex
		sessionWg     sync.WaitGroup

		ConnectChannelFunc func(addr string) types.IChannel
		AcceptWg           sync.WaitGroup
		IsRun              bool
		Ctx                context.Context
		CtxCancelFunc      context.CancelFunc
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

//Init 服务初始化
func (ser *Service) Init(addr string) {
	ser.Ctx, ser.CtxCancelFunc = context.WithCancel(context.TODO())
	ser.addr = addr
	ser.idToSession = make(map[string]*Session)
	ser.addrToSession = make(map[string]*Session)
}

//GetAddr 获取地址
func (ser *Service) GetAddr() string {
	return ser.addr
}

//OnAccept 新链接回调
func (ser *Service) OnAccept(channel types.IChannel) {
	session := ser.createSession(channel, consts.Accept)
	if session != nil {
		ser.idMutex.Lock()
		ser.idToSession[session.UID()] = session
		ser.idMutex.Unlock()
		session.start()
	}
}

//GetSession 通过id获取Session
func (ser *Service) GetSession(sid string) types.ISession {
	defer ser.idMutex.Unlock()
	ser.idMutex.Lock()
	session, ok := ser.idToSession[sid]
	if ok {
		return session
	}
	return nil
}

//GetSessionByAddr 通过addr地址获取Session
func (ser *Service) GetSessionByAddr(addr string) types.ISession {
	defer ser.addrMutex.Unlock()
	ser.addrMutex.Lock()
	if s, ok := ser.addrToSession[addr]; ok {
		return s
	}
	session := ser.onConnect(addr)
	if session == nil {
		xlog.Error("create session fail addr:[%s]", addr)
		return nil
	}
	ser.idMutex.Lock()
	ser.idToSession[session.UID()] = session
	ser.idMutex.Unlock()
	ser.addrToSession[addr] = session
	session.start()
	return session
}

//Stop 停止服务
func (ser *Service) Stop() {
	ser.idMutex.Lock()
	for k := range ser.idToSession {
		ser.idToSession[k].stop()
	}
	ser.idMutex.Unlock()

	ser.addrMutex.Lock()
	for k := range ser.addrToSession {
		ser.addrToSession[k].stop()
	}
	ser.addrMutex.Unlock()
	ser.sessionWg.Wait()
}

func (ser *Service) onDelete(session types.ISession) {
	if ser.delSession(session) {
		ser.sessionWg.Done()
	}
}

func (ser *Service) delSession(session types.ISession) bool {
	ser.delSessionByAddr(session.RemoteAddr())
	if ok := ser.delSessionByID(session.UID()); ok {
		return true
	}
	return false
}

func (ser *Service) delSessionByID(id string) bool {
	defer ser.idMutex.Unlock()
	ser.idMutex.Lock()
	if _, ok := ser.idToSession[id]; ok {
		delete(ser.idToSession, id)
		return true
	}
	return false
}

func (ser *Service) delSessionByAddr(addr string) {
	ser.addrMutex.Lock()
	_, ok := ser.addrToSession[addr]
	if ok {
		delete(ser.addrToSession, addr)
	}
	ser.addrMutex.Unlock()
}

func (ser *Service) onConnect(addr string) *Session {
	channel := ser.ConnectChannelFunc(addr)
	if channel != nil {
		return ser.createSession(channel, consts.Connector)
	}
	return nil
}

func (ser *Service) createSession(channel types.IChannel, tag consts.SessionTag) *Session {
	session := sessionPool.Get().(*Session)
	session.init(ser, channel, tag)
	if session != nil {
		ser.sessionWg.Add(1)
	}
	return session
}
