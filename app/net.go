package app

import (
	"encoding/binary"
	"time"

	"github.com/gorilla/websocket"
)

var (

	//NetEndian 传输类型
	NetEndian binary.ByteOrder

	//WebSocketMessageType ws使用的消息类型(使用websocket才有效)
	WebSocketMessageType int = websocket.BinaryMessage

	//SendMsgMaxLen 发送最大长度(websocket的话不能超过126) 默认0 不分片
	SendMsgMaxLen int = 0

	//ReadMsgMaxLen 包体最大长度
	ReadMsgMaxLen int = 1 << 24 //16MB

	//ConnectInterval 链接间隔
	ReConnectInterval time.Duration = (1 * time.Second)

	//ConnectMax 尝试链接最大次数
	ReConnectMax int = 3

	//Heartbeat 心跳时间
	Heartbeat time.Duration = (30 * time.Second)

	//ConnectTimeout 链接超时
	ConnectTimeout time.Duration = (3 * time.Second)

	//ReadTimeout
	ReadTimeout time.Duration = (35 * time.Second)
)
