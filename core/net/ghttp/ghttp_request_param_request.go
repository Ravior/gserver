package ghttp

import (
	"github.com/Ravior/gserver/core/util/gcontainer/gvar"
)

// GetRequest retrieves and returns the parameter named <key> passed from client and
// custom params as interface{}, no matter what HTTP method the client is using. The
// parameter <def> specifies the default value if the <key> does not exist.
//
// GetRequest is one of the most commonly used functions for retrieving parameters.
//
// Note that if there're multiple parameters with the same name, the parameters are
// retrieved and overwrote in order of priority: router < query < body < form < custom.
func (r *Request) GetRequest(key string, def ...interface{}) interface{} {
	value := r.GetParam(key)
	if value == nil {
		value = r.GetForm(key)
	}
	if value == nil {
		r.parseBody()
		if len(r.bodyMap) > 0 {
			value = r.bodyMap[key]
		}
	}
	if value == nil {
		value = r.GetQuery(key)
	}
	if value != nil {
		return value
	}
	if len(def) > 0 {
		return def[0]
	}
	return value
}

// GetRequestVar retrieves and returns the parameter named <key> passed from client and
// custom params as gvar.Var, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetRequest(key, def...))
}

// GetRequestString retrieves and returns the parameter named <key> passed from client and
// custom params as string, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestString(key string, def ...interface{}) string {
	return r.GetRequestVar(key, def...).String()
}

// GetRequestBool retrieves and returns the parameter named <key> passed from client and
// custom params as bool, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestBool(key string, def ...interface{}) bool {
	return r.GetRequestVar(key, def...).Bool()
}

// GetRequestInt retrieves and returns the parameter named <key> passed from client and
// custom params as int, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInt(key string, def ...interface{}) int {
	return r.GetRequestVar(key, def...).Int()
}

// GetRequestInt32 retrieves and returns the parameter named <key> passed from client and
// custom params as int32, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInt32(key string, def ...interface{}) int32 {
	return r.GetRequestVar(key, def...).Int32()
}

// GetRequestInt64 retrieves and returns the parameter named <key> passed from client and
// custom params as int64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestInt64(key string, def ...interface{}) int64 {
	return r.GetRequestVar(key, def...).Int64()
}

// GetRequestUint retrieves and returns the parameter named <key> passed from client and
// custom params as uint, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint(key string, def ...interface{}) uint {
	return r.GetRequestVar(key, def...).Uint()
}

// GetRequestUint retrieves and returns the parameter named <key> passed from client and
// custom params as uint, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint8(key string, def ...interface{}) uint8 {
	return r.GetRequestVar(key, def...).Uint8()
}

// GetRequestUint32 retrieves and returns the parameter named <key> passed from client and
// custom params as uint32, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestVar(key, def...).Uint32()
}

// GetRequestUint64 retrieves and returns the parameter named <key> passed from client and
// custom params as uint64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestVar(key, def...).Uint64()
}

// GetRequestFloat32 retrieves and returns the parameter named <key> passed from client and
// custom params as float32, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestVar(key, def...).Float32()
}

// GetRequestFloat64 retrieves and returns the parameter named <key> passed from client and
// custom params as float64, no matter what HTTP method the client is using. The parameter
// <def> specifies the default value if the <key> does not exist.
func (r *Request) GetRequestFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestVar(key, def...).Float64()
}
