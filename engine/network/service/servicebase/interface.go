package servicebase

type (
	//AcceptFn 收到新链接回调
	AcceptFn func(IChannel)
	//WriteFn 写回调
	WriteFn func([]byte)
	//ReadFn 读回调
	ReadFn func([]byte)
	//CloseFn 关闭回调
	CloseFn func()

	//IChannel 信道接口
	IChannel interface {
		Start()
		Send(data []byte)
		Stop()
		RemoteAddr() string
		LocalAddr() string
		SetCallBackFn(ReadFn, CloseFn)
	}
	//IDefaultRPC rpc
	IDefaultRPC interface {
		Await() bool
	}
)
