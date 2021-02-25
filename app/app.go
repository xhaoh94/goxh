package app

import (
	"runtime"
)

type ServiceType string

var (
	//Version 版本
	Version string
	//SID 服务id
	SID string
	//SType 服务类型
	SType ServiceType
	//OutsideAddr 外部服务地址
	OutsideAddr string
	//InteriorAddr 内部服务地址
	InteriorAddr string
	//RPCAddr GRPC服务地址
	RPCAddr string
)

//GetRuntime 运行平台
func GetRuntime() string {
	return runtime.GOOS
}
