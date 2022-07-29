package gtcp

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ravior/gserver/net/gnet"
	"github.com/Ravior/gserver/os/glog"
	"github.com/gorilla/websocket"
	"io"
	"net"
	"sync/atomic"
)

type Connection struct {
	isClosed   int32              // 当前链接的关闭状态(采用原子操作处理)
	connID     uint32             // 当前链接的ID, 也可以称作为SessionID，ID全局唯一
	msgChan    chan []byte        // 缓冲管道，用于读、写两个goroutine之间的消息通信
	msgHandler *gnet.MsgHandler   // 消息处理模块
	socket     gnet.ISocket       // 当前链接关联的Socket
	conn       *net.TCPConn       // 当前链接的TCP套接字
	ctx        context.Context    // 告知该链接已经退出/停止的channel
	cancel     context.CancelFunc // cancelFunc
}

// NewConnection 创建新的链接对象
func NewConnection(socket gnet.ISocket, conn *net.TCPConn, connID uint32, msgHandler *gnet.MsgHandler, maxMsgChanLen uint32) *Connection {
	c := &Connection{
		socket:     socket,
		conn:       conn,
		connID:     connID,
		isClosed:   0,
		msgHandler: msgHandler,
		msgChan:    make(chan []byte, maxMsgChanLen),
	}

	if conn != nil {
		// 将新创建的Conn添加到链接管理器
		c.socket.GetConnMgr().Add(c)
	} else {
		c.isClosed = 1
	}

	return c
}

// StartWriter 写消息Goroutine， 用户将数据发送给客户端
func (c *Connection) StartWriter() {
	defer func() {
		glog.Infof("Connection Writer Close, ConnId:%d", c.connID)
		if err := recover(); err != nil {
			e := fmt.Sprintf("%v", err)
			glog.Errorf("Connection write loop has error:%v, ConnId:%d, Addr:%s", e, c.connID, c.RemoteAddr())
		}
	}()

	defer c.Stop()

	for {
		select {
		case data, ok := <-c.msgChan:
			if ok {
				// 有数据要写给客户端
				if _, err := c.conn.Write(data); err != nil {
					glog.Warnf("Connection write message has error: %s, ConnId:%d, Addr:%s 即将断开", err.Error(), c.connID, c.RemoteAddr())
					return
				}
			} else {
				glog.Warnf("MsgChan has been closed, ConnId:%d, Addr:%s 即将断开", c.connID, c.RemoteAddr())
				break
			}
		case <-c.ctx.Done():
			glog.Debugf("Connection Context Is Cancel, Stop Writer, ConnId:%d, Addr:%s 即将断开", c.connID, c.RemoteAddr())
			return
		}
	}
}

// StartReader 读消息Goroutine，用于从客户端中读取数据
func (c *Connection) StartReader() {
	defer func() {
		glog.Infof("Connection Reader Close, ConnId:%d", c.connID)
		if err := recover(); err != nil {
			e := fmt.Sprintf("%v", err)
			glog.Errorf("Connection read loop has error:%v, ConnId:%d", e, c.connID)
		}
	}()

	defer c.Stop()

	for {
		select {
		case <-c.ctx.Done():
			glog.Debugf("Connection Context Is Cancel, Stop Reader, ConnId:%d, Addr:%s 即将断开", c.connID, c.RemoteAddr())
			return
		default:
			// 读取客户端的Msg head
			headData := make([]byte, MsgPack.GetHeadLen())

			if _, err := io.ReadFull(c.conn, headData); err != nil {
				glog.Warnf("Connection read message has error: %s, ConnId:%d, (%s)即将断开", err.Error(), c.connID, c.RemoteAddr())
				return
			}

			// 普通socket拆包，得到nameLen 和 dataLen 放在msg中
			msg, err := MsgPack.Unpack(headData)
			if err != nil {
				glog.Errorf("Connection unpack message fail，err:%s, ConnId:%d, (%s)即将断开", err.Error(), c.connID, c.RemoteAddr())
				return
			}

			// 根据 nameLen bodyLen 读取 CaveList
			var data []byte
			if msg.GetDataLen() > 0 {
				data = make([]byte, msg.GetDataLen())
				if _, err := io.ReadFull(c.conn, data); err != nil {
					glog.Errorf("Connection read message fail，err:%s, ConnId:%d, (%s)即将断开", err.Error(), c.connID, c.RemoteAddr())
					return
				}
			}
			msg.SetData(data)

			// 得到当前客户端请求的Request数据
			req := gnet.NewRequest(c, msg)
			// 将收到消息交给Worker处理
			c.msgHandler.SendMsgToTaskQueue(req)
		}
	}
}

