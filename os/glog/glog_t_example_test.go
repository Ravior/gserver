package glog

func Example_debug() {
	Debug("gserver")
	// May Output:
	// {"level":"debug","ts":"2022-08-01T16:24:07.151+0800","caller":"glog/glog_t_test.go:6","msg":"gserver"}
}

func Example_debugf() {
	Debugf("gserver, server:%d", 10001)
	// May Output:
	// {"level":"debug","ts":"2022-08-01T16:24:07.151+0800","caller":"glog/glog_t_test.go:6","msg":"gserver, server:10001"}
}
