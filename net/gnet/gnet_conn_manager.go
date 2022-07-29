package gnet

import (
	"errors"
	"go.uber.org/atomic"
	"sync"
)

// ConnManager 链接管理器
type ConnManager struct {
	connections sync.Map
	connNum     atomic.Int32 // 当前连接数
}

// NewConnManager 创建新的链接管理器
func NewConnManager() *ConnManager {
	return &ConnManager{}
}

// Add 添加链接
func (c *ConnManager) Add(conn IConnection) {
	// 保存链接信息
	c.connections.Store(conn.GetConnID(), conn)
	// 链接数+1
	c.connNum.Inc()
}

// Remove 删除链接
func (c *ConnManager) Remove(conn IConnection) {
	if _, ok := c.connections.Load(conn.GetConnID()); ok {
		// 链接数-1
		c.connNum.Dec()
	}
	c.connections.Delete(conn.GetConnID())
}

// Get 根据ConnID获取链接
func (c *ConnManager) Get(connID uint32) (IConnection, error) {
	if _conn, ok := c.connections.Load(connID); ok {
		return _conn.(IConnection), nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// ClearConn 删除并停止所有链接
func (c *ConnManager) ClearConn() {
	c.connections.Range(func(k interface{}, v interface{}) bool {
		if connId, ok := k.(uint32); ok {
			c.connections.Delete(connId)
			c.connNum.Dec()
		}
		if conn, ok := v.(IConnection); ok {
			conn.Stop()
		}
		return true
	})
}

func (c *ConnManager) Len() int32 {
	return c.connNum.Load()
}

// ClearOneConn 删除指定ConnID的链接
func (c *ConnManager) ClearOneConn(connID uint32) {
	conn, _ := c.Get(connID)
	if conn != nil {
		c.connNum.Dec()
		conn.Stop()
	}
	c.connections.Delete(connID)
}

// BroadcastMsg 广播数据
func (c *ConnManager) BroadcastMsg(msgId uint32, data []byte) {
	if data != nil && msgId > 0 {
		c.connections.Range(func(k interface{}, v interface{}) bool {
			if conn, ok := v.(IConnection); ok {
				conn.SendMsg(msgId, data)
			}
			return true
		})
	}
}
