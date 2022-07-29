package ghttp

// GetQuery retrieves and returns parameter with given name <key> from query string
// and request body. It returns <def> if <key> does not exist in the query and <def> is given,
// or else it returns nil.
//
// Note that if there're multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: query > body.
func (r *Request) GetQuery(key string, def ...interface{}) interface{} {
	r.parseQuery()
	if len(r.queryMap) > 0 {
		if v, ok := r.queryMap[key]; ok {
			return v
		}
	}
	if r.Method == "GET" {
		r.parseBody()
	}
	if len(r.bodyMap) > 0 {
		if v, ok := r.bodyMap[key]; ok {
			return v
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}
