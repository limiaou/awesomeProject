package rdb

import (
	_ "github.com/gogf/gf/contrib/drivers/clickhouse/v2"
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/database/gdb"
)

func SetDb() {
	gdb.SetConfig(gdb.Config{
		"mysql": gdb.ConfigGroup{
			gdb.ConfigNode{
				Type:  "mysql",
				Host:  "10.255.202.20",
				Port:  "3306",
				User:  "root",
				Pass:  "Connext@1qaz@WSX",
				Name:  "db_bill",
				Debug: false,
			},
		},
		"clickhouse": gdb.ConfigGroup{
			//gdb.ConfigNode{
			//	Type:  "clickhouse",
			//	Host:  "101.132.45.204",
			//	Port:  "9000",
			//	User:  "default",
			//	Pass:  "12345678",
			//	Name:  "mdlz",
			//	Debug: false,
			//},
			gdb.ConfigNode{
				Type:  "clickhouse",
				Host:  "10.255.202.73",
				Port:  "9000",
				User:  "default",
				Pass:  "12345678",
				Name:  "finance",
				Debug: false,
			},
			//gdb.ConfigNode{
			//	Type:  "clickhouse",
			//	Host:  "10.255.201.106",
			//	Port:  "9000",
			//	User:  "default",
			//	Pass:  "12345678",
			//	Name:  "mdlz",
			//	Debug: false,
			//},
		},
	})
}
