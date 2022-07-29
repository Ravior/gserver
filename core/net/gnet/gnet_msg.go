package gnet

// uint8  : 0 to 255
// uint16 : 0 to 65535
// uint32 : 0 to 4294967295
// uint64 : 0 to 18446744073709551615
// int8   : -128 to 127
// int16  : -32768 to 32767
// int32  : -2147483648 to 2147483647
// int64  : -9223372036854775808 to 9223372036854775807

//|------------head-------------|-----body-------|
//|---4 bytes---|----4 bytes----|-----dataLen----|
//|----------------------------------------------|
//|---dataLen---|-----msgId-----|-------Msg------|
//|----------------------------------------------|

type Msg struct {
	ID      uint32 // 消息ID
	DataLen uint32 // 消息内容的长度
	Data    []byte // 消息内容
}

// NewMsg 创建一个Message对象
func NewMsg(msgId uint32, data []byte) *Msg {
	m := Msg{}
	m.SetMsgId(msgId)
	m.SetData(data)
	return &m
}

// GetMsgId 获取消息名称
func (m *Msg) GetMsgId() uint32 {
	return m.ID
}

// SetMsgId 设置消息名称
func (m *Msg) SetMsgId(msgId uint32) {
	m.ID = msgId
}

// GetData 获取消息内容
func (m *Msg) GetData() []byte {
	return m.Data
}

// SetData 设置消息内容
func (m *Msg) SetData(data []byte) {
	m.DataLen = uint32(len(data))
	m.Data = data

}

// WithData 携带消息内容
func (m *Msg) WithData(data []byte) *Msg {
	m.SetData(data)
	return m
}

// WithMsgId 携带消息名称
func (m *Msg) WithMsgId(msgId uint32) *Msg {
	m.SetMsgId(msgId)
	return m
}

// GetDataLen 获取消息内容段长度
func (m *Msg) GetDataLen() uint32 {
	return m.DataLen
}

// SetDataLen 设置消息数据段长度
func (m *Msg) SetDataLen(len uint32) {
	m.DataLen = len
}
