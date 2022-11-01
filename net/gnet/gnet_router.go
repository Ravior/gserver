package gnet

import (
	"errors"
	"fmt"
	"github.com/Ravior/gserver/crypto/gcrc32"
	"github.com/Ravior/gserver/os/glog"
	"github.com/Ravior/gserver/util/gserialize"
	"github.com/Ravior/gserver/util/gutil"
	"github.com/golang/protobuf/proto"
	"reflect"
	"strings"
	"sync"
)

type HandlerCallback interface{}
type HandlerFunc func(request *Request, msg proto.Message)
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// 路由组
type Group struct {
	prefix     string
	hMap       map[string]HandlerFunc
	hMapMidd   map[string][]MiddlewareFunc
	middleware []MiddlewareFunc
}

// AddRoute 添加路由
func (g *Group) AddRoute(path string, callback HandlerCallback, middleware ...MiddlewareFunc) {
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
	g.hMapMidd[path] = middleware
}

// Use 全局中间件
func (g *Group) Use(middleware ...MiddlewareFunc) *Group {
	g.middleware = append(g.middleware, middleware...)
	return g
}

func (g *Group) applyMiddleware(name string) HandlerFunc {

	h, ok := g.hMap[name]
	if ok == false {
		// 通配符
		h, ok = g.hMap["*"]
	}

	if ok {
		for i := len(g.middleware) - 1; i >= 0; i-- {
			h = g.middleware[i](h)
		}

		for i := len(g.hMapMidd[name]) - 1; i >= 0; i-- {
			h = g.hMapMidd[name][i](h)
		}
	}

	return h
}

func (g *Group) exec(name string, req *Request) {
	msgType, ok := RouteItemMgr.msgIdMap[req.GetMessage().GetMsgId()]
	if !ok {
		glog.Errorf("RouteItemMgr.msgIdMap Not Found MsgId: %d", req.GetMessage().GetMsgId())
		return
	}

	msg := reflect.New(msgType.Elem()).Interface().(proto.Message)
	err := gserialize.Protobuf.Unmarshal(req.GetMessage().GetData(), msg)
	if err != nil {
		glog.Errorf("unmarshal message error: %v", err)
		return
	}
	gutil.NiceCallFunc(func() {
		h := g.applyMiddleware(name)
		if h == nil {
			glog.Debug("Router Msg Handler Miss, MsgId:", req.GetMessage().GetMsgId())
		} else {
			defer func() {
				if err := recover(); err != nil {
					e := fmt.Sprintf("%v", err)
					glog.Errorf("handler msg has err:%v", e)
				}
			}()
			h(req, msg)
		}
	})
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

	msgId := gcrc32.Encrypt(proto)
	router := &RouteItem{
		MsgId: msgId,
		Proto: proto,
		Route: route,
	}
	if !gutil.IsEmpty(proto) {
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
		prefix:   prefix,
		hMap:     make(map[string]HandlerFunc),
		hMapMidd: make(map[string][]MiddlewareFunc),
	}

	r.groups = append(r.groups, g)
	return g
}

func (r *Router) Run(req *Request) {
	msgId := req.GetMessage().GetMsgId()
	route := RouteItemMgr.GetRoute(msgId)
	if gutil.IsEmpty(route) {
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
