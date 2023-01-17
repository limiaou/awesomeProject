package split

import (
	"awesomeProject/log"
	"awesomeProject/snow_flake"
	"awesomeProject/split/rdb"
	typeMNLZ "awesomeProject/split/type"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"math"
	"strconv"
	"time"
)

func SubMonthlyAndYearlyAmortization(currentMonth, date string) error {
	log.Info("", "SubMonthlyAndYearlyAmortization method,currentMonth : %s , date : %s", currentMonth, date)
	clickhouse, err := gdb.NewByGroup("mysql")
	if err != nil {
		log.Error("", "%v", err)
		return err
	}
	defer clickhouse.Close(clickhouse.GetCtx())
	clickhouse.SetDebug(false)
	// 清洗包年包月分摊
	aliyunSummaryDays, err := rdb.QueryTInfBillAliyunSummaryDay(date)
	if err != nil {
		log.Error("", "SubMonthlyAndYearlyAmortization method", "QueryTInfBillAliyunSummaryDay err : %s", err.Error())
		return err
	}
	var result, results, tempInserts, tempInserts2 []*typeMNLZ.AliyunBillSummaryDay
	for i := range aliyunSummaryDays {
		tempInserts, tempInserts2, err = execWashingData(aliyunSummaryDays[i], 32, currentMonth)
		if err != nil {
			log.Error("", "execWashingData method", "execWashingData err : %s", err.Error())
			return err
		}
		result = append(result, tempInserts...)
		results = append(results, tempInserts2...)
	}
	err = rdb.DeleteSummaryDay(date, `t_inf_bill_aliyun_summary_split_positive_day`, clickhouse)
	if err != nil {
		log.Error("", "aliyun-DeletePositiveDay失败", err)
		return err
	}
	if len(result) > 0 {
		err = preInsert(result, clickhouse)
		//_, err = rdb.SaveSummaryDay(result, `t_inf_bill_aliyun_summary_split_positive_day`, clickhouse)
		if err != nil {
			log.Error("", "aliyun-SavePositiveDay失败", err)
			return err
		}
	}
	err = rdb.DeleteSummaryDay(date, `t_inf_bill_aliyun_summary_split_negative_day`, clickhouse)
	if err != nil {
		log.Error("", "aliyun-DeleteNegativeDay失败", err)
		return err
	}
	if len(results) > 0 {
		_, err = rdb.SaveSummaryDay(results, `t_inf_bill_aliyun_summary_split_negative_day`, clickhouse)
		if err != nil {
			log.Error("", "aliyun-SaveNegativeDay失败", err)
			return err
		}
	}
	parse, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Error("", "aliyun-日期格式转换失败", err)
		return err
	}

	// START 动态更新t_inf_bill_aliyun_summary_split_positive_day表中的值
	positiveData, err := rdb.FindDateDataByPositive("t_inf_bill_aliyun_summary_split_positive_day", parse, clickhouse)
	if err != nil {
		log.Error("", "aliyun-日期格式转换失败", err)
		return err
	}
	for i := range positiveData {

		if positiveData[i].Deleted == 2 {
			positiveData[i].Deleted = 0
		}
		currentTime, err := time.Parse("2006-01-02", positiveData[i].BillingDate)
		if err != nil {
			log.Error("", "aliyun-日期格式转换失败", err)
		}
		if currentTime.Format("2006-01") == currentMonth {
			positiveData[i].PaymentAmount = positiveData[i].PretaxGrossAmount * float64(parse.Day()-currentTime.Day()+1)
			positiveData[i].PretaxAmount = positiveData[i].PretaxGrossAmount * float64(parse.Day()-currentTime.Day()+1)
		} else {
			positiveData[i].PaymentAmount = positiveData[i].PretaxGrossAmount * float64(parse.Day())
			positiveData[i].PretaxAmount = positiveData[i].PretaxGrossAmount * float64(parse.Day())
		}
	}
	err = preInsert(positiveData, clickhouse)
	//_, err = rdb.SaveSummaryDay(positiveData, `t_inf_bill_aliyun_summary_split_positive_day`, clickhouse)
	if err != nil {
		log.Error("", "aliyun-实时更新金额失败", err)
		return err
	}
	// ALL_SHARE table
	// 共享资源费用分摊
	//err = SplitAllShareCosts(currentMonth, clickhouse)
	//if err != nil {
	//	return err
	//}
	//saveOrUpdateData(currentMonth, clickhouse)
	return nil
}

