package gconsole

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/Ravior/gserver/core/net/gnet"
	"github.com/Ravior/gserver/core/net/gtcp"
	"github.com/Ravior/gserver/core/os/glog"
	"io"
	"net"
	"os"
)

type Client struct {
	ipVersion  string
	remoteIP   string
	remotePort int32
	conn       *net.TCPConn
	reader     *bufio.Reader
}

func NewClient(remoteIP string, remotePort int32) *Client {
	return &Client{
		ipVersion:  "tcp4",
		remoteIP:   remoteIP,
		remotePort: remotePort,
	}
}

func (c *Client) Start() {
	addr, err := net.ResolveTCPAddr(c.ipVersion, fmt.Sprintf("%s:%d", c.remoteIP, c.remotePort))
	if err != nil {
		fmt.Printf("Resolve TCP Addr Err:%v\r\n", err.Error())
		os.Exit(0)
	}

	conn, err := net.DialTCP(c.ipVersion, nil, addr)
	if err != nil {
		fmt.Printf("Connect To Server Fail, Addr: %v, Err:%v\r\n", addr, err.Error())
		os.Exit(0)
	}
	fmt.Println("链接Console Server成功，请输入命令继续...")
	c.conn = conn

	go c.StartTcpReader()
	go c.StartConsoleReader()

}

func (c *Client) StartTcpReader() {
	for {
		headData := make([]byte, gtcp.MsgPack.GetHeadLen())

		if _, err := io.ReadFull(c.conn, headData); err != nil {
			break
		}

		msg, err := gtcp.MsgPack.Unpack(headData)
		if err != nil {
			break
		}

		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.conn, data); err != nil {
				break
			}
		}

		msg.SetData(data)
		line := string(msg.Data)
		fmt.Println(line)
		fmt.Print(">> 请输入命令继续，或输入【exit】退出终端...\r\n\r\n")
	}
}

func (c *Client) StartConsoleReader() {
	for {
		in := bufio.NewReader(os.Stdin)
		data, _, err := in.ReadLine()
		if err != nil {
			break
		}
		if string(data) == "exit" {
			fmt.Println(">> 终端即将退出...")
			os.Exit(0)
		}
		fmt.Print(">> 获取Console Server数据中，请稍后...\r\n\r\n")

		if err := c.SendMsg(data); err != nil {
			fmt.Printf("发送命令到Console Server失败，Err:%v", err)
			os.Exit(0)
		}
	}
}
func (c *Client) Run() {
	c.Start()
}

// SendMsg 发送消息
func (c *Client) SendMsg(data []byte) error {
	// 将data封包，并且发送
	p := gnet.NewMsg(0, data)
	msg, err := gtcp.MsgPack.Pack(p)

	if err != nil {
		glog.Errorf("connection pack message fail，msgData:%v, err:%s", data, err.Error())
		return errors.New(fmt.Sprintf("connection pack message fail, err:%v", err))
	}

	_, err = c.conn.Write(msg)
	return err
}
