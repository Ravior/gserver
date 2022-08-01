package glog

import "testing"

func Test_Debug(t *testing.T) {
	Debug("gserver")
	Debugf("gserver, serverId:%d", 10001)
}

func Test_Info(t *testing.T) {
	Info("gserver")
	Infof("gserver, serverId:%d", 10001)
}

func Test_Warn(t *testing.T) {
	Warn("gserver")
	Warnf("gserver, serverId:%d", 10001)
}

func Test_Error(t *testing.T) {
	Error("gserver")
	Errorf("gserver, serverId:%d", 10001)
}