func execWashingData(temp *typeMNLZ.AliyunBillSummaryDay, seed int, currentMonth string) ([]*typeMNLZ.AliyunBillSummaryDay, []*typeMNLZ.AliyunBillSummaryDay, error) {
	var tempInserts, tempInserts2 []*typeMNLZ.AliyunBillSummaryDay
	if temp.PretaxAmount >= 0 {
		var deleted int8
		switch temp.ServicePeriodUnit {
		case "年":
			month, _ := strconv.Atoi(temp.ServicePeriod)
			moneys, months, tempMoney := subByYear(temp.PaymentTime, month, temp.PretaxAmount, currentMonth)
			tempID, err := snow_flake.GenerateServiceIdentities(seed, len(moneys))
			if err != nil {
				return nil, nil, err
			}
			for j := range months {
				if moneys[j] == 0 {
					deleted = 2
				} else {
					deleted = 0
				}
				tempInserts = append(tempInserts, &typeMNLZ.AliyunBillSummaryDay{
					SummaryId:             tempID[j],
					ProductId:             temp.ProductId,
					ResourceId:            temp.ResourceId,
					AccountId:             temp.AccountId,
					CloudAccountId:        temp.CloudAccountId,
					Application:           temp.Application,
					ProductSource:         temp.ProductSource,
					CloudRoomId:           temp.CloudRoomId,
					CloudRoomName:         temp.CloudRoomName,
					ResourceGroupId:       temp.ResourceGroupId,
					BillingCycle:          months[j],
					BillingDate:           temp.BillingDate,
					BillingWeek:           temp.BillingWeek,
					ProductCode:           temp.ProductCode,
					ProductName:           temp.ProductName,
					ProductType:           temp.ProductType,
					ProductDetail:         temp.ProductDetail,
					SubscriptionType:      temp.SubscriptionType,
					PaymentTime:           temp.PaymentTime,
					PretaxGrossAmount:     tempMoney,
					PretaxAmount:          moneys[j],
					PaymentAmount:         moneys[j],
					OutstandingAmount:     temp.OutstandingAmount,
					DeductedByPrepaidCard: temp.DeductedByPrepaidCard,
					DeductedByCashCoupons: temp.DeductedByCashCoupons,
					DeductedByCredit:      temp.DeductedByCredit,
					Currency:              temp.Currency,
					PipCode:               temp.PipCode,
					InstanceId:            temp.InstanceId,
					NickName:              temp.NickName,
					Region:                temp.Region,
					ServicePeriod:         temp.ServicePeriod,
					ServicePeriodUnit:     temp.ServicePeriodUnit,
					UsageUnit:             temp.UsageUnit,
					CreatedTime:           temp.CreatedTime,
					UpdatedTime:           temp.UpdatedTime,
					TenantId:              temp.TenantId,
					RecordId:              temp.RecordId,
					CloudAccountName:      temp.CloudAccountName,
					ResourceGroupName:     temp.ResourceGroupName,
					ProductTypeCnName:     temp.ProductTypeCnName,
					GroupCnName:           temp.GroupCnName,
					ProductTypeEnName:     temp.ProductTypeEnName,
					GroupEnName:           temp.GroupEnName,
					AppId:                 temp.AppId,
					Deleted:               deleted,
				})
			}
		case "月":
			month, _ := strconv.Atoi(temp.ServicePeriod)
			moneys, months, tempMoney := subByMonth(temp.PaymentTime, month, temp.PretaxAmount, currentMonth)
			tempID, err := snow_flake.GenerateServiceIdentities(seed, len(moneys))
			if err != nil {
				panic(err)
			}
			for j := range months {
				if moneys[j] == 0 {
					deleted = 2
				} else {
					deleted = 0
				}
				tempInserts = append(tempInserts, &typeMNLZ.AliyunBillSummaryDay{
					SummaryId:             tempID[j],
					ProductId:             temp.ProductId,
					ResourceId:            temp.ResourceId,
					AccountId:             temp.AccountId,
					CloudAccountId:        temp.CloudAccountId,
					Application:           temp.Application,
					ProductSource:         temp.ProductSource,
					CloudRoomId:           temp.CloudRoomId,
					CloudRoomName:         temp.CloudRoomName,
					ResourceGroupId:       temp.ResourceGroupId,
					BillingCycle:          months[j],
					BillingDate:           temp.BillingDate,
					BillingWeek:           temp.BillingWeek,
					ProductCode:           temp.ProductCode,
					ProductName:           temp.ProductName,
					ProductType:           temp.ProductType,
					ProductDetail:         temp.ProductDetail,
					SubscriptionType:      temp.SubscriptionType,
					PaymentTime:           temp.PaymentTime,
					PretaxGrossAmount:     tempMoney,
					PretaxAmount:          moneys[j],
					PaymentAmount:         moneys[j],
					OutstandingAmount:     temp.OutstandingAmount,
					DeductedByPrepaidCard: temp.DeductedByPrepaidCard,
					DeductedByCashCoupons: temp.DeductedByCashCoupons,
					DeductedByCredit:      temp.DeductedByCredit,
					Currency:              temp.Currency,
					PipCode:               temp.PipCode,
					InstanceId:            temp.InstanceId,
					NickName:              temp.NickName,
					Region:                temp.Region,
					ServicePeriod:         temp.ServicePeriod,
					ServicePeriodUnit:     temp.ServicePeriodUnit,
					UsageUnit:             temp.UsageUnit,
					CreatedTime:           temp.CreatedTime,
					UpdatedTime:           temp.UpdatedTime,
					TenantId:              temp.TenantId,
					RecordId:              temp.RecordId,
					CloudAccountName:      temp.CloudAccountName,
					ResourceGroupName:     temp.ResourceGroupName,
					ProductTypeCnName:     temp.ProductTypeCnName,
					GroupCnName:           temp.GroupCnName,
					ProductTypeEnName:     temp.ProductTypeEnName,
					GroupEnName:           temp.GroupEnName,
					AppId:                 temp.AppId,
					Deleted:               deleted,
				})
			}
		case "天":
			days, _ := strconv.Atoi(temp.ServicePeriod)
			moneys, months, tempMoney := subByDay(temp.PaymentTime, days, temp.PretaxAmount, currentMonth)
			tempID, err := snow_flake.GenerateServiceIdentities(seed, len(moneys))
			if err != nil {
				panic(err)
			}
			for j := range months {
				if moneys[j] == 0 {
					deleted = 2
				} else {
					deleted = 0
				}
				tempInserts = append(tempInserts, &typeMNLZ.AliyunBillSummaryDay{
					SummaryId:             tempID[j],
					ProductId:             temp.ProductId,
					ResourceId:            temp.ResourceId,
					AccountId:             temp.AccountId,
					CloudAccountId:        temp.CloudAccountId,
					Application:           temp.Application,
					ProductSource:         temp.ProductSource,
					CloudRoomId:           temp.CloudRoomId,
					CloudRoomName:         temp.CloudRoomName,
					ResourceGroupId:       temp.ResourceGroupId,
					BillingCycle:          months[j],
					BillingDate:           temp.BillingDate,
					BillingWeek:           temp.BillingWeek,
					ProductCode:           temp.ProductCode,
					ProductName:           temp.ProductName,
					ProductType:           temp.ProductType,
					ProductDetail:         temp.ProductDetail,
					SubscriptionType:      temp.SubscriptionType,
					PaymentTime:           temp.PaymentTime,
					PretaxGrossAmount:     tempMoney,
					PretaxAmount:          moneys[j],
					PaymentAmount:         moneys[j],
					OutstandingAmount:     temp.OutstandingAmount,
					DeductedByPrepaidCard: temp.DeductedByPrepaidCard,
					DeductedByCashCoupons: temp.DeductedByCashCoupons,
					DeductedByCredit:      temp.DeductedByCredit,
					Currency:              temp.Currency,
					PipCode:               temp.PipCode,
					InstanceId:            temp.InstanceId,
					NickName:              temp.NickName,
					Region:                temp.Region,
					ServicePeriod:         temp.ServicePeriod,
					ServicePeriodUnit:     temp.ServicePeriodUnit,
					UsageUnit:             temp.UsageUnit,
					CreatedTime:           temp.CreatedTime,
					UpdatedTime:           temp.UpdatedTime,
					TenantId:              temp.TenantId,
					RecordId:              temp.RecordId,
					CloudAccountName:      temp.CloudAccountName,
					ResourceGroupName:     temp.ResourceGroupName,
					ProductTypeCnName:     temp.ProductTypeCnName,
					GroupCnName:           temp.GroupCnName,
					ProductTypeEnName:     temp.ProductTypeEnName,
					GroupEnName:           temp.GroupEnName,
					AppId:                 temp.AppId,
					Deleted:               deleted,
				})
			}
		default:
			tempInserts = append(tempInserts, temp)
		}
	} else {
		tempInserts2 = append(tempInserts2, temp)
	}
	return tempInserts, tempInserts2, nil
}

