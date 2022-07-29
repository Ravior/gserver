package gnet

// Request 请求抽象
type Request struct {
	Conn IConnection // 已经和客户端建立好的 链接
	Msg  *Msg        // 客户端请求的数据
}

func NewRequest(conn IConnection, msg *Msg) *Request {
	return &Request{
		Conn: conn,
		Msg:  msg,
	}
}

// GetConnection 获取连接信息
func (r *Request) GetConnection() IConnection {
	return r.Conn
}

// GetConnId 获取链接ID
func (r *Request) GetConnId() uint32 {
	return r.Conn.GetConnID()
}

// GetMessage 获取请求的消息
func (r *Request) GetMessage() *Msg {
	return r.Msg
}
