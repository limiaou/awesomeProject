package log

import "github.com/gogf/gf/os/glog"

type LOGFLAG int

const (
	// LONG Print full file name and line number: /a/b/c/d.go:23.
	LONG LOGFLAG = 2 << iota
	// SHORT Print final file name element and line number: d.go:23. overrides F_FILE_LONG.
	SHORT
	// STD Print the date in the local time zone: 2009-01-23  01:23:23
	STD
)

var serviceName string = "Unkown_Service"

var LOGFLAGMAP = map[LOGFLAG]int{
	LONG:        glog.F_FILE_LONG,
	SHORT:       glog.F_FILE_SHORT,
	STD:         glog.F_TIME_STD,
	SHORT | STD: glog.F_FILE_SHORT | glog.F_TIME_STD,
	LONG | STD:  glog.F_FILE_LONG | glog.F_TIME_STD,
}
