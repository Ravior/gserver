package ghttp

import (
	"fmt"
	"github.com/Ravior/gserver/os/glog"
	"github.com/Ravior/gserver/util/gconfig"
	"net/http"
	"os"
)

const (
	ServerStatusStopped = 0
	ServerStatusRunning = 1
)

type Server struct {
	status int       // Status of current server.
	port   int32     // 端口
	name   string    // 服务器名称
	id     string    // 服务器ID
	ip     string    // Host
	exit   chan bool // 退出通道
	router *Router   // 消息路由器
}

func NewServer() *Server {
	return &Server{
		name:   gconfig.Global.ServerId,
		id:     gconfig.Global.ServerName,
		ip:     gconfig.Global.HttpServer.IP,
		port:   gconfig.Global.HttpServer.Port,
		status: ServerStatusStopped,
		exit:   make(chan bool, 1),
		router: NewRouter(),
	}
}

// 获取服务器名称
func (s *Server) GetName() string {
	return s.name
}

// 获取服务器ID
func (s *Server) GetId() string {
	return s.id
}

// 获取服务器IP地址
func (s *Server) GetHost() string {
	return s.ip
}

// 获取服务器端口
func (s *Server) GetPort() int32 {
	return s.port
}

// 启动服务器
func (s *Server) Start() {
	if s.status == ServerStatusRunning {
		return
	}

	// 标记服务器状态为运行中
	s.status = ServerStatusRunning

	// 启动全局监听
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		request := NewRequest(req, resp)
		s.router.Run(request)
	})

	// 开启一个Go协程去监听服务器端口
	go func() {
		addr := fmt.Sprintf("%s:%d", s.ip, s.port)
		// ListenAndServe会阻塞进程
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			// 标记服务器状态为停止
			s.status = ServerStatusStopped
			glog.Errorf("启动Http服务器失败， %+v, 程序即将退出", err)
			os.Exit(0)
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	if s.status == ServerStatusStopped {
		return
	}
	s.exit <- true
}

// 运行服务器
func (s *Server) Run() {
	s.Start()
	// 阻塞, 否则主Go退出
	select {
	case exit := <-s.exit:
		if exit {
			return
		}
	}
}

// 添加路由
func (s *Server) AddRoute(path string, handler HandlerFunc, middleware ...MiddlewareFunc) {
	s.router.group.AddRoute(path, handler, middleware...)
}
