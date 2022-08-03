package gwebsocket

import (
	"fmt"
	"github.com/Ravior/gserver/net/gnet"
	"github.com/Ravior/gserver/os/glog"
	"github.com/Ravior/gserver/os/gtimer"
	"github.com/Ravior/gserver/util/gconfig"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

// IOBufferBytesSize will be used when reading messages from clients
// 参数值参考PitaYa
const IOBufferBytesSize = 4096

type Server struct {
	name          string                                                 // 服务器名称
	id            string                                                 // 服务器ID
	ipVersion     string                                                 // IP版本，"tcp"、"tcp4"或"tcp6"
	ip            string                                                 // Host
	certFile      string                                                 // SSL证书文件路径
	keyFile       string                                                 // SSL密钥文件路径
	port          int32                                                  // 端口
	connID        uint32                                                 // 链接ID
	exit          chan bool                                              // 退出通道
	listener      *net.TCPListener                                       // 服务器TCP监听器
	connMgr       *gnet.ConnManager                                      // 链接管理器
	router        *gnet.Router                                           // 消息路由器
	msgHandler    *gnet.MsgHandler                                       // 当前Server的消息管理模块，用来绑定消息ID和对应的处理方法
	onConnCheck   func(resp http.ResponseWriter, req *http.Request) bool // WebSocket链接校验判断
	onConnUpgrade func(conn *Connection, req *http.Request)              // Http协议升级为WebSocket协议触发的Hook函数
	onConnStart   gnet.ConnCallback                                      // 有新的客户端链接时触发的Hook函数
	onConnStop    gnet.ConnCallback                                      // 当客户端链接断开时触发的Hook函数
}

func NewServer() *Server {
	server := &Server{
		id:         gconfig.Global.ServerId,
		name:       gconfig.Global.ServerName,
		ipVersion:  "tcp4",
		ip:         gconfig.Global.WsServer.IP,
		port:       gconfig.Global.WsServer.Port,
		certFile:   gconfig.Global.WsServer.CertFile,
		keyFile:    gconfig.Global.WsServer.KeyFile,
		connMgr:    gnet.NewConnManager(),
		exit:       make(chan bool, 1),
		router:     &gnet.Router{},
		msgHandler: gnet.NewMsgHandler(gconfig.Global.WsServer.WorkerPoolSize, gconfig.Global.WsServer.WorkerTaskLen),
	}
	server.msgHandler.SetRouter(server.router)

	// 每30秒输出一次当前在线人数
	gtimer.Add(30*time.Second, func() {
		glog.Infof("Server Online:%d", server.connMgr.Len())
	})

	return server
}

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	HandshakeTimeout: 2 * time.Second,   // 设置2秒钟超时时间
	ReadBufferSize:   IOBufferBytesSize, // io 操作的缓存大小，如果不指定就会自动分配
	WriteBufferSize:  IOBufferBytesSize, // 写数据操作的缓存池，如果没有设置值，write buffers 将会分配到链接生命周期里
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// SetOnConnCheck 设置链接检查函数
func (s *Server) SetOnConnCheck(onConnCheck func(resp http.ResponseWriter, req *http.Request) bool) {
	s.onConnCheck = onConnCheck
}

// SetOnConnUpgrade 设置链接协议升级回调函数
func (s *Server) SetOnConnUpgrade(onConnUpgrade func(conn *Connection, req *http.Request)) {
	s.onConnUpgrade = onConnUpgrade
}

func (s *Server) wsHandler(resp http.ResponseWriter, req *http.Request) {
	glog.Debugf("收到WebSocket链接请求,Addr:%s", req.RemoteAddr)
	if s.onConnCheck != nil {
		// 如果校验失败则不进行链接
		if s.onConnCheck(resp, req) == false {
			glog.Errorf("[服务器错误] 链接校验失败,Addr:%s", req.RemoteAddr)
			return
		}
	}

	glog.Debugf("WebSocket链接参数校验成功,Addr:%s", req.RemoteAddr)

	// 最大连接数判断
	if gconfig.Global.WsServer.MaxConn > 0 && s.GetConnMgr().Len() >= gconfig.Global.WsServer.MaxConn {
		glog.Errorf("[超出服务器设定最大连接数限制] 当前连接数:%d, 最大连接数:%d", s.GetConnMgr().Len(), gconfig.Global.WsServer.MaxConn)
		return
	}

	// 应答客户端告知升级连接为websocket
	conn, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		glog.Warnf("WebSocket升级协议失败, Err:%v, Addr:%s", err, req.RemoteAddr)
		return
	}

	// 原子+1(此处需要先操作，避免出现相同ConnID的问题，故ConnID从1开始)
	atomic.AddUint32(&s.connID, 1)
	// 创建链接对象
	dealConn := NewConnection(s, conn, s.connID, s.msgHandler, gconfig.Global.WsServer.MaxMsgChanLen)

	if s.onConnUpgrade != nil {
		// 执行链接升级回调
		s.onConnUpgrade(dealConn, req)
	}

	glog.Infof("WebSocket链接成功, 开始处理业务, ConnId:%d, Addr:%s", dealConn.connID, req.RemoteAddr)

	// 启动当前链接的处理业务
	go dealConn.Start()
}