func subByYear(startTime string, year int, money float64, current string) (moneys []float64, months []string, oneDay float64) {
	allMonth := year * 12
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	currentTime, err := time.Parse("2006-01", current)
	if err != nil {
		panic(err)
	}
	currentYear := currentTime.Year()
	currentMonth := currentTime.Month()
	monthNow := parse.Month()
	yearNow := parse.Year()
	day := parse.Day()
	allDay := (parse.AddDate(year, 0, 0).Sub(parse).Hours()) / 24
	oneDay = money / allDay
	for i := 0; i < allMonth; i++ {
		if yearNow <= currentYear && monthNow <= currentMonth {
			moneys = append(moneys, getMouthDay(parse)*oneDay)
			months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
		} else {
			moneys = append(moneys, 0)
			months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
		}
		monthNow++
		if monthNow > 12 {
			yearNow++
			monthNow %= 12
		}
		parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
		parse = parse.AddDate(0, 1, 0)
	}
	if day != 1 {
		if yearNow <= currentYear && monthNow <= currentMonth {
			moneys = append(moneys, float64(day-1)*oneDay)
			months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
		} else {
			moneys = append(moneys, 0)
			months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
		}
	}
	return
}
func subByMonth(startTime string, month int, money float64, current string) (moneys []float64, months []string, oneDay float64) {
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	currentTime, err := time.Parse("2006-01", current)
	if err != nil {
		panic(err)
	}
	currentYear := currentTime.Year()
	currentMonth := currentTime.Month()
	monthNow := parse.Month()
	year := parse.Year()
	day := parse.Day()
	allDay := (parse.AddDate(0, month, 0).Sub(parse).Hours()) / 24
	oneDay = money / allDay
	for i := 0; i < month; i++ {
		if year <= currentYear && monthNow <= currentMonth {
			moneys = append(moneys, getMouthDay(parse)*oneDay)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		} else {
			moneys = append(moneys, 0)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		}
		monthNow++
		if monthNow > 12 {
			year++
			monthNow %= 12
		}
		parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
		parse = parse.AddDate(0, 1, 0)
	}
	if day != 1 {
		if year <= currentYear && monthNow <= currentMonth {
			moneys = append(moneys, float64(day-1)*oneDay)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		} else {
			moneys = append(moneys, 0)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		}
	}
	return
}

