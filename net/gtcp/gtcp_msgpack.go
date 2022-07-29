package gtcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/Ravior/gserver/net/gnet"
)

//|------------head-------------|-----body-------|
//|---4 bytes---|----4 bytes----|-----dataLen----|
//|----------------------------------------------|
//|---dataLen---|-----msgId-----|-------body-----|
//|----------------------------------------------|

var (
	defaultHeaderLen     uint32 = 8 // 不可修改
	defaultMaxPacketSize uint32 = 2048
)

var MsgPack = NewMsgPack()

// 数据包 拆包/封包
type msgPack struct {
	maxPacketSize uint32
	littleEndian  bool
}

func NewMsgPack() *msgPack {
	return &msgPack{
		maxPacketSize: defaultMaxPacketSize,
	}
}

// GetHeadLen 获取包头长度方法
func (mp *msgPack) GetHeadLen() uint32 {
	// dataLen uint32(4字节) + msgId uint16(2字节)
	return defaultHeaderLen
}

// SetByteOrder 设置大小端
func (mp *msgPack) SetByteOrder(littleEndian bool) {
	mp.littleEndian = littleEndian
}

// SetMaxPacketSize 设置最大包体长度
func (mp *msgPack) SetMaxPacketSize(maxPacketSize uint32) {
	if maxPacketSize > 0 {
		mp.maxPacketSize = maxPacketSize
	}
}

// Pack 封包方法(压缩数据)
func (mp *msgPack) Pack(msg *gnet.Msg) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	if mp.littleEndian == true {
		// 小端
		// 写dataLen
		if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
			return nil, err
		}

		// 写msgId
		if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
			return nil, err
		}

		// 写data数据
		if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
			return nil, err
		}
	} else {
		// 写dataLen
		if err := binary.Write(dataBuff, binary.BigEndian, msg.GetDataLen()); err != nil {
			return nil, err
		}

		// 写msgId
		if err := binary.Write(dataBuff, binary.BigEndian, msg.GetMsgId()); err != nil {
			return nil, err
		}

		// 写data数据
		if err := binary.Write(dataBuff, binary.BigEndian, msg.GetData()); err != nil {
			return nil, err
		}
	}

	return dataBuff.Bytes(), nil
}

// Unpack 拆包方法(解压数据)
func (mp *msgPack) Unpack(headData []byte) (*gnet.Msg, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(headData)

	// 只解压head的信息，得到dataLen和msgID
	msg := &gnet.Msg{}

	if mp.littleEndian == true {
		// 读dataLen
		if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
			return nil, err
		}

		// 读msgId
		if err := binary.Read(dataBuff, binary.LittleEndian, &msg.ID); err != nil {
			return nil, err
		}
	} else {
		// 大端
		// 读dataLen
		if err := binary.Read(dataBuff, binary.BigEndian, &msg.DataLen); err != nil {
			return nil, err
		}

		// 读msgId
		if err := binary.Read(dataBuff, binary.BigEndian, &msg.ID); err != nil {
			return nil, err
		}
	}

	var dataLen = msg.GetDataLen() + mp.GetHeadLen()

	// 判断dataLen的长度是否超出我们允许的最大包长度
	if mp.maxPacketSize > 0 && dataLen > mp.maxPacketSize {
		return nil, errors.New("too large proto CaveList received")
	}

	// 这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil
}
