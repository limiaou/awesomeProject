package rdb

import (
	"awesomeProject/log"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"testing"
	"time"
)

func TestFindDateDataByPositive(t *testing.T) {
	SetDb()
	clickhouse, err := gdb.NewByGroup("mysql")
	if err != nil {
		panic(err)
	}
	parse, err := time.Parse("2006-01-02", "2022-12-13")
	positiveData, err := FindDateDataByPositive("t_inf_bill_aliyun_summary_split_positive_day", parse, clickhouse)
	for _, v := range positiveData {
		fmt.Println(v)
	}
}

func TestUpdatePositiveDay(t *testing.T) {
	SetDb()
	mysql, err := gdb.NewByGroup("mysql")
	if err != nil {
		log.Error("", "%v", err)
	}
	clickhouse, err := gdb.NewByGroup("clickhouse")
	if err != nil {
		log.Error("", "%v", err)
	}
	//err = rdb.UpdateNegativeDay(mysql, clickhouse)
	parse, err := time.Parse("2006-01-02", "2022-12-01")
	if err != nil {
		panic(err)
	}
	for err == nil {
		curMonth := parse.Format("2006-01")
		err = SavePositiveDay(curMonth, mysql, clickhouse)
		parse = parse.AddDate(0, 1, 0)
	}
}

func TestQueryUpdatePositive(t *testing.T) {
	SetDb()
	clickhouse, err := gdb.NewByGroup("clickhouse")
	if err != nil {
		panic(err)
	}
	parse, err := time.Parse("2006-01-02", "2022-12-13")
	curMonth := parse.Format("2006-01")
	positiveData, err := QueryUpdatePositive(curMonth, clickhouse)
	for _, v := range positiveData {
		fmt.Println(v)
	}
}
