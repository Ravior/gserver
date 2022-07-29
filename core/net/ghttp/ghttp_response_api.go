package ghttp

import (
	"fmt"
	"github.com/Ravior/gserver/core/errors/gcode"
	"github.com/Ravior/gserver/core/errors/gerror"
	"github.com/Ravior/gserver/core/internal/json"
	"strings"
	"time"
)

type ApiJsonResult struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg,omitempty"`
	Data  interface{} `json:"data,omitempty"`
	MsgId string      `json:"msgId"`
	Time  int64       `json:"time"`
}

func (r *Response) ApiSuccess(data interface{}) {
	msgId := strings.Replace(r.Request.URL.Path, "/", "", 1)
	res := ApiJsonResult{
		Code:  0,
		Data:  data,
		Time:  time.Now().Unix(),
		MsgId: msgId,
	}
	if b, err := json.Marshal(res); err != nil {
		fmt.Println("JSON序列化失败", err)
	} else {
		CORSDefault(r.Header()).Set("Content-Type", "application/json")
		r.Write(b)
		r.Flush()
	}
}

func (r *Response) ApiFail(err error) {
	errCode := gerror.Code(err).Code()
	errMsg := err.Error()
	if errCode == gcode.CodeNil.Code() {
		errMsg = "服务器错误"
	}

	msgId := strings.Replace(r.Request.URL.Path, "/", "", 1)
	res := ApiJsonResult{
		Code:  errCode,
		Msg:   errMsg,
		Time:  time.Now().Unix(),
		MsgId: msgId,
	}

	if b, err := json.Marshal(res); err != nil {
		fmt.Println("JSON序列化失败", err)
	} else {
		CORSDefault(r.Header()).Set("Content-Type", "application/json")
		r.Write(b)
		r.Flush()
	}
}

func (r *Response) ApiRes(data []byte) {
	r.Write(data)
}