func subByDay(startTime string, days int, money float64, current string) (moneys []float64, months []string, oneDay float64) {
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	currentTime, err := time.Parse("2006-01", current)
	if err != nil {
		panic(err)
	}
	currentYear := currentTime.Year()
	currentMonth := currentTime.Month()
	monthNow := parse.Month()
	year := parse.Year()
	day := parse.Day()
	oneDay = money / float64(days)
	for days > 30 {
		mouthDay := getMouthDay(parse)
		if year <= currentYear && monthNow <= currentMonth {
			moneys = append(moneys, mouthDay*oneDay)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		} else {
			moneys = append(moneys, 0)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		}
		monthNow++
		if monthNow > 12 {
			year++
			monthNow %= 12
		}
		days -= int(mouthDay)
		parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
		parse = parse.AddDate(0, 1, 0)
	}
	if day != 1 {
		if year <= currentYear && monthNow <= currentMonth {
			moneys = append(moneys, float64(days)*oneDay)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		} else {
			moneys = append(moneys, 0)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
		}
	}
	return
}

func getMouthDay(t time.Time) (day float64) {
	parse, _ := time.Parse("2006-01-02", t.Format("2006-01")+"-01")
	return math.Ceil(parse.AddDate(0, 1, -1).Sub(t).Hours()/24) + 1
}

