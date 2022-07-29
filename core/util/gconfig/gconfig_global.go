package gconfig

// 存储一切全局参数，供其他模块使用
type globalConfig struct {
	Config
}

var Global *globalConfig

/*
	提供init方法，默认加载
*/
func init() {
	// 初始化GlobalObject变量，设置一些默认值
	Global = &globalConfig{
		Config: NewConfig(),
	}
}

// Load 读取配置文件
func (g *globalConfig) Load(configFile string) error {
	return LoadJsonConfig(configFile, &g.Config)
}
