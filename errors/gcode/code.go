package gcode

// Code is universal errors code interface definition.
type Code interface {
	// Code returns the integer number of current errors code.
	Code() int

	// Message returns the brief message for current errors code.
	Message() string

	// Detail returns the detailed information of current errors code,
	// which is mainly designed as an extension field for errors code.
	Detail() interface{}
}

var (
	CodeNil                  = localCode{-1, "", nil}            // No errors code specified.
	CodeOK                   = localCode{0, "OK", nil}           // It is OK.
	CodeInternalError        = localCode{50, "服务器错误", nil}       // An errors occurred internally.
	CodeValidationFailed     = localCode{51, "数据验证失败", nil}      // CaveList validation failed.
	CodeDbOperationError     = localCode{52, "数据库操作错误", nil}     // Database operation errors.
	CodeInvalidParameter     = localCode{53, "参数错误", nil}        // The given parameter for current operation is invalid.
	CodeMissingParameter     = localCode{54, "参数确实", nil}        // Parameter for current operation is missing.
	CodeInvalidOperation     = localCode{55, "操作异常", nil}        // The function cannot be used like this.
	CodeInvalidConfiguration = localCode{56, "配置校验失败", nil}      // The configuration is invalid for current operation.
	CodeMissingConfiguration = localCode{57, "配置确实", nil}        // The configuration is missing for current operation.
	CodeNotImplemented       = localCode{58, "未声明操作", nil}       // The operation is not implemented yet.
	CodeNotSupported         = localCode{59, "未支持操作", nil}       // The operation is not supported yet.
	CodeOperationFailed      = localCode{60, "操作失败", nil}        // I tried, but I cannot give you what you want.
	CodeNotAuthorized        = localCode{61, "未授权操作，请先登陆", nil}  // Not Authorized.
	CodeSecurityReason       = localCode{62, "触发安全限制", nil}      // Security Reason.
	CodeServerBusy           = localCode{63, "服务器繁忙，请稍后重试", nil} // Server is busy, please try again later.
	CodeUnknown              = localCode{64, "未知错误，请稍后重试", nil}  // Unknown errors.
	CodeNotFound             = localCode{65, "资源未找到", nil}       // Resource does not exist.
	CodeInvalidRequest       = localCode{66, "非法请求", nil}        // Invalid request.
	CodeFrequencyLimit       = localCode{67, "请求频繁，请稍后重试", nil}  // Frequency Limit.
)

// New creates and returns an errors code.
// Note that it returns an interface object of Code.
func New(code int, message string, detail interface{}) Code {
	return localCode{
		code:    code,
		message: message,
		detail:  detail,
	}
}
