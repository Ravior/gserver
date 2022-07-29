package grpc

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"os"
)

type Client struct {
	name       string
	id         string
	ipVersion  string
	remoteIP   string
	remotePort int32
	conn       *grpc.ClientConn
}

func NewClient(clientName string, clientId string, remoteIP string, remotePort int32) *Client {
	client := &Client{
		name:       clientName,
		id:         clientId,
		ipVersion:  "tcp4",
		remoteIP:   remoteIP,
		remotePort: remotePort,
	}

	return client
}

func (c *Client) Start() {
	// 获取一个TCP的Addr
	addr, err := net.ResolveTCPAddr(c.ipVersion, fmt.Sprintf("%s:%d", c.remoteIP, c.remotePort))
	if err != nil {
		os.Exit(0)
	}
	conn, err := grpc.Dial(addr.String(), grpc.WithInsecure())
	if err != nil {
		fmt.Println("链接失败")
	}
	c.conn = conn
}

func (c *Client) Stop() {
	c.conn.Close()
}
