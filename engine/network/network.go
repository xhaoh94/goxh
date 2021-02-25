package network

import (
	"sync"

	"github.com/xhaoh94/goxh/app"
	"github.com/xhaoh94/goxh/consts"
	"github.com/xhaoh94/goxh/engine/network/rpc"
	"github.com/xhaoh94/goxh/engine/network/service/servicebase"
	"github.com/xhaoh94/goxh/engine/xlog"
)

type (
	//Network 服务结构体
	Network struct {
		service IService

		idToSession map[string]ISession //Accept Map
		idMutex     sync.Mutex
	}
	//OutsideNetwork 外部服务体
	OutsideNetwork struct {
		Network
	}
)

var (
	isRun    bool
	wg       sync.WaitGroup
	outside  *OutsideNetwork
	interior *InteriorNetwork
)

//SetOutsideService 设置外部服务
func SetOutsideService(ser IService) {
	outside = &OutsideNetwork{}
	outside.service = ser
}

//SetInteriorService 设置内部服务
func SetInteriorService(ser IService) {
	interior = &InteriorNetwork{}
	interior.service = ser
	interior.addrToSession = make(map[string]ISession, 0)
}

//Start 网络初始化入口
func Start() {
	if isRun {
		return
	}
	if interior == nil {
		xlog.Fatal("service is nil")
		return
	}
	if outside != nil {
		outside.start(app.OutsideAddr)
	}

	isRun = true
	interior.start(app.InteriorAddr)
	rpc.Start()
	registerService()
}

//Stop 销毁
func Stop() {
	if !isRun {
		return
	}
	isRun = false
	unRegisterService()
	rpc.Stop()
	if outside != nil {
		outside.stop()
	}
	interior.stop()
	wg.Wait()
}

//GetSession 通过id获取Session
func GetSession(sid string) ISession {
	defer outside.idMutex.Unlock()
	outside.idMutex.Lock()
	session, ok := outside.idToSession[sid]
	if ok {
		return session
	}

	defer interior.idMutex.Unlock()
	interior.idMutex.Lock()
	session, ok = interior.idToSession[sid]
	if ok {
		return session
	}
	return nil
}

func (nw *Network) start(addr string) {
	if addr == "" {
		xlog.Fatal("addr is nil")
		return
	}
	nw.idToSession = make(map[string]ISession)
	nw.service.Init(addr, nw.onAccept)
	nw.service.Start()
}

func (nw *Network) stop() {
	nw.idMutex.Lock()
	for k := range nw.idToSession {
		nw.idToSession[k].Stop()
	}
	nw.idMutex.Unlock()
	nw.service.Stop()
}

func (nw *Network) onDelete(session ISession) {
	if nw.delSessionByID(session.UID()) {
		wg.Done()
	}
}
func (nw *Network) onAccept(channel servicebase.IChannel) {
	session := nw.onCreateSession(channel, consts.Accept, nw.onDelete)
	if session != nil {
		nw.idMutex.Lock()
		nw.idToSession[session.UID()] = session
		nw.idMutex.Unlock()
		session.Start()
	}
}
func (nw *Network) onCreateSession(channel servicebase.IChannel, t int, fn func(ISession)) ISession {
	session := nw.service.NewSession(channel, t, fn)
	if session != nil {
		wg.Add(1)
	}
	return session
}

func (nw *Network) delSessionByID(id string) bool {
	defer nw.idMutex.Unlock()
	nw.idMutex.Lock()
	if _, ok := nw.idToSession[id]; ok {
		delete(nw.idToSession, id)
		return true
	}
	return false
}
