package gnet

import (
	"fmt"
	"github.com/Ravior/gserver/os/glog"
)

var (
	defaultWorkerPoolSize uint32 = 8
	defaultWorkerTaskSize uint32 = 1024
)

// MsgHandler 消息处理器模块，msgHandler会有多个worker回来同时处理消息，conn发送的消息通过connId取模落入到某一个worker处理
type MsgHandler struct {
	WorkerPoolSize uint32          // 业务工作Worker池的数量
	WorkerTaskSize uint32          // 每个Worker的可等待执行Task数量
	TaskQueue      []chan *Request // Worker负责取任务的消息队列
	TaskExit       []chan bool
	Router         *Router // 路由
}

func NewMsgHandler(workerPoolSize uint32, workerTaskSize uint32) *MsgHandler {
	// 如果传入的值为0，则采用默认值
	if workerPoolSize == 0 {
		workerPoolSize = defaultWorkerPoolSize
	}
	if workerTaskSize == 0 {
		workerTaskSize = defaultWorkerTaskSize
	}

	return &MsgHandler{
		WorkerPoolSize: workerPoolSize,
		WorkerTaskSize: workerTaskSize,
		// 一个worker对应一个queue
		TaskQueue: make([]chan *Request, workerPoolSize),
		TaskExit:  make([]chan bool, workerPoolSize),
	}
}

// HandleMsg 已非阻塞方式处理消息
func (mh *MsgHandler) HandleMsg(req *Request) {
	defer func() {
		if err := recover(); err != nil {
			message := fmt.Sprintf("%s", err)
			glog.Error(message)
		}
	}()

	if mh.Router != nil {
		mh.Router.Run(req)
	}
}

// SetRouter 为消息添加具体的处理逻辑(路由)
func (mh *MsgHandler) SetRouter(router *Router) {
	mh.Router = router
}

// StartWorkerPool 启动worker工作池
func (mh *MsgHandler) StartWorkerPool() {
	glog.Debug("StartWork Worker Pool, Worker Num:", mh.WorkerPoolSize)
	// 遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskQueue[i] = make(chan *Request, mh.WorkerTaskSize)
		mh.TaskExit[i] = make(chan bool, 1)

		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.startOneWorker(i, mh.TaskQueue[i], mh.TaskExit[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandler) startOneWorker(workID int, taskQueue chan *Request, taskExit chan bool) {
	// 不断循环等待队列中的消息
	for {
		select {
		case request := <-taskQueue:
			mh.HandleMsg(request)
		case isExit := <-taskExit:
			if isExit {
				glog.Debugf("Worker ID: %d Exit", workID)
				return
			}
		}
	}
}

// StopWorkerPool 关闭Work线程池
func (mh *MsgHandler) StopWorkerPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		mh.TaskExit[i] <- true
		close(mh.TaskExit[i])
	}
}

// SendMsgToTaskQueue 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request *Request) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	mh.TaskQueue[workerID] <- request
}
