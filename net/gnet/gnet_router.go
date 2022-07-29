package gnet

import (
	"errors"
	"fmt"
	"github.com/Ravior/gserver/crypto/gcrc32"
	"github.com/Ravior/gserver/internal/empty"
	"github.com/Ravior/gserver/os/glog"
	"github.com/Ravior/gserver/util/gconfig"
	"github.com/Ravior/gserver/util/gserialize"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MaxResTime = 10 // 最大响应时长,单位ms

type HandlerCallback interface{}
type HandlerFunc func(request *Request, msg proto.Message)

// Group 路由组
type Group struct {
	prefix string
	hMap   map[string]HandlerFunc
}

// AddRoute 添加路由
func (g *Group) AddRoute(path string, callback HandlerCallback) {
	err, funValue, msgType := checkMsgRouteCallback(callback)
	if err != nil {
		glog.Errorf("AddRoute Error: %v， %v", err, funValue)
		return
	}
	msg := reflect.New(msgType).Elem().Interface().(proto.Message)
	msgName := gserialize.Protobuf.GetMessageName(msg)
	route := fmt.Sprintf("%s.%s", g.prefix, path)
	RouteItemMgr.AddRoute(msgName, route, reflect.TypeOf(msg))

	g.hMap[path] = func(request *Request, msg2 proto.Message) {
		funValue.Call([]reflect.Value{reflect.ValueOf(request), reflect.ValueOf(msg2)})
	}
}

func (g *Group) getHandler(name string) HandlerFunc {
	h, ok := g.hMap[name]
	if ok == false {
		// 通配符
		h, ok = g.hMap["*"]
	}

	return h
}

func (g *Group) exec(name string, req *Request) {
	// 只有在调试环境下执行，避免空耗CPU
	if gconfig.Global.Debug {
		glog.Debugf("Handle Msg, Route:%s.%s, ConnId:%d, Addr:%s", g.prefix, name, req.GetConnId(), req.GetConnection().RemoteAddr())
	}

	msgType, ok := RouteItemMgr.msgIdMap[req.GetMessage().GetMsgId()]
	if !ok {
		glog.Errorf("RouteItemMgr.msgIdMap Not Found MsgId: %d", req.GetMessage().GetMsgId())
		return
	}

	h := g.getHandler(name)
	if h == nil {
		glog.Errorf("Router Msg Handler Miss, MsgId:%d", req.GetMessage().GetMsgId())
	} else {
		go func() {
			defer func() {
				if err := recover(); err != nil {
					glog.Errorf("Router Msg Handler Has Error, Route:%s.%s, ConnId:%d, Error:%v", g.prefix, name, req.GetConnId(), err)
				}
			}()

			msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
			err := gserialize.Protobuf.Unmarshal(req.GetMessage().GetData(), msg)
			if err != nil {
				glog.Errorf("Route:%s.%s, ConnId:%d，unmarshal message error: %v", g.prefix, name, req.GetConnId(), err)
				return
			}
			// 只有在调试环境下执行，避免空耗CPU
			if gconfig.Global.Debug {
				glog.Debugf("[ElapsedTime] Start. Route:%s.%s ｜ ConnId:%d | Param: [%+v]", g.prefix, name, req.GetConnId(), msg)
			}
			bt := time.Now().UnixNano()
			h(req, msg)
			et := time.Now().UnixNano()
			diff := (et - bt) / int64(time.Millisecond)

			if diff >= MaxResTime {
				// 超过10的日志记录处理
				glog.Warnf("[ElapsedTime] MaxResTime, End. Route:%s.%s | ConnId:%d | Cost: %dms", g.prefix, name, req.GetConnId(), diff)
			} else {
				// 只有在调试环境下执行，避免空耗CPU
				if gconfig.Global.Debug {
					glog.Debugf("[ElapsedTime] End. Route:%s.%s | ConnId:%d | Cost: %dms", g.prefix, name, req.GetConnId(), diff)
				}
			}
		}()
	}
}

var RouteItemMgr = &msgRouteMgr{
	routes:   make([]*RouteItem, 0),
	protoMap: make(map[string]uint32),
	msgIdMap: make(map[uint32]reflect.Type),
}

type RouteItem struct {
	MsgId uint32 `json:"msgId"`
	Proto string `json:"proto"`
	Route string `json:"route"`
}

type msgRouteMgr struct {
	sync.Mutex
	routes   []*RouteItem
	protoMap map[string]uint32
	msgIdMap map[uint32]reflect.Type
}

func (m *msgRouteMgr) Init() {
	m.reload()
}

func (m *msgRouteMgr) reload() {
	m.routes = m.routes[0:0]
}

func (m *msgRouteMgr) AddRoute(proto string, route string, msg reflect.Type) {
	m.Lock()
	defer m.Unlock()

	msgId := gcrc32.EncryptString(proto)
	router := &RouteItem{
		MsgId: msgId,
		Proto: proto,
		Route: route,
	}
	if !empty.IsEmpty(proto) {
		m.protoMap[proto] = msgId
		m.msgIdMap[msgId] = msg
	}

	m.routes = append(m.routes, router)
}

func (m *msgRouteMgr) GetRoute(msgId uint32) string {
	for _, route := range m.routes {
		if route.MsgId == msgId {
			return route.Route
		}
	}
	return ""
}

// Router 路由器
type Router struct {
	groups []*Group
}

func (r *Router) Group(prefix string) *Group {
	g := &Group{
		prefix: prefix,
		hMap:   make(map[string]HandlerFunc),
	}

	r.groups = append(r.groups, g)
	return g
}

func (r *Router) Run(req *Request) {
	msgId := req.GetMessage().GetMsgId()
	route := RouteItemMgr.GetRoute(msgId)
	if empty.IsEmpty(route) {
		return
	}

	prefix := ""
	msgName := route
	sArr := strings.Split(route, ".")
	if len(sArr) == 2 {
		prefix = sArr[0]
		msgName = sArr[1]
	}

	for _, g := range r.groups {
		if g.prefix == prefix {
			g.exec(msgName, req)
		} else if g.prefix == "*" {
			g.exec(msgName, req)
		}
	}
}

// 检查形如 func(arg0, proto.Message)
func checkMsgRouteCallback(cb interface{}) (err error, funValue reflect.Value, msgType reflect.Type) {
	cbType := reflect.TypeOf(cb)
	if cbType.Kind() != reflect.Func {
		err = errors.New("callback not a func")
		return
	}

	numArgs := cbType.NumIn()
	if numArgs != 2 {
		err = errors.New("callback param num must greater than 0")
		return
	}
	req := cbType.In(0)
	if req.Kind() != reflect.Ptr {
		err = errors.New("callback param args0 not ptr")
		return
	}

	msgType = cbType.In(1)
	if msgType.Kind() != reflect.Ptr {
		err = errors.New("callback param args1 not ptr")
		return
	}

	funValue = reflect.ValueOf(cb)

	return
}
