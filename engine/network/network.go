package network

import (
	"sync"

	"github.com/xhaoh94/goxh/consts"
	"github.com/xhaoh94/goxh/engine/network/rpc"
	"github.com/xhaoh94/goxh/engine/network/types"
	"github.com/xhaoh94/goxh/engine/xlog"
)

var (
	idToSession   map[string]types.ISession //Accept Map
	idMutex       sync.Mutex
	addrToSession map[string]types.ISession //Connect Map
	addrMutex     sync.Mutex
	isRun         bool
	wg            sync.WaitGroup
	outside       types.IService
	interior      types.IService
)

//SetOutsideService 设置外部服务
func SetOutsideService(ser types.IService, addr string) {
	outside = ser
	outside.Init(addr, onAccept)
}

//SetInteriorService 设置内部服务
func SetInteriorService(ser types.IService, addr string) {
	interior = ser
	interior.Init(addr, onAccept)
}

//SetGrpcAddr 设置grpc服务
func SetGrpcAddr(addr string) {
	rpc.Init(addr)
}

//Start 网络初始化入口
func Start() {
	if isRun {
		return
	}
	idToSession = make(map[string]types.ISession)
	addrToSession = make(map[string]types.ISession)
	var outsideAddr, interiorAddr, rpcAddr string
	if interior == nil {
		xlog.Fatal("service is nil")
		return
	}
	interiorAddr = interior.GetAddr()
	rpcAddr = rpc.GetAddr()

	isRun = true
	if outside != nil {
		outsideAddr = outside.GetAddr()
		outside.Start()
	}
	interior.Start()
	rpc.Start()

	registerService(outsideAddr, interiorAddr, rpcAddr)
}

//Stop 销毁
func Stop() {
	if !isRun {
		return
	}

	isRun = false
	unRegisterService()
	rpc.Stop()

	idMutex.Lock()
	for k := range idToSession {
		idToSession[k].Stop()
	}
	idMutex.Unlock()

	addrMutex.Lock()
	for k := range addrToSession {
		addrToSession[k].Stop()
	}
	addrMutex.Unlock()

	if outside != nil {
		outside.Stop()
	}
	interior.Stop()
	wg.Wait()
}

//GetSession 通过id获取Session
func GetSession(sid string) types.ISession {
	defer idMutex.Unlock()
	idMutex.Lock()
	session, ok := idToSession[sid]
	if ok {
		return session
	}
	return nil
}

//GetSessionByAddr 通过地址获取Session
func GetSessionByAddr(addr string) types.ISession {

	defer addrMutex.Unlock()
	addrMutex.Lock()
	if s, ok := addrToSession[addr]; ok {
		return s
	}
	s := onConnect(addr)
	if s == nil {
		xlog.Error("create session fail addr:[%s]", addr)
		return nil
	}
	addrToSession[addr] = s
	return s
}
func onConnect(addr string) types.ISession {
	channel := interior.ConnectChannel(addr)
	if channel != nil {
		session := onCreateSession(channel, consts.Connector)
		if session != nil {
			idMutex.Lock()
			idToSession[session.UID()] = session
			idMutex.Unlock()
			session.Start()
			return session
		}
	}
	return nil
}

func onDelete(session types.ISession) {
	if delSession(session) {
		wg.Done()
	}
}

func delSession(session types.ISession) bool {
	delSessionByAddr(session.RemoteAddr())
	if ok := delSessionByID(session.UID()); ok {
		return true
	}
	return false
}

func delSessionByID(id string) bool {
	defer idMutex.Unlock()
	idMutex.Lock()
	if _, ok := idToSession[id]; ok {
		delete(idToSession, id)
		return true
	}
	return false
}

func delSessionByAddr(addr string) {
	addrMutex.Lock()
	_, ok := addrToSession[addr]
	if ok {
		delete(addrToSession, addr)
	}
	addrMutex.Unlock()
}

func onAccept(channel types.IChannel) {
	session := onCreateSession(channel, consts.Accept)
	if session != nil {
		idMutex.Lock()
		idToSession[session.UID()] = session
		idMutex.Unlock()
		session.Start()
	}
}

func onCreateSession(channel types.IChannel, tag int) types.ISession {
	session := channel.GetService().NewSession(channel, tag, onDelete)
	if session != nil {
		wg.Add(1)
	}
	return session
}
