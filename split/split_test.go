package split

import (
	"awesomeProject/log"
	"awesomeProject/split/rdb"
	typeMNLZ "awesomeProject/split/type"
	"github.com/gogf/gf/v2/database/gdb"
	"testing"
	"time"
)

func TestSplit(t *testing.T) {
	rdb.SetDb()
	startDate, err := time.Parse("2006-01-02", "2022-01-01")
	log.Assert(err)
	now := time.Now()
	for startDate.Before(now) {
		err = SubMonthlyAndYearlyAmortization(startDate.Format("2006-01"), startDate.Format("2006-01-02"))
		log.Assert(err)
		startDate = startDate.AddDate(0, 0, 1)
	}
}

func TestUpdate(t *testing.T) {
	rdb.SetDb()
	mysql, err := gdb.NewByGroup("mysql")
	if err != nil {
		log.Error("", "%v", err)
	}
	clickhouse, err := gdb.NewByGroup("clickhouse")
	if err != nil {
		log.Error("", "%v", err)
	}
	err = rdb.TruncatePositiveDay(clickhouse)
	//err = rdb.UpdateNegativeDay(mysql, clickhouse)
	parse, err := time.Parse("2006-01-02", "2022-01-01")
	if err != nil {
		panic(err)
	}
	for err == nil {
		curMonth := parse.Format("2006-01")
		err = rdb.UpdatePositiveDay(curMonth, mysql, clickhouse)
		parse = parse.AddDate(0, 1, 0)
	}
}

func TestSaveOrUpdate(t *testing.T) {
	rdb.SetDb()
	orm, err := gdb.NewByGroup("clickhouse")
	if err != nil {
		log.Error("", "%v", err)
		//return err
	}
	//begin := returnError(orm.Begin(orm.Transaction())).(*gdb.TX)
	parse, err := time.Parse("2006-01-02", "2022-01-01")
	if err != nil {
		panic(err)
	}
	for parse.Before(time.Now()) {
		curMonth := parse.Format("2006-01")
		err = SplitAllShareCosts(curMonth, orm)
		saveOrUpdateData(curMonth, orm)
		parse = parse.AddDate(0, 1, 0)
	}
}

func TestSplitEveryDay(t *testing.T) {
	startTime := "2022-01-01"
	rdb.SetDb()
	mysql, err := gdb.NewByGroup("mysql")
	if err != nil {
		log.Error("", "%v", err)
	}
	clickhouse, err := gdb.NewByGroup("clickhouse")
	if err != nil {
		log.Error("", "%v", err)
	}
	startDate, err := time.Parse("2006-01-02", startTime)
	log.Assert(err)
	now := time.Now()
	for startDate.Before(now) {
		err = SubMonthlyAndYearlyAmortization(startDate.Format("2006-01"), startDate.Format("2006-01-02"))
		log.Assert(err)
		startDate = startDate.AddDate(0, 0, 1)
	}
	//mysql同步clickhouse
	err = rdb.UpdateNegativeDay(mysql, clickhouse)
	if err != nil {
		panic(err)
	}
	err = rdb.TruncatePositiveDay(clickhouse)
	if err != nil {
		panic(err)
	}
	parseU, err := time.Parse("2006-01-02", typeMNLZ.BEGINNING_OF_THE_YEAR)
	for err == nil {
		curMonth := parseU.Format("2006-01")
		err = rdb.UpdatePositiveDay(curMonth, mysql, clickhouse)
		parseU = parseU.AddDate(0, 1, 0)
	}
	//分摊
	parseS, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	for parseS.Before(time.Now()) {
		curMonth := parseS.Format("2006-01")
		err = SplitAllShareCosts(curMonth, clickhouse)
		log.Assert(err)
		saveOrUpdateData(curMonth, clickhouse)
		parseS = parseS.AddDate(0, 1, 0)
	}
}
