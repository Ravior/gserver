package ghttp

import (
	"net/http"
)

var (
	// defaultAllowHeaders is the default allowed headers for CORS.
	// It's defined another map for better header key searching performance.
	defaultSupportedHttpMethods = "GET,PUT,POST,DELETE,PATCH,HEAD,CONNECT,OPTIONS,TRACE"
	defaultAllowHeaders         = "Origin,Content-Type,Accept,User-Agent,Cookie,Authorization,X-Auth-Token,X-Requested-With"
)

func CORSDefault(header http.Header) http.Header {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Credentials", "true")
	header.Set("Access-Control-Allow-Methods", defaultSupportedHttpMethods)
	header.Set("Access-Control-Allow-Headers", defaultAllowHeaders)
	return header
}
