package network

import (
	"sync"

	"github.com/xhaoh94/goxh/consts"
	"github.com/xhaoh94/goxh/engine/xlog"
)

type (
	//InteriorNetwork 内部服务体
	InteriorNetwork struct {
		Network
		addrToSession map[string]ISession //Connect Map
		addrMutex     sync.Mutex
	}
)

//GetSessionByAddr 通过地址获取Session
func GetSessionByAddr(addr string) ISession {

	defer interior.addrMutex.Unlock()
	interior.addrMutex.Lock()
	if s, ok := interior.addrToSession[addr]; ok {
		return s
	}
	s := interior.createSession(addr)
	if s == nil {
		xlog.Error("create session fail addr:[%s]", addr)
		return nil
	}
	interior.addrToSession[addr] = s
	return s
}

func (nw *InteriorNetwork) stop() {
	nw.addrMutex.Lock()
	for k := range nw.addrToSession {
		nw.addrToSession[k].Stop()
	}
	nw.addrMutex.Unlock()
	nw.Network.stop()
}

func (nw *InteriorNetwork) createSession(addr string) ISession {
	channel := nw.service.ConnectChannel(addr)
	if channel != nil {
		session := nw.onCreateSession(channel, consts.Connector, nw.onDelete)
		if session != nil {
			session.Start()
			return session
		}
	}
	return nil
}

func (nw *InteriorNetwork) delSessionByAddr(addr string) bool {
	defer nw.addrMutex.Unlock()
	nw.addrMutex.Lock()
	_, ok := nw.addrToSession[addr]
	if ok {
		delete(nw.addrToSession, addr)
		return true
	}
	return false
}

func (nw *InteriorNetwork) onDelete(session ISession) {
	if nw.delSession(session) {
		wg.Done()
	}
}

func (nw *InteriorNetwork) delSession(session ISession) bool {
	if ok := nw.delSessionByID(session.UID()); !ok {
		return nw.delSessionByAddr(session.RemoteAddr())
	}
	return true
}
