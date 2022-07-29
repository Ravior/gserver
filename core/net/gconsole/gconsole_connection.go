package gconsole

import (
	"errors"
	"fmt"
	"github.com/Ravior/gserver/core/net/gconsole/command"
	"github.com/Ravior/gserver/core/net/gnet"
	"github.com/Ravior/gserver/core/net/gtcp"
	"github.com/Ravior/gserver/core/os/glog"
	"io"
	"net"
	"strings"
)

type Connection struct {
	server *Server      // 当前链接关联的Socket
	conn   *net.TCPConn // 当前链接的TCP套接字
	connID uint32       // 当前链接的ID, 也可以称作为SessionID，ID全局唯一
}

// NewConnection 创建新的链接对象
func NewConnection(server *Server, conn *net.TCPConn, connID uint32) *Connection {
	c := &Connection{
		server: server,
		conn:   conn,
		connID: connID,
	}

	return c
}

// StartReader StartTcpReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	defer c.Stop()

	for {
		// 读取客户端的Msg head
		headData := make([]byte, gtcp.MsgPack.GetHeadLen())

		if _, err := io.ReadFull(c.conn, headData); err != nil {
			break
		}

		// 普通socket拆包，得到nameLen 和 dataLen 放在msg中
		msg, err := gtcp.MsgPack.Unpack(headData)
		if err != nil {
			break
		}

		// 根据 nameLen bodyLen 读取 CaveList
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.conn, data); err != nil {
				break
			}
		}
		msg.SetData(data)

		line := string(msg.Data)

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		if args[0] == "exit" {
			break
		}

		var cmd command.Command
		for _, _cmd := range command.Mgr.GetCommands() {
			if _cmd.Name() == args[0] {
				cmd = _cmd
				break
			}
		}
		if cmd == nil {
			if err := c.SendMsg([]byte("command not found, try `help` for help\r\n")); err != nil {
				break
			}
			continue
		}
		output := cmd.Run(args[1:])
		if output != "" {
			if err := c.SendMsg([]byte(output + "\r\n")); err != nil {
				break
			}
		}
	}
}

func (c *Connection) Start() {
	if c.conn != nil {
		// 开启一个Go协程，从客户端读取数据
		go c.StartReader()

		// 触发Socket中Conn Start钩子方法
		c.server.CallOnConnStart(c)
	}
}

func (c *Connection) Stop() {
	// 关闭socket链接
	_ = c.conn.Close()
	// 触发Socket中Conn Stop钩子方法(放在go内，防止回调出现死锁）
	go c.server.CallOnConnStop(c)
}

func (c *Connection) RemoteAddr() net.Addr {
	if c.conn != nil {
		return c.conn.RemoteAddr()
	}
	return nil
}

// SendMsg 发送消息
func (c *Connection) SendMsg(data []byte) error {
	// 将data封包，并且发送
	p := gnet.NewMsg(0, data)
	msg, err := gtcp.MsgPack.Pack(p)

	if err != nil {
		glog.Errorf("connection pack message fail，msgData:%v, err:%s, connId:%d, addr:%s", data, err.Error(), c.connID, c.RemoteAddr())
		return errors.New(fmt.Sprintf("connection pack message fail, err:%v", err))
	}

	_, err = c.conn.Write(msg)

	return err
}
