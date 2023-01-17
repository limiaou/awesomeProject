package rdb

import (
	"awesomeProject/log"
	typeMNLZ "awesomeProject/split/type"
	"errors"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/v2/database/gdb"
	"time"
)

func SaveSummaryDay(data []*typeMNLZ.AliyunBillSummaryDay, tableName string, tx gdb.DB) (int64, error) {
	_, err := tx.Model(tableName).Data(&data).Save()
	if err != nil {
		return 0, err
	}
	return int64(len(data)), err
}

func DeleteSummaryDay(date, tableName string, tx gdb.DB) (err error) {
	model := tx.Model(tableName)
	model.Where("billing_date =", date)
	_, err = model.Delete()
	return err
}

func QueryUpdatePositive(currentMonth string, tx gdb.DB) ([]*typeMNLZ.AliyunBillSummaryDay, error) {
	results := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := tx.Model("t_inf_bill_aliyun_summary_split_positive_day").
		Where("application!='ALL_SHARE'").
		Where("billing_cycle=?", currentMonth)
	err := model.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func DeleteSpiltShare(currentMonth string, tx gdb.DB) (err error) {
	model := tx.Model("t_inf_bill_aliyun_summary_split_share")
	model.Where("billing_cycle =", currentMonth)
	_, err = model.Delete()
	return err
}
func SaveSplitShare(data []*typeMNLZ.AliyunBillSummaryDay, tx gdb.DB) (int64, error) {
	_, err := tx.Model("t_inf_bill_aliyun_summary_split_share").Data(&data).Save()
	if err != nil {
		return 0, err
	}
	return int64(len(data)), err
}

func QueryAllShareData(currentMonth string, tx gdb.DB) ([]*typeMNLZ.BillOverviewShare, error) {
	results := make([]*typeMNLZ.BillOverviewShare, 0)
	model := tx.Model("t_inf_bill_overview_all_share").
		Where("billing_cycle=?", currentMonth).
		Where("(pay_as_you_go_cost>0 or all_share>0 or subscription_cost>0 )").Order("app_name")
	err := model.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func QueryAllShareByProduct(billingCycle string, tx gdb.DB) (allAllShareMoney map[string]*typeMNLZ.ShareByProduct, err error) {
	allAllShareMoney = make(map[string]*typeMNLZ.ShareByProduct)
	var temp []*typeMNLZ.ShareByProduct
	model := tx.Model("t_inf_bill_aliyun_summary_split_positive_day").
		Fields("sum(pretax_amount) as share_money", "product_name", "subscription_type").
		Where("pretax_amount>0").
		Where("application='ALL_SHARE'").
		Where("billing_cycle=?", billingCycle).Group("subscription_type").Group("product_name")
	err = model.Scan(&temp)
	if err != nil {
		panic(err)
	}
	for i := range temp {
		allAllShareMoney[temp[i].ProductName] = temp[i]
	}
	return
}

// FindDateDataByPositive
// @Description: 实时更新当前日期账单数据
// @Author: Evan Lu
// @Date: 2022-12-01 11:35:00
// @Param tableName string
// @Param today time.Time
// @Return: int64
// @Return: error
func FindDateDataByPositive(tableName string, today time.Time, clickhouse gdb.DB) ([]*typeMNLZ.AliyunBillSummaryDay, error) {
	result := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := clickhouse.Model(tableName)
	model.Where("(service_period_unit='月' and (TIMESTAMPDIFF(day,date_add(payment_time,interval service_period month),'" + today.Format("2006-01-02") + "')<0) AND (`billing_cycle`='" + today.Format("2006-01") + "')) or (service_period_unit='年' and (TIMESTAMPDIFF(day,date_add(payment_time,interval service_period*12 month),'" + today.Format("2006-01-02") + "')<0) AND (`billing_cycle`='" + today.Format("2006-01") + "')) or (service_period_unit='天' and (TIMESTAMPDIFF(day,date_add(payment_time,interval service_period day),'" + today.Format("2006-01-02") + "')<0) AND (`billing_cycle`='" + today.Format("2006-01") + "'))")
	//model.Where("(service_period_unit='月' and (DATE_DIFF(day,DATE_ADD(month,toInt32(service_period),toDate(payment_time)),toDate('" + today.Format("2006-01-02") + "'))<0) AND (`billing_cycle`='" + today.Format("2006-01") + "')) or (service_period_unit='年' and (DATE_DIFF(day,DATE_ADD(month,toInt32(service_period)*12,toDate(payment_time)),toDate('" + today.Format("2006-01-02") + "'))<0) AND (`billing_cycle`='" + today.Format("2006-01") + "')) or (service_period_unit='天' and (DATE_DIFF(day,DATE_ADD(day,toInt32(service_period),toDate(payment_time)),toDate('" + today.Format("2006-01-02") + "'))<0) AND (`billing_cycle`='" + today.Format("2006-01") + "'))")
	err := model.Scan(&result)
	return result, err
}

func GetShareCostList(billingCycle string, tenantID int64, tx *gdb.TX) ([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, error) {
	result := make([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, 0)
	model := tx.Model("t_inf_bill_overview_all_share").Fields("SUM(total_cost) as total_cost, app_name, billing_cycle").
		Where("billing_cycle = ?", billingCycle).
		Where("tenant_id", tenantID).
		Group("app_id,billing_cycle")
	err := model.Scan(&result)
	return result, err
}

func GetTotalShareCostList(billingCycle string, tenantID int64, tx *gdb.TX) ([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, error) {
	result := make([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, 0)
	model := tx.Model("t_inf_bill_overview_all_share").Fields("SUM(total_cost) as total_cost, billing_cycle,subscription_type")
	model.Where("billing_cycle = ?", billingCycle).Where("app_name =", "ALL_SHARE")
	model.Where("tenant_id", tenantID)
	model.Group("billing_cycle, subscription_type")
	err := model.Scan(&result)
	return result, err
}

func GetTotalAppCostList(billingCycle string, tenantID int64, tx *gdb.TX) ([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, error) {
	result := make([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, 0)
	model := tx.Model("t_inf_bill_overview_all_share").Fields("SUM(total_cost) as total_cost,account_id,account_name,app_name,cloud_room_id,cloud_room_name,billing_cycle,subscription_type")
	model.Where("billing_cycle = ?", billingCycle)
	model.Where("tenant_id", tenantID)
	model.Group("app_id,cloud_room_id,billing_cycle,subscription_type")
	err := model.Scan(&result)
	return result, err
}

func GetOverviewListInfo(billingCycle string, tenantID int64, tx *gdb.TX) ([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, error) {
	result := make([]*typeMNLZ.TInfBillOverviewSummaryMonthItem, 0)
	model := tx.Model("t_inf_bill_overview_all_share")
	model.Where("billing_cycle =", billingCycle)
	model.Where("tenant_id", tenantID)
	err := model.Scan(&result)
	return result, err
}

func QueryTInfBillAliyunSummaryDay(date string) ([]*typeMNLZ.AliyunBillSummaryDay, error) {
	mysql, err := gdb.NewByGroup("mysql")
	if err != nil {
		log.Error("", "%v", err)
		return nil, err
	}
	results := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := mysql.Model("t_inf_bill_aliyun_summary_day").
		Where("billing_date = ?", date)
	err = model.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func UpdateTaskStatus(taskStatus int) {
	mysql, err := gdb.NewByGroup("mysql")
	if err != nil {
		log.Error("", "%v", err)
	}

	result, _ := mysql.Schema("cloudproject").Model("s_tasks").
		Where("task_id", 101).Data(
		g.Map{
			"task_state":          taskStatus,
			"task_exec_last_Date": time.Now().Format("2006-01-02 15:04:05"),
			"updated_time":        time.Now().Format("2006-01-02 15:04:05"),
		}).Update()
	log.Info("result", fmt.Sprintf("%v", result))
}

func SavePositiveDay(date string, mysql, clickhouse gdb.DB) error {
	log.Info("", "UpdatePositiveDay method,currentMonth : %s ", date)
	results := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := mysql.Model("t_inf_bill_aliyun_summary_split_positive_day").
		Where("billing_cycle = ?", date)
	err := model.Scan(&results)
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		_, err = clickhouse.Save(clickhouse.GetCtx(), "t_inf_bill_aliyun_summary_split_positive_day", &results)
		if err != nil {
			panic(err)
		}
		return nil
	}
	return errors.New("data is nil")
}

func TruncatePositiveDay(clickhouse gdb.DB) error {
	result := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := clickhouse.Model("t_inf_bill_aliyun_summary_split_positive_day")
	err := model.Scan(&result)
	if err != nil {
		return err
	}
	if len(result) > 0 {
		_, err := clickhouse.Exec(clickhouse.GetCtx(), "TRUNCATE table t_inf_bill_aliyun_summary_split_positive_day")
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func UpdatePositiveDay(date string, mysql, clickhouse gdb.DB) error {
	log.Info("", "SavePositiveDay method,currentMonth : %s ", date)
	results := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := mysql.Model("t_inf_bill_aliyun_summary_split_positive_day").
		Where("billing_cycle = ?", date)
	err := model.Scan(&results)
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		_, err = clickhouse.Save(clickhouse.GetCtx(), "t_inf_bill_aliyun_summary_split_positive_day", &results)
		if err != nil {
			panic(err)
		}
		return nil
	}
	return errors.New("data is nil")
}

func UpdateNegativeDay(mysql, clickhouse gdb.DB) error {
	log.Info("", "UpdateNegativeDay method start")
	results := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	result := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := clickhouse.Model("t_inf_bill_aliyun_summary_split_negative_day")
	err := model.Scan(&result)
	if err != nil {
		panic(err)
	}
	if len(result) > 0 {
		_, err := clickhouse.Exec(clickhouse.GetCtx(), "TRUNCATE table t_inf_bill_aliyun_summary_split_negative_day")
		if err != nil {
			panic(err)
		}
	}
	model = mysql.Model("t_inf_bill_aliyun_summary_split_negative_day").Order("billing_cycle")
	err = model.Scan(&results)
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		_, err = clickhouse.Save(clickhouse.GetCtx(), "t_inf_bill_aliyun_summary_split_negative_day", &results)
		if err != nil {
			panic(err)
		}
		return nil
	}
	return errors.New("data is nil")
}

func SaveOverviewShare(data []*typeMNLZ.BillOverviewShare, tableName string, tx gdb.DB) (int64, error) {
	_, err := tx.Model(tableName).Data(&data).Save()
	if err != nil {
		return 0, err
	}
	return int64(len(data)), err
}

func DeleteOverviewAllShareByAppName(tenantID int64, appName string, tx *gdb.TX) (err error) {
	model := tx.Model("t_inf_bill_overview_all_share")
	model.Where("tenant_id=?", tenantID)
	model.Where("app_name=?", appName)
	_, err = model.Delete()
	return err
}

func DeleteOverviewAllShare(currentMonth string, tx gdb.DB) (err error) {
	model := tx.Model("t_inf_bill_overview_all_share")
	model.Where("billing_cycle =", currentMonth)
	_, err = model.Delete()
	return err
}

func QueryAllAppPositive(currentMonth string, tx gdb.DB) ([]*typeMNLZ.AliyunBillSummaryDay, error) {
	results := make([]*typeMNLZ.AliyunBillSummaryDay, 0)
	model := tx.Model("t_inf_bill_aliyun_summary_split_positive_day").
		Where("billing_cycle=?", currentMonth).
		Where("application!='ALL_SHARE'")
	err := model.Scan(&results)
	if err != nil {
		return nil, err
	}
	return results, nil
}

func QueryAllShareMoney(billingCycle string, tx gdb.DB) (allAllShareMoney []float64, err error) {
	data, err := tx.Model("t_inf_bill_aliyun_summary_split_positive_day").
		Fields("sum(pretax_amount) as share_money").
		Where("pretax_amount>0").
		Where("application='ALL_SHARE'").
		Where("billing_cycle=?", billingCycle).Group("subscription_type").Order("subscription_type").Limit(2).All()
	for i := range data.Array() {
		allAllShareMoney = append(allAllShareMoney, data.Array()[i].Float64())
	}
	return
}