//============== 实现 interfaces.IConnection 里的全部接口方法 ========

func (c *Connection) GetProtocolType() gnet.ProtocolType {
	return gnet.WebSocket
}

func (c *Connection) Start() {
	if c.conn != nil {
		c.ctx, c.cancel = context.WithCancel(context.Background())
		// 开启一个Go协程，从客户端读取数据
		go c.StartReader()
		// 开启一个Go协程，写回数据到客户端
		go c.StartWriter()

		// 触发Socket中Conn Start钩子方法
		c.socket.CallOnConnStart(c)
	}
}

//Stop 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	glog.Debugf("停止连接, ConnId:%d, Addr:%s", c.connID, c.RemoteAddr())

	// 标记为已关闭状态, 如果已经关闭则不执行任何操作
	if c.SetClosed() == false {
		return
	}

	glog.Infof("执行关闭链接操作, ConnId:%d, Addr:%s, IsClosed:%v", c.connID, c.RemoteAddr(), c.isClosed)

	c.cancel()

	// 关闭该链接全部管道
	close(c.msgChan)

	// 关闭socket链接
	_ = c.conn.Close()

	// 将链接从连接管理器中删除
	c.socket.GetConnMgr().Remove(c)

	// 触发Socket中Conn Stop钩子方法(放在go内，防止回调出现死锁）
	go c.socket.CallOnConnStop(c)
}

func (c *Connection) GetConnID() uint32 {
	return c.connID
}

func (c *Connection) SetClosed() bool {
	return atomic.SwapInt32(&c.isClosed, 1) == 0
}

func (c *Connection) IsClosed() bool {
	return atomic.LoadInt32(&c.isClosed) == 1
}

func (c *Connection) GetTcpConnection() *net.TCPConn {
	return c.conn
}

func (c *Connection) GetWsConnection() *websocket.Conn {
	return nil
}

func (c *Connection) GetSocket() gnet.ISocket {
	return c.socket
}

func (c *Connection) RemoteAddr() net.Addr {
	if c.conn != nil {
		return c.conn.RemoteAddr()
	}
	return nil
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if data == nil {
		return errors.New("connection send nil msg")
	}

	// 链接已关闭
	if c.IsClosed() {
		glog.Warnf("Connection has been closed when send msg, ConnId:%d, Addr:%s", c.connID, c.RemoteAddr())
		return errors.New("connection has been closed when send msg")
	}

	// 将data封包，并且发送
	p := gnet.NewMsg(msgId, data)
	msg, err := MsgPack.Pack(p)
	if err != nil {
		glog.Errorf("Connection pack message fail，msgId:%d, msgData:%v, err:%s, ConnId:%d, Addr:%s", msgId, data, err.Error(), c.connID, c.RemoteAddr())
		return errors.New(fmt.Sprintf("connection pack message fail, err:%v", err))
	}

	// 检测通道关闭Panic
	defer func() {
		if err := recover(); err != nil {
			e := fmt.Sprintf("%v", err)
			glog.Errorf("Connection SendMsg has error:%v, ConnId:%d, Addr:%s", e, c.connID, c.RemoteAddr())
		}
	}()

	select {
	case c.msgChan <- msg:
	case <-c.ctx.Done():
		glog.Debugf("Connection Context Done, ConnId:%d, Addr:%s", c.connID, c.RemoteAddr())
	default:
		glog.Warnf("Connection Send Msg Error, I/O Buff Is Full Or Can't Write Msg To Client, MsgId:%d, ConnId:%d, Addr:%s", msgId, c.connID, c.RemoteAddr())
		return errors.New("I/O Buff Is Full Or Can't Write Msg To Client")
	}

	return nil
}
