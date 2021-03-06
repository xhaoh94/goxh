package actor

import (
	"encoding/json"
	"sync"

	"github.com/xhaoh94/goxh/engine/etcd"
	"github.com/xhaoh94/goxh/engine/xlog"

	"github.com/coreos/etcd/mvcc/mvccpb"
	"go.etcd.io/etcd/clientv3"
)

type actorReg struct {
	actorID   uint32
	serviceID string
	sessionID string
}

var (
	actorPrefix   = "location/actor"
	actorEs       *etcd.EtcdService
	actorRegLock  sync.RWMutex
	keyToActorReg map[string]*actorReg
)

func registerActor() {
	var err error
	actorEs, err = etcd.NewEtcdService(get, put, del)
	if err != nil {
		xlog.Fatal("es etcd group err [%v]", err)
		return
	}
	actorEs.Get(actorPrefix, true)
}
func unRegisterActor() {
	if actorEs != nil {
		actorEs.Close()
	}
}
func newSpaceReg(val []byte) (*actorReg, error) {
	actor := &actorReg{}
	if err := json.Unmarshal(val, actor); err != nil {
		return nil, err
	}
	return actor, nil
}

func get(resp *clientv3.GetResponse) {
	if resp == nil || resp.Kvs == nil {
		return
	}
	for i := range resp.Kvs {
		put(resp.Kvs[i])
	}
}
func put(kv *mvccpb.KeyValue) {
	actorRegLock.Lock()
	defer actorRegLock.Unlock()
	if kv.Value == nil {
		return
	}
	key := string(kv.Key)
	value, err := newSpaceReg(kv.Value)
	if err != nil {
		xlog.Error("put actor err[%v]", err)
		return
	}
	keyToActorReg[key] = value
}
func del(kv *mvccpb.KeyValue) {
	actorRegLock.Lock()
	defer actorRegLock.Unlock()
	key := string(kv.Key)
	if _, ok := keyToActorReg[key]; ok {
		delete(keyToActorReg, key)
	}
}
