package gtcp

import (
	"fmt"
	"github.com/Ravior/gserver/core/net/gnet"
	"github.com/Ravior/gserver/core/os/glog"
	"github.com/Ravior/gserver/core/util/gconfig"
	"net"
	"os"
	"sync/atomic"
)

// Server 定义一个Server服务类，实现interfaces.IServer接口
type Server struct {
	name        string            // 服务器名称
	id          string            // 服务器ID
	ipVersion   string            // IP版本，"tcp"、"tcp4"或"tcp6"
	ip          string            // Host
	port        int32             // 端口
	exit        chan bool         // 退出通道
	listener    *net.TCPListener  // 服务器TCP监听器
	connMgr     *gnet.ConnManager // 链接管理器
	router      *gnet.Router      // 消息路由器
	msgHandler  *gnet.MsgHandler  // 当前Server的消息管理模块，用来绑定消息ID和对应的处理方法
	onConnStart gnet.ConnCallback // 有新的客户端链接时触发的Hook函数
	onConnStop  gnet.ConnCallback // 当客户端链接断开时触发的Hook函数
}

func NewServer() *Server {
	server := &Server{
		id:         gconfig.Global.ServerId,
		name:       gconfig.Global.ServerName,
		ipVersion:  "tcp4",
		ip:         gconfig.Global.TcpServer.IP,
		port:       gconfig.Global.TcpServer.Port,
		connMgr:    gnet.NewConnManager(),
		router:     &gnet.Router{},
		exit:       make(chan bool, 1),
		msgHandler: gnet.NewMsgHandler(gconfig.Global.TcpServer.WorkerPoolSize, gconfig.Global.TcpServer.WorkerTaskLen),
	}
	server.msgHandler.SetRouter(server.router)
	return server
}

//============== 实现 interfaces.INetEndPoint 里的全部接口方法 ========

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

// Start 启动服务器
func (s *Server) Start() {
	glog.Infof("Server: %s StartWork", s.GetName())
	// 开启一个Go协程去监听服务器端口
	go func() {
		// 启动消息Worker工作池
		s.msgHandler.StartWorkerPool()

		// 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			os.Exit(0)
		}

		// 监听服务器地址
		listener, err := net.ListenTCP(s.ipVersion, addr)
		if err != nil {
			os.Exit(0)
		}
		s.listener = listener

		var connID uint32 = 0

		// 处理TCP链接
		for {
			// 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				return
			}

			// 最大连接数判断
			if gconfig.Global.TcpServer.MaxConn > 0 && s.GetConnMgr().Len() >= gconfig.Global.TcpServer.MaxConn {
				continue
			}

			// 创建链接对象
			dealConn := NewConnection(s, conn, connID, s.msgHandler, gconfig.Global.TcpServer.MaxMsgChanLen)

			// 原子+1
			atomic.AddUint32(&connID, 1)

			// 启动当前链接的处理业务
			go dealConn.Start()

		}
	}()
}

// Stop 停止服务器
func (s *Server) Stop() {
	//关闭worker工作池
	s.msgHandler.StopWorkerPool()
	s.connMgr.ClearConn()
	_ = s.listener.Close()

	s.exit <- true
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