func returnStrMonth(i int) string {
	if i < 10 {
		return "0" + strconv.Itoa(i)
	} else {
		return strconv.Itoa(i)
	}
}

func SplitAllShareCosts(currentMonth string, clickhouse gdb.DB) error {
	log.Info("", "SplitAllShareCosts method,currentMonth : %s ", currentMonth)
	err := rdb.DeleteOverviewAllShare(currentMonth, clickhouse)
	if err != nil {
		panic(err)
	}
	after := make(map[string][]*typeMNLZ.BillOverviewShare)
	positive, err := rdb.QueryAllAppPositive(currentMonth, clickhouse)
	if err != nil {
		panic(err)
	}
	allShareMoney, err := rdb.QueryAllShareMoney(currentMonth, clickhouse)
	if err != nil {
		panic(err)
	}
	var subSum, paySum float64
	for i := range positive {
		if _, boo := after[positive[i].Application]; !boo {
			tempID, err := snow_flake.GenerateServiceIdentities(33, 5)
			if err != nil {
				panic(err)
			}
			after[positive[i].Application] = returnSlice(positive[i], tempID)
		}
		if positive[i].SubscriptionType == "Subscription" {
			subSum += positive[i].PretaxAmount
			after[positive[i].Application][0].SubscriptionCost += positive[i].PretaxAmount
		} else {
			paySum += positive[i].PretaxAmount
			after[positive[i].Application][1].PayAsYouGoCost += positive[i].PretaxAmount
		}
	}
	var saveDate []*typeMNLZ.BillOverviewShare
	for _, shares := range after {
		ratioSub, ratioPay := shares[0].SubscriptionCost/subSum, shares[1].PayAsYouGoCost/paySum
		shareMoneySub, shareMoneyPay := ratioSub*allShareMoney[0], ratioPay*allShareMoney[1]
		shares[2].AllShare = shareMoneyPay + shareMoneySub
		if shares[0].SubscriptionCost > 0 {
			shares[3].AllSubscription = shares[0].SubscriptionCost + shareMoneySub
		}
		if shares[1].PayAsYouGoCost > 0 {
			shares[4].AllPayAsYouGo = shares[1].PayAsYouGoCost + shareMoneyPay
		}
		shares[0].TotalCost = shares[0].SubscriptionCost
		shares[1].TotalCost = shares[1].PayAsYouGoCost
		shares[2].TotalCost = shares[2].AllShare
		saveDate = append(saveDate, shares...)
	}
	_, err = rdb.SaveOverviewShare(saveDate, "t_inf_bill_overview_all_share", clickhouse)
	if err != nil {
		log.Error("SaveOverviewSummaryMonth Error ", "发生错误:", err)
		return err
	}
	return err
}

func returnSlice(positive *typeMNLZ.AliyunBillSummaryDay, tempID []int64) []*typeMNLZ.BillOverviewShare {
	return []*typeMNLZ.BillOverviewShare{
		{
			Uuid:             tempID[0],
			AccountId:        positive.AccountId,
			BillingCycle:     positive.BillingCycle,
			AppId:            positive.AppId,
			AppName:          positive.Application,
			CloudRoomName:    positive.CloudRoomName,
			TotalCost:        0,
			SubscriptionType: 2,
			CloudType:        0,
			BillingType:      0,
		},
		{
			Uuid:             tempID[1],
			AccountId:        positive.AccountId,
			BillingCycle:     positive.BillingCycle,
			AppId:            positive.AppId,
			AppName:          positive.Application,
			CloudRoomName:    positive.CloudRoomName,
			TotalCost:        0,
			SubscriptionType: 2,
			CloudType:        0,
			BillingType:      0,
		},
		{
			Uuid:             tempID[2],
			AccountId:        positive.AccountId,
			BillingCycle:     positive.BillingCycle,
			AppId:            positive.AppId,
			AppName:          positive.Application,
			CloudRoomName:    positive.CloudRoomName,
			TotalCost:        0,
			SubscriptionType: 2,
			CloudType:        0,
			BillingType:      0,
		},
		{
			Uuid:             tempID[3],
			AccountId:        positive.AccountId,
			BillingCycle:     positive.BillingCycle,
			AppId:            positive.AppId,
			AppName:          positive.Application,
			CloudRoomName:    positive.CloudRoomName,
			TotalCost:        0,
			SubscriptionType: 2,
			CloudType:        0,
			BillingType:      0,
		},
		{
			Uuid:             tempID[4],
			AccountId:        positive.AccountId,
			BillingCycle:     positive.BillingCycle,
			AppId:            positive.AppId,
			AppName:          positive.Application,
			CloudRoomName:    positive.CloudRoomName,
			TotalCost:        0,
			SubscriptionType: 2,
			CloudType:        0,
			BillingType:      0,
		},
	}
}

