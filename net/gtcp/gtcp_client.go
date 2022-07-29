package gtcp

import (
	"fmt"
	gnet2 "github.com/Ravior/gserver/net/gnet"
	"github.com/Ravior/gserver/os/glog"
	"net"
)

var (
	defaultWorkerPoolSize uint32 = 1
	defaultWorkerTaskSize uint32 = 10
	defaultMaxMsgChanLen  uint32 = 10
)

// Client 客户端暂时只支持TCP链接
type Client struct {
	name       string
	id         string
	ipVersion  string
	remoteIP   string
	remotePort int32
	msgHandler *gnet2.MsgHandler
	router     *gnet2.Router // 消息路由器

	connMgr     *gnet2.ConnManager
	onConnStart gnet2.ConnCallback
	onConnStop  gnet2.ConnCallback
}

func NewClient(clientId string, clientName string, remoteIP string, remotePort int32) *Client {
	client := &Client{
		id:         clientId,
		name:       clientName,
		ipVersion:  "tcp4",
		remoteIP:   remoteIP,
		remotePort: remotePort,
		router:     &gnet2.Router{},
		msgHandler: gnet2.NewMsgHandler(defaultWorkerPoolSize, defaultWorkerTaskSize),
		connMgr:    gnet2.NewConnManager(),
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

	addr, err := net.ResolveTCPAddr(c.ipVersion, fmt.Sprintf("%s:%d", c.GetHost(), c.GetPort()))
	if err != nil {
		glog.Warnf("Resolve TCP Addr Err:%v", err.Error())
		return
	}

	connServer, err := net.DialTCP(c.ipVersion, nil, addr)
	if err != nil {
		glog.Warnf("Connect To Server Fail, Addr: %v, Err:%v", addr, err.Error())
		conn := NewConnection(c, nil, 0, c.msgHandler, defaultMaxMsgChanLen)
		c.CallOnConnStop(conn)
		return
	}

	// 保证Client的时候只有一个Conn
	c.connMgr.ClearConn()
	conn := NewConnection(c, connServer, 0, c.msgHandler, defaultMaxMsgChanLen)
	conn.Start()
}

func (c *Client) GetRouter() *gnet2.Router {
	return c.router
}

func (c *Client) GetConn() gnet2.IConnection {
	conn, err := c.connMgr.Get(0)
	if err == nil {
		return conn
	} else {
		return nil
	}
}

func (c *Client) GetConnMgr() *gnet2.ConnManager {
	return c.connMgr
}

func (c *Client) SetOnConnStart(connCallback gnet2.ConnCallback) {
	c.onConnStart = connCallback
}

func (c *Client) SetOnConnStop(connCallback gnet2.ConnCallback) {
	c.onConnStop = connCallback
}

func (c *Client) CallOnConnStart(conn gnet2.IConnection) {
	if c.onConnStart != nil {
		c.onConnStart(conn)
	}
}

func (c *Client) CallOnConnStop(conn gnet2.IConnection) {
	if c.onConnStop != nil {
		c.onConnStop(conn)
	}
}
