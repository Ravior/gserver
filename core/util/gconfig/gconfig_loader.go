package gconfig

import (
	"errors"
	"fmt"
	"github.com/Ravior/gserver/core/os/gfile"
	jsonLib "github.com/json-iterator/go"
	"io/ioutil"
)

// LoadJsonConfigFromBasePath 读取配置文件
func LoadJsonConfigFromBasePath(configFile string, data interface{}) error {
	configPath := fmt.Sprintf("%s%s", Global.BasePath, configFile)
	if ok := gfile.Exists(configPath); ok != true {
		text := fmt.Sprintf("Config File %s is not exist!!", configPath)
		return errors.New(text)
	}

	_data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err
	}
	// 将json数据解析到struct中
	err = jsonLib.Unmarshal(_data, data)
	if err != nil {
		return err
	}

	return nil
}

// LoadJsonConfig 读取配置文件
func LoadJsonConfig(configFile string, data interface{}) error {
	if ok := gfile.Exists(configFile); ok != true {
		text := fmt.Sprintf("Config File %s is not exist!!", configFile)
		return errors.New(text)
	}

	_data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	// 将json数据解析到struct中
	err = jsonLib.Unmarshal(_data, data)
	if err != nil {
		return err
	}

	return nil
}
