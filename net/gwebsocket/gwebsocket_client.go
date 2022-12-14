package gwebsocket

import (
	"fmt"
	"github.com/Ravior/gserver/internal/empty"
	"github.com/Ravior/gserver/net/gnet"
	"github.com/Ravior/gserver/os/glog"
	"github.com/gorilla/websocket"
)

var (
	defaultWorkerPoolSize uint32 = 1
	defaultWorkerTaskSize uint32 = 10
	defaultMaxMsgChanLen  uint32 = 100
)

type Client struct {
	name        string
	id          string
	scheme      string
	remotePath  string
	remoteIP    string
	remotePort  int32
	router      *gnet.Router
	msgHandler  *gnet.MsgHandler
	connMgr     *gnet.ConnManager
	onConnStart func(conn gnet.IConnection)
	onConnStop  func(conn gnet.IConnection)
}

func NewClient(clientName string, clientId string, remoteIP string, remotePort int32) *Client {
	client := &Client{
		name:       clientName,
		id:         clientId,
		remoteIP:   remoteIP,
		remotePort: remotePort,
		router:     &gnet.Router{},
		msgHandler: gnet.NewMsgHandler(defaultWorkerPoolSize, defaultWorkerTaskSize),
		connMgr:    gnet.NewConnManager(),
	}
	client.msgHandler.SetRouter(client.router)
	return client
}

//============== 实现 interfaces.ISocket 里的全部接口方法 ========

// GetName 获取客户端名称
func (c *Client) GetName() string {
	return c.name
}

// GetId 获取客户端ID
func (c *Client) GetId() string {
	return c.id
}

func (c *Client) SetSchema(scheme string) {
	c.scheme = scheme
}

func (c *Client) SetRemotePath(path string) {
	c.remotePath = path
}

func (c *Client) GetSchema() string {
	if empty.IsEmpty(c.scheme) {
		c.scheme = "ws"
	}
	return c.scheme
}

func (c *Client) GetRemotePath() string {
	if empty.IsEmpty(c.remotePath) {
		c.remotePath = "/"
	}
	return c.remotePath
}

func (c *Client) GetHost() string {
	return c.remoteIP
}

func (c *Client) GetPort() int32 {
	return c.remotePort
}

// Start 开启
func (c *Client) Start() {
	c.msgHandler.StartWorkerPool()
}

// Stop 关闭
func (c *Client) Stop() {
	c.msgHandler.StopWorkerPool()
	c.connMgr.ClearConn()
}

func (c *Client) Run() {
	c.Start()

	addr := fmt.Sprintf("%s://%s:%d%s", c.GetSchema(), c.GetHost(), c.GetPort(), c.GetRemotePath())
	connServer, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		glog.Warnf("WebSocket客户端链接失败，错误消息:%v", err.Error())
		conn := NewConnection(c, nil, 0, c.msgHandler, defaultMaxMsgChanLen)
		c.CallOnConnStop(conn)
		return
	}

	// 保证Client的时候只有一个Conn
	c.connMgr.ClearConn()
	conn := NewConnection(c, connServer, 0, c.msgHandler, defaultMaxMsgChanLen)
	conn.Start()
}

func (c *Client) GetRouter() *gnet.Router {
	return c.router
}

func (c *Client) GetConn() gnet.IConnection {
	conn, err := c.connMgr.Get(0)
	if err == nil {
		return conn
	} else {
		return nil
	}
}

func (c *Client) GetConnMgr() *gnet.ConnManager {
	return c.connMgr
}

func (c *Client) SetOnConnStart(connCallback gnet.ConnCallback) {
	c.onConnStart = connCallback
}

func (c *Client) SetOnConnStop(connCallback gnet.ConnCallback) {
	c.onConnStop = connCallback
}

func (c *Client) CallOnConnStart(conn gnet.IConnection) {
	if c.onConnStart != nil {
		glog.Infof("Client CallOnConnStart, ConnId:%d, Addr:%s", conn.GetConnID(), conn.RemoteAddr())
		c.onConnStart(conn)
	}
}

func (c *Client) CallOnConnStop(conn gnet.IConnection) {
	if c.onConnStop != nil {
		glog.Infof("Client CallOnConnStop, ConnId:%d, Addr:%s", conn.GetConnID(), conn.RemoteAddr())
		c.onConnStop(conn)
	}
}