func saveOrUpdateData(curMonth string, tx gdb.DB) {
	log.Info("", "SaveOrUpdateData method,currentMonth : %s ", curMonth)
	// 查询positive金额
	positive, err := rdb.QueryUpdatePositive(curMonth, tx)
	if err != nil {

		panic("select positive data err ")
	}
	// 删除当前月份历史数据
	err = rdb.DeleteSpiltShare(curMonth, tx)
	if err != nil {

		panic("select positive data err ")
	}
	// 生成五个维度的共享资源 计算比例
	proportion := calculationProportion(curMonth, tx)

	allShare, err := rdb.QueryAllShareData(curMonth, tx)
	if err != nil {
		panic("select positive data err ")
	}
	shareMoney := make(map[string]*typeMNLZ.TempShare)
	for i := range allShare {
		if v, ok := shareMoney[allShare[i].AppName]; ok {
			v.SaveMoney(allShare[i])
		} else {
			shareMoney[allShare[i].AppName] = (&typeMNLZ.TempShare{}).SaveMoney(allShare[i])
		}
	}
	after := make(map[typeMNLZ.ByProduct]*typeMNLZ.AliyunBillSummaryDay)
	for i := range positive {
		if v, ok := shareMoney[positive[i].Application]; ok {
			if _, boo := after[typeMNLZ.ByProduct{CloudRoomName: positive[i].Application, ProductName: "网络"}]; !boo {
				// 根据比例计算分摊金额
				addMap(after, positive[i], v.AllShare, proportion)
			}
			if v1, boo := after[typeMNLZ.ByProduct{ProductName: positive[i].ProductName, CloudRoomName: positive[i].CloudRoomName}]; boo {
				if positive[i].SubscriptionType == "Subscription" {
					v1.DeductedByPrepaidCard += positive[i].PretaxAmount
				} else {
					v1.DeductedByCashCoupons += positive[i].PretaxAmount
				}
				v1.PretaxAmount += positive[i].PretaxAmount
			} else {
				trs := &typeMNLZ.AliyunBillSummaryDay{
					ProductId:        positive[i].ProductId,
					BillingCycle:     positive[i].BillingCycle,
					BillingDate:      positive[i].BillingDate,
					ProductName:      positive[i].ProductName,
					ProductType:      positive[i].ProductType,
					ProductDetail:    positive[i].ProductDetail,
					SubscriptionType: positive[i].SubscriptionType,
					PaymentTime:      positive[i].PaymentTime,
					CloudRoomName:    positive[i].CloudRoomName,
					PretaxAmount:     positive[i].PretaxAmount,
					AppId:            positive[i].AppId,
					Application:      positive[i].Application,
				}
				if positive[i].SubscriptionType == "Subscription" {
					trs.DeductedByPrepaidCard = positive[i].PretaxAmount
				} else {
					trs.DeductedByCashCoupons = positive[i].PretaxAmount
				}
				after[typeMNLZ.ByProduct{ProductName: positive[i].ProductName, CloudRoomName: positive[i].CloudRoomName}] = trs
			}
		}
	}
	tempID, err := snow_flake.GenerateServiceIdentities(33, len(after))
	if err != nil {

		panic(err)
	}
	i := 0
	for _, value := range after {
		value.SummaryId = tempID[i]
		i++
	}
	var answerSlice []*typeMNLZ.AliyunBillSummaryDay
	for _, moneys := range after {
		answerSlice = append(answerSlice, moneys)
	}
	_, err = rdb.SaveSplitShare(answerSlice, tx)
	if err != nil {

		panic(err)
	}
}