// GetName 获取服务器名称
func (s *Server) GetName() string {
	return s.name
}

// GetId 获取服务器ID
func (s *Server) GetId() string {
	return s.id
}

// GetHost 获取服务器IP地址
func (s *Server) GetHost() string {
	return s.ip
}

// GetPort 获取服务器端口
func (s *Server) GetPort() int32 {
	return s.port
}

// Run 运行服务器
func (s *Server) Run() {
	s.Start()
	// 阻塞,否则主Go退出， listener的go将会退出
	select {
	case exit := <-s.exit:
		if exit {
			return
		}
	}
}

// Start 启动服务器
func (s *Server) Start() {
	glog.Debugf("Server: %s StartWork", s.GetName())
	// 开启一个Go协程去监听服务器端口
	go func() {
		// 启动消息Worker工作池
		s.msgHandler.StartWorkerPool()

		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			fmt.Println("服务器地址格式错误, Error:", err)
			os.Exit(0)
		}

		http.HandleFunc("/", s.wsHandler)

		// 监听服务器地址
		if s.certFile != "" && s.keyFile != "" {
			glog.Debugf("Websocket Server StartWork. URL:wss://%s, certFile = %s, keyFile = %s", addr.String(), s.certFile, s.keyFile)
			err = http.ListenAndServeTLS(addr.String(), s.certFile, s.keyFile, nil)
		} else {
			glog.Debugf("Websocket Server StartWork. URL:ws://%s", addr.String())
			err = http.ListenAndServe(addr.String(), nil)
		}

		if err != nil {
			fmt.Println("服务器启动失败, Error:", err)
			os.Exit(0)
		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	// 关闭worker工作池
	s.msgHandler.StopWorkerPool()
	s.connMgr.ClearConn()
	_ = s.listener.Close()

	s.exit <- true
}

func (s *Server) GetRouter() *gnet.Router {
	return s.router
}

// GetConnMgr 获取链接管理器
func (s *Server) GetConnMgr() *gnet.ConnManager {
	return s.connMgr
}

// SetOnConnStart 设置服务器有新的链接Hook函数
func (s *Server) SetOnConnStart(connCallback gnet.ConnCallback) {
	s.onConnStart = connCallback
}

// SetOnConnStop 设置服务器有链接断开Hook函数
func (s *Server) SetOnConnStop(connCallback gnet.ConnCallback) {
	s.onConnStop = connCallback
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn gnet.IConnection) {
	if s.onConnStart != nil {
		s.onConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn gnet.IConnection) {
	if s.onConnStop != nil {
		s.onConnStop(conn)
	}
}
