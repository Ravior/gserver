package gnsq

import (
	"github.com/Ravior/gserver/core/os/glog"
	"github.com/Ravior/gserver/core/os/gtimer"
	"github.com/Ravior/gserver/core/util/gcontainer/gmap"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

// NsqCli GMQ主要用于消息队列处理
// 目前主要自用的消息队列为Nsq
type NsqCli struct {
	addr     string
	producer *nsq.Producer
	consumer *nsq.Consumer
	lock     sync.Mutex
	buffs    *gmap.StrAnyMap // 客户端消息缓存
	buffSize int             // 客户端消息缓存大小(默认100)
}

func NewNsqCli(addr string, buffSize ...int) *NsqCli {
	_buffSize := 100
	if len(buffSize) > 0 {
		_buffSize = buffSize[0]
	}
	cli := &NsqCli{
		addr:     addr,
		buffs:    gmap.NewStrAnyMap(true),
		buffSize: _buffSize,
	}
	cli.init()

	return cli
}

func (n *NsqCli) initProducer() error {
	if n.producer == nil {
		n.lock.Lock()
		defer n.lock.Unlock()

		config := nsq.NewConfig()
		producer, err := nsq.NewProducer(n.addr, config)
		if err != nil {
			glog.Errorf("create producer failed, err:%v", err)
			return err
		}
		n.producer = producer
	}
	return nil
}

type NsqMessageHandler interface {
	HandleMessage(msg *nsq.Message) error
	GetName() string
}

// myNsqMsgHandler 是一个消费者类型
type myNsqMsgHandler struct {
	Title   string
	Handler NsqMessageHandler
}

// HandleMessage 是需要实现的处理消息的方法
func (m *myNsqMsgHandler) HandleMessage(msg *nsq.Message) error {
	if m.Handler != nil {
		err := m.Handler.HandleMessage(msg)
		if err != nil {
			glog.Errorf("NSQCli处理消息出现错误，err:%v", err)
			return err
		}
	}
	return nil
}

func (n *NsqCli) init() {
	if n.buffSize > 0 {
		// 定时清理缓存
		gtimer.Add(60*time.Second, func() {
			n.sync()
		})
	}
}

func (n *NsqCli) Stop() {
	n.producer.Stop()
	n.consumer.Stop()
	if n.buffSize > 0 {
		n.sync()
	}
}

func (n *NsqCli) initConsumer(topic string, channel string, handlerName string, handler NsqMessageHandler) error {
	config := nsq.NewConfig()
	config.LookupdPollInterval = 15 * time.Second
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		glog.Errorf("NSQCli初始化消费者失败 err:%v", err)
		return err
	}
	msgHandler := &myNsqMsgHandler{
		Title:   handlerName,
		Handler: handler,
	}
	c.AddHandler(msgHandler)

	// 直接链接NsqD
	if err := c.ConnectToNSQD(n.addr); err != nil {
		return err
	}

	n.consumer = c
	return nil
}

func (n *NsqCli) Publish(topic string, body []byte) error {
	if err := n.initProducer(); err != nil {
		glog.Errorf("NSQCli初始化生产者失败，err:%v", err)
		return err
	}

	buffs := make([][]byte, 0)
	if _buffs := n.buffs.Get(topic); _buffs != nil {
		if _, ok := _buffs.([][]byte); ok {
			buffs = _buffs.([][]byte)
		}
	}

	buffs = append(buffs, body)
	msgCount := len(buffs)
	if msgCount >= n.buffSize {
		responseChan := make(chan *nsq.ProducerTransaction, len(buffs))
		err := n.producer.MultiPublishAsync(topic, buffs, responseChan)
		if err != nil {
			glog.Errorf("NSQCli创建投递任务失败, topic:%s, err: %v", topic, err)
			return err
		}
		trans := <-responseChan
		if trans.Error != nil {
			glog.Errorf("NSQCli投递消息失败, topic:%s, err: %v", topic, trans.Error.Error())
		}
		buffs = make([][]byte, 0)
	}
	n.buffs.Set(topic, buffs)

	return nil
}

func (n *NsqCli) sync() {
	syncTopics := make([]string, 0)
	n.buffs.Iterator(func(topic string, _buffs interface{}) bool {
		if buffs, ok := _buffs.([][]byte); ok {
			if len(buffs) > 0 {
				syncTopics = append(syncTopics, topic)
			}
		}
		return true
	})
	for _, topic := range syncTopics {
		if _buffs := n.buffs.Get(topic); _buffs != nil {
			if buffs, ok := _buffs.([][]byte); ok {
				responseChan := make(chan *nsq.ProducerTransaction, len(buffs))
				err := n.producer.MultiPublishAsync(topic, buffs, responseChan)
				if err != nil {
					glog.Errorf("NSQCli创建投递任务失败, topic:%s, err: %v", topic, err)
				}
				trans := <-responseChan
				if trans.Error != nil {
					glog.Errorf("NSQCli投递消息失败, topic:%s, err: %v", topic, trans.Error.Error())
				}
				buffs = make([][]byte, 0)
				n.buffs.Set(topic, buffs)
			}
		}
	}
}

func (n *NsqCli) Consume(topic string, channel string, handlerName string, handler NsqMessageHandler) {
	err := n.initConsumer(topic, channel, handlerName, handler)
	if err != nil {
		glog.Errorf("NSQCli初始化消费者失败，err:%v", err)
		return
	}
}