func calculationProportion(curMonth string, tx gdb.DB) (proportion map[string]float64) {
	proportion = map[string]float64{
		"IT支持系统": 0.0,
		"安全":     0.0,
		"存储":     0.0,
		"计算":     0.0,
		"网络":     0.0,
	}
	data, err := rdb.QueryAllShareByProduct(curMonth, tx)
	if err != nil {
		panic(err)
	}
	maxMoney := 0.0
	for k, v := range data {
		if _, ok := proportion[data[k].ProductName]; ok {
			proportion[data[k].ProductName] += v.ShareMoney
			maxMoney += v.ShareMoney
		} else {
			proportion["IT支持系统"] += v.ShareMoney
			maxMoney += v.ShareMoney
		}
	}
	for s, f := range proportion {
		proportion[s] = f / maxMoney
	}
	return
}

// addMap 生成五个维度的共享资源 计算比例
func addMap(m map[typeMNLZ.ByProduct]*typeMNLZ.AliyunBillSummaryDay, positive *typeMNLZ.AliyunBillSummaryDay, money float64, proportion map[string]float64) {
	val := &typeMNLZ.AliyunBillSummaryDay{
		ProductId:        positive.ProductId,
		BillingCycle:     positive.BillingCycle,
		BillingDate:      positive.BillingDate,
		ProductName:      "网络",
		ProductType:      positive.ProductType,
		ProductDetail:    "共享资源——网络",
		SubscriptionType: "Share",
		PaymentTime:      positive.PaymentTime,
		PretaxAmount:     proportion["网络"] * money,
		InvoiceDiscount:  proportion["网络"] * money,
		Application:      positive.Application,
		CloudRoomName:    positive.CloudRoomName,
		AppId:            positive.AppId,
	}
	m[typeMNLZ.ByProduct{CloudRoomName: positive.CloudRoomName, ProductName: "网络"}] = val
	list := []string{"IT支持系统", "安全", "存储", "计算"}
	for i := range list {
		temp := *val
		temp.ProductName, temp.ProductDetail, temp.PretaxAmount, temp.InvoiceDiscount = list[i], "共享资源——"+list[i], proportion[list[i]]*money, proportion[list[i]]*money
		m[typeMNLZ.ByProduct{CloudRoomName: positive.CloudRoomName, ProductName: list[i]}] = &temp
	}

}

func returnError(data ...interface{}) interface{} {
	if data[len(data)-1] != nil {
		panic(data[len(data)-1].(error))
	}
	return data[0]
}
func getMouthStartDay(t time.Time) (day float64) {
	parse, _ := time.Parse("2006-01-02", t.Format("2006-01")+"-01")
	return math.Ceil(t.Sub(parse).Hours() / 24)
}

//插入前判断是否大于900条,如果大于900就分批
func preInsert(acc []*typeMNLZ.AliyunBillSummaryDay, clickhouse gdb.DB) error {
	index := 0
	//如果不超过两千条即直接插入即可
	if len(acc) < 900 { //这里分为两千条插入一次
		//直接插入	//todo
		_, err := rdb.SaveSummaryDay(acc, `t_inf_bill_aliyun_summary_split_positive_day`, clickhouse)
		if err != nil {
			log.Error("", "SaveSummaryDay 900 Error")
		}
		return nil
	}
	//分批插入
	for index < len(acc) {
		if index+900 > len(acc) {
			fmt.Println("分批插入——插入区间", index, len(acc))
			//todo
			_, err := rdb.SaveSummaryDay(acc[index:], `t_inf_bill_aliyun_summary_split_positive_day`, clickhouse)
			if err != nil {
				log.Error("", "SaveSummaryDay 900 Error")
			}
			index = len(acc)
			break
		} else {
			fmt.Println("分批插入——插入区间", index, index+900)
			//todo
			_, err := rdb.SaveSummaryDay(acc[index:index+900], `t_inf_bill_aliyun_summary_split_positive_day`, clickhouse)
			if err != nil {
				log.Error("", "SaveSummaryDay %d Error", index+900)
			}
			index += 900
		}
	}
	return nil
}
