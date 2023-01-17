package log

import (
	"fmt"
	"github.com/gogf/gf/os/glog"
	"strings"
)

// Info logger
func Info(tid string, f string, args ...interface{}) {
	setFlags(STD)
	str := formatStr(f, args...)
	glog.Info(baseInfo(tid) + str)
}

// Debug logger
func Debug(tid string, f string, args ...interface{}) {
	setFlags(STD)
	str := formatStr(f, args...)
	glog.Debug(baseInfo(tid) + str)
}

// Warn logger
func Warn(tid string, f string, args ...interface{}) {
	setFlags(STD)
	str := formatStr(f, args...)
	glog.Warning(baseInfo(tid) + str)
}

// Error logger
func Error(tid string, f string, args ...interface{}) {
	setFlags(STD)
	str := formatStr(f, args...)
	glog.Error(baseInfo(tid) + str)
}

// Assert error
func Assert(err error) {
	if err != nil {
		panic(err)
	}
}

func setFlags(flag LOGFLAG) {
	glog.SetFlags(LOGFLAGMAP[flag])
}

func formatStr(f string, args ...interface{}) string {
	return fmt.Sprintf(strings.Replace(f, "\n", "\n\t", -1), args...)
}

func baseInfo(traceID string) (base string) {
	base = "[" + serviceName + "]"
	if len(traceID) > 0 {
		base = base + " [" + traceID + "] "
	}
	return
}
