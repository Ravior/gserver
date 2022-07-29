package ghttp

// GetQueryParams return all parameters by GET Http Method
func (r *Request) GetQueryParams() interface{} {
	return r.queryMap
}

// GetParam returns custom parameter with given name <key>.
// It returns <def> if <key> does not exist.
// It returns nil if <def> is not passed.
func (r *Request) GetParam(key string, def ...interface{}) interface{} {
	if r.paramsMap != nil {
		return r.paramsMap[key]
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// GetString is an alias and convenient function for GetRequestString.
// See GetRequestString.
func (r *Request) GetString(key string, def ...interface{}) string {
	return r.GetRequestString(key, def...)
}

// GetBool is an alias and convenient function for GetRequestBool.
// See GetRequestBool.
func (r *Request) GetBool(key string, def ...interface{}) bool {
	return r.GetRequestBool(key, def...)
}

// GetInt is an alias and convenient function for GetRequestInt.
// See GetRequestInt.
func (r *Request) GetInt(key string, def ...interface{}) int {
	return r.GetRequestInt(key, def...)
}

// GetInt32 is an alias and convenient function for GetRequestInt32.
// See GetRequestInt32.
func (r *Request) GetInt32(key string, def ...interface{}) int32 {
	return r.GetRequestInt32(key, def...)
}

// GetInt64 is an alias and convenient function for GetRequestInt64.
// See GetRequestInt64.
func (r *Request) GetInt64(key string, def ...interface{}) int64 {
	return r.GetRequestInt64(key, def...)
}

// GetUint is an alias and convenient function for GetRequestUint.
// See GetRequestUint.
func (r *Request) GetUint(key string, def ...interface{}) uint {
	return r.GetRequestUint(key, def...)
}

// GetUint8 is an alias and convenient function for GetRequestUint32.
// See GetRequestUint32.
func (r *Request) GetUint8(key string, def ...interface{}) uint8 {
	return r.GetRequestUint8(key, def...)
}

// GetUint32 is an alias and convenient function for GetRequestUint32.
// See GetRequestUint32.
func (r *Request) GetUint32(key string, def ...interface{}) uint32 {
	return r.GetRequestUint32(key, def...)
}

// GetUint64 is an alias and convenient function for GetRequestUint64.
// See GetRequestUint64.
func (r *Request) GetUint64(key string, def ...interface{}) uint64 {
	return r.GetRequestUint64(key, def...)
}

// GetFloat32 is an alias and convenient function for GetRequestFloat32.
// See GetRequestFloat32.
func (r *Request) GetFloat32(key string, def ...interface{}) float32 {
	return r.GetRequestFloat32(key, def...)
}

// GetFloat64 is an alias and convenient function for GetRequestFloat64.
// See GetRequestFloat64.
func (r *Request) GetFloat64(key string, def ...interface{}) float64 {
	return r.GetRequestFloat64(key, def...)
}
