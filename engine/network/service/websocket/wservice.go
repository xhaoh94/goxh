package websocket

import (
	"net/http"
	"net/url"
	"time"

	"github.com/xhaoh94/goxh/app"
	"github.com/xhaoh94/goxh/engine/network/service"
	"github.com/xhaoh94/goxh/engine/network/service/servicebase"
	"github.com/xhaoh94/goxh/engine/xlog"

	"github.com/gorilla/websocket"
)

//WService WebSocket服务器
type WService struct {
	service.Service
	upgrader websocket.Upgrader
	server   *http.Server
}

//Init 服务初始化
func (ws *WService) Init(addr string, accept servicebase.AcceptFn) {
	ws.Service.Init(addr, accept)
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", ws.wsPage)
	ws.server = &http.Server{Addr: ws.GetAddr(), Handler: mux}
	ws.upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

}

//Start 启动
func (ws *WService) Start() {
	xlog.Info("websocket service Waiting for clients. -> [%s]", ws.GetAddr())
	go ws.accept()
}
func (ws *WService) accept() {
	defer ws.Wg.Done()
	ws.IsRun = true
	ws.Wg.Add(1)

	err := ws.server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			xlog.Info("websocket close")
		} else {
			xlog.Error("websocket ListenAndServe err: [%s]", err.Error())
		}
	}
}
func (ws *WService) wsPage(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		xlog.Error("websocket wsPage: [%s]", err.Error())
		return
	}
	go ws.connection(conn)
}

func (ws *WService) connection(conn *websocket.Conn) {
	wChannel := ws.addChannel(conn)
	ws.OnAccept(wChannel)
}
func (ws *WService) addChannel(conn *websocket.Conn) (wChannel *WChannel) {
	wChannel = channelPool.Get().(*WChannel)
	wChannel.init(conn)
	return
}

//ConnectChannel 链接新信道
func (ws *WService) ConnectChannel(addr string) servicebase.IChannel {
	var connCount int
	for {
		u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
		var dialer *websocket.Dialer
		conn, _, err := dialer.Dial(u.String(), nil)
		if err == nil {
			return ws.addChannel(conn)
		}
		if connCount > app.ReConnectMax {
			xlog.Info("websocket create channel fail addr:[%s] err:[%v]", addr, err)
			return nil
		}
		time.Sleep(app.ReConnectInterval)
		connCount++
		continue
	}

}

//Stop 停止服务
func (ws *WService) Stop() {
	if !ws.IsRun {
		return
	}
	ws.IsRun = false
	ws.server.Shutdown(ws.Ctx)
	ws.CtxCancelFunc()
	// 等待线程结束
	ws.Wg.Wait()

}
