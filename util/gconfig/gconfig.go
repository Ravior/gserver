package gconfig

type LogConfig struct {
	Level           string // 输出日志等级
	StackLevel      string // 堆栈输出日志等级
	EnableWriteFile bool   // 是否输出文件(必需配置FilePath)
	EnableConsole   bool   // 是否控制台输出
	FilePath        string // 日志文件输出路径
	FileFormat      string // 日志文件格式
	MaxAge          int    // 最大保留天数  maxAge达到限制，则会被清理
	RotationTime    int    // 日志自动切割时长，单位小时
	TimeFormat      string // 时间输出格式
	PrintCaller     bool   // 是否打印调用函数
}

type addr struct {
	IP   string
	Port int32
}

// TcpServerConfig Tcp服务器配置
type TcpServerConfig struct {
	addr
	MaxConn        int32  // 当前服务器允许的最大链接数
	WorkerPoolSize uint32 // 业务工作Worker池的数量
	WorkerTaskLen  uint32 // 业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen  uint32 // BuffMsg长度
}

// WsServerConfig Websocket服务器配置
type WsServerConfig struct {
	addr
	MaxConn        int32  // 当前服务器允许的最大链接数
	WorkerPoolSize uint32 // 业务工作Worker池的数量
	WorkerTaskLen  uint32 // 业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen  uint32 // BuffMsg长度

	CertFile string // SSL证书地址
	KeyFile  string // SSL证书密钥地址
}

// HttpServerConfig Http服务器配置
type HttpServerConfig struct {
	addr
	CertFile string // SSL证书地址
	KeyFile  string // SSL证书密钥地址
}

// RpcServerConfig RPC服务器配置
type RpcServerConfig struct {
	addr
}

// DBConfig 数据库配置
type DBConfig struct {
	Host        string
	Port        int
	User        string
	Password    string
	DbName      string
	Prefix      string
	MaxIdleConn int  // 最大空闲连接数
	MaxOpenConn int  // 最大打开连接数
	ShowLog     bool // 是否现实日志
}

type PprofConfig struct {
	addr
}

// ClientConfig 客户端配置
type ClientConfig struct {
	RemoteHost    string // 远程服务器主机ip
	RemoteTcpPort int32  // 远程服务器主机端口号
	ClientName    string
	ClientId      string
}

type Config struct {
	Debug      bool   // 是否开启Debug模式
	ServerCode string // 服务器编号
	BasePath   string // 程序根目录
	ServerName string // 当前服务器名称
	ServerId   string // 服务器id
	Version    string // 服务器版本

	Master   ClientConfig          // 中心服务器链接地址配置
	DataBase map[string]*DBConfig  // 数据库配置
	Log      map[string]*LogConfig // 日志配置

	TcpServer  TcpServerConfig  // TCP服务器配置
	WsServer   WsServerConfig   // WebSocket服务器配置
	HttpServer HttpServerConfig // HTTP服务器配置
	RpcServer  RpcServerConfig  // RPC服务器配置
	Pprof      PprofConfig      // Pprof控制
}

func NewConfig() Config {
	c := Config{}
	c.ServerName = "Default Server"
	c.ServerId = "Server"
	c.Version = "1.0.0"

	c.Log = make(map[string]*LogConfig)
	return c
}
