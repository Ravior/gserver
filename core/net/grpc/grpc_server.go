package grpc

import (
	"fmt"
	"github.com/Ravior/gserver/core/os/glog"
	"github.com/Ravior/gserver/core/util/gconfig"
	"google.golang.org/grpc"
	"net"
	"os"
)

type Server struct {
	name      string           // 服务器名称
	ip        string           // Host
	port      int32            // 端口
	ipVersion string           // IP版本，"gtcp"、"tcp4"或"tcp6"
	listener  *net.TCPListener // 服务器TCP监听器
	Server    *grpc.Server     // GRPC服务器
	exit      chan bool        // 退出通道
}

func NewServer() *Server {
	server := &Server{
		name:      gconfig.Global.ServerName,
		ipVersion: "tcp4",
		ip:        gconfig.Global.RpcServer.IP,
		port:      gconfig.Global.RpcServer.Port,
		exit:      make(chan bool, 1),
	}
	return server
}

// GetHost 获取服务器IP地址
func (s *Server) GetHost() string {
	return s.ip
}

// GetPort 获取服务器端口
func (s *Server) GetPort() int32 {
	return s.port
}

// GetName 获取服务器名称
func (s *Server) GetName() string {
	return s.name
}

// Start 启动服务器
func (s *Server) Start() {
	glog.Debugf("Server: %s StartWork", s.GetName())
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

		// 实例化grpc Server
		grpServer := grpc.NewServer()

		s.Server = grpServer

		grpServer.Serve(listener)

	}()
}

func (s *Server) Stop() {
	_ = s.listener.Close()

	s.exit <- true
}

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
