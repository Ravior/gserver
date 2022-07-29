package ghttp

import "sync"

type HandlerFunc func(req *Request)
type MiddlewareFunc func(HandlerFunc) HandlerFunc

// 路由组
type Group struct {
	sync.Mutex
	hMap       map[string]HandlerFunc
	hMapMidd   map[string][]MiddlewareFunc
	middleware []MiddlewareFunc
}

// AddRoute 添加路由
func (g *Group) AddRoute(name string, handlerFunc HandlerFunc, middleware ...MiddlewareFunc) {
	g.Lock()
	defer g.Unlock()

	g.hMap[name] = handlerFunc
	g.hMapMidd[name] = middleware
}

func (g *Group) Use(middleware ...MiddlewareFunc) *Group {
	g.middleware = append(g.middleware, middleware...)
	return g
}

func (g *Group) applyMiddleware(name string) HandlerFunc {
	h, ok := g.hMap[name]
	if ok == false {
		return nil
	}

	for i := len(g.middleware) - 1; i >= 0; i-- {
		h = g.middleware[i](h)
	}

	for i := len(g.hMapMidd[name]) - 1; i >= 0; i-- {
		h = g.hMapMidd[name][i](h)
	}

	return h
}

func (g *Group) exec(name string, req *Request) {
	h := g.applyMiddleware(name)
	if h == nil {
		req.Response.Write("404")
		req.Response.Flush()
	} else {
		h(req)
	}

}

func NewRouter() *Router {
	return &Router{
		group: &Group{
			hMap:     make(map[string]HandlerFunc),
			hMapMidd: make(map[string][]MiddlewareFunc),
		},
	}
}

// 路由器
type Router struct {
	group *Group
}

func (r *Router) Run(req *Request) {
	path := req.URL.Path

	if r.group != nil {
		r.group.exec(path, req)
	}
}
