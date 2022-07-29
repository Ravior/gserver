package gconsole

import (
	"fmt"
	"github.com/Ravior/gserver/core/net/gnet"
	"github.com/Ravior/gserver/core/os/glog"
	"github.com/Ravior/gserver/core/util/gconfig"
	"net"
	"os"
	"sync/atomic"
)

// ConnCallback 连接回调
type ConnCallback func(conn *Connection)

// Server 定义一个Server服务类，实现interfaces.IServer接口
type Server struct {
	ipVersion string // IP版本，"tcp"、"tcp4"或"tcp6"
	ip        string // Host
	port      int32  // 端口

	listener *net.TCPListener // 服务器TCP监听器
	exit     chan bool        // 退出通道

	onConnStart ConnCallback // 有新的客户端链接时触发的Hook函数
	onConnStop  ConnCallback // 当客户端链接断开时触发的Hook函数
}

func NewServer() *Server {
	server := &Server{
		ipVersion: "tcp4",
		ip:        gconfig.Global.Console.IP,
		port:      gconfig.Global.Console.Port,
		exit:      make(chan bool, 1),
	}
	return server
}

// Start 启动服务器
func (s *Server) Start() {

	// 开启一个Go协程去监听服务器端口
	go func() {
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
		glog.Infof("Console Server: %s StartWork", addr)

		var connID uint32 = 0

		for {
			// 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				return
			}

			// 创建链接对象
			dealConn := NewConnection(s, conn, connID)

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
	_ = s.listener.Close()
	s.exit <- true
}

// Run 运行服务器
func (s *Server) Run() {
	s.Start()
	// 阻塞,否则主Go退出， listenner的go将会退出
	select {
	case exit := <-s.exit:
		if exit {
			return
		}
	}
}

func (s *Server) GetRouter() *gnet.Router {
	return nil
}

// SetOnConnStart 设置服务器有新的链接Hook函数
func (s *Server) SetOnConnStart(connCallback ConnCallback) {
	s.onConnStart = connCallback
}

// SetOnConnStop 设置服务器有链接断开Hook函数
func (s *Server) SetOnConnStop(connCallback ConnCallback) {
	s.onConnStop = connCallback
}

// CallOnConnStart 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn *Connection) {
	if s.onConnStart != nil {
		s.onConnStart(conn)
	}
}

// CallOnConnStop 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn *Connection) {
	if s.onConnStop != nil {
		s.onConnStop(conn)
	}
}
