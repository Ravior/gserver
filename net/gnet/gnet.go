package gnet

import (
	"github.com/gorilla/websocket"
	"net"
)

// ProtocolType 协议类型
type ProtocolType uint8

// ConnCallback 连接回调
type ConnCallback func(conn IConnection)

// HeartBeatTime 心跳时长 单位:秒
const HeartBeatTime = 300

const (
	Tcp = iota
	WebSocket
	Http
)

// IConnection 定义连接接口
type IConnection interface {
	Start()                                  // 启动连接，让当前连接开始工作
	Stop()                                   // 停止连接，结束当前连接状态
	GetTcpConnection() *net.TCPConn          // 从当前连接获取原始的socket TCPConn
	GetWsConnection() *websocket.Conn        // 从当前连接获取原始的websocket conn
	GetProtocolType() ProtocolType           // 获取链接协议类型, TCP/WebSocket
	GetSocket() ISocket                      // 获取链接的Socket对象
	GetConnID() uint32                       // 获取当前连接ID
	IsClosed() bool                          // 当前链接是否已关闭
	SetClosed() bool                         // 设置关闭状态，设置成功返回true,已关闭则返回false
	RemoteAddr() net.Addr                    // 获取远程客户端地址信息
	SendMsg(msgId uint32, data []byte) error // 发送消息
}

// ISocket Socket抽象接口，可以是Sever端/Client端
type ISocket interface {
	GetName() string                           // 获取Socket对象名称
	GetId() string                             // 获取Socket对象ID
	GetHost() string                           // 获取监听(Server)/链接(Client)的主机地址
	GetPort() int32                            // 获取监听(Server)/链接(Client)的端口
	Start()                                    // 启动
	Stop()                                     // 停止运行
	Run()                                      // 运行
	GetRouter() *Router                        // 获取消息路由器
	GetConnMgr() *ConnManager                  // 获取链接管理器(Client的ConnMgr只会有一个Conn)
	SetOnConnStart(startCallBack ConnCallback) // 设置有新的链接Hook函数
	SetOnConnStop(stopCallback ConnCallback)   // 设置有链接断开Hook函数
	CallOnConnStart(conn IConnection)          // 调用链接OnConnStart Hook函数
	CallOnConnStop(conn IConnection)           // 调用链接OnConnStop Hook函数
}
