package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	//"log"
	"math"
	"strconv"
	"time"
)

//
//import (
//	"gorm.io/gorm"
//)
//
//type Product struct {
//	gorm.Model
//	Code  string
//	Price uint
//}
//
//func main() {
//	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
//	DSN := "root:root@tcp(127.0.0.1:3306)/giligili?charset=utf8mb4&parseTime=True&loc=Local"
//	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect database")
//	}
//	//
//	//// 迁移 schema
//	//db.AutoMigrate(&Product{})
//	//
//	//// Create
//	//db.Create(&Product{Code: "D42", Price: 100})
//	//
//	//// Read
//	//var product Product
//	//db.First(&product, 1)                 // 根据整型主键查找
//	//db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录
//	//
//	//// Update - 将 product 的 price 更新为 200
//	//db.Model(&product).Update("Price", 200)
//	//// Update - 更新多个字段
//	//db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
//	//db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
//	//
//	//// Delete - 删除 product
//	//db.Delete(&product, 1)
//}
type AliyunBillWashingQuery struct {
	ItemUuid              int64   `orm:"item_uuid" json:"item_uuid" desc:"详情账单自己生成uuid"`
	BillAccountId         int64   `orm:"bill_account_id" json:"bill_account_id" desc:"账号与账单关联uuid"`
	BillingCycle          string  `orm:"billing_cycle" json:"billing_cycle" desc:"账期，格式：YYYY-MM"`
	ItemType              string  `orm:"item_type" json:"item_type" desc:"账单类型：SubscriptionOrder／PayAsYouGoBill／Refund／Adjustment"`
	Application           string  `orm:"application" json:"application" desc:"应用（清洗)"`
	Environment           string  `orm:"environment" json:"environment" desc:"环境（清洗)"`
	Department            string  `orm:"department" json:"department" desc:"部门（清洗)"`
	Owner                 string  `orm:"owner" json:"owner" desc:"责任人"`
	Business              string  `orm:"business" json:"business" desc:"业务（清洗）"`
	CloudRoomId           string  `orm:"cloud_room_id" json:"cloud_room_id" desc:"云单元（清洗）"`
	ProductSource         string  `orm:"product_source" json:"product_source" desc:"产品（清洗）"`
	ProductCode           string  `orm:"product_code" json:"product_code" desc:"产品代码"`
	ProductName           string  `orm:"product_name" json:"product_name" desc:"产品名称"`
	ProductType           string  `orm:"product_type" json:"product_type" desc:"产品类型"`
	ProductDetail         string  `orm:"product_detail" json:"product_detail" desc:"产品明细"`
	SubscriptionType      string  `orm:"subscription_type" json:"subscription_type" desc:"订阅类型 Subscription／PayAsYouGo"`
	PaymentTime           string  `orm:"payment_time" json:"payment_time" desc:"订单支付时间"`
	UsageStartTime        string  `orm:"usage_start_time" json:"usage_start_time" desc:"账单开始时间"`
	UsageEndTime          string  `orm:"usage_end_time" json:"usage_end_time" desc:"账单结束时间"`
	RecordId              string  `orm:"record_id" json:"record_id" desc:"订单号/账单号"`
	PretaxGrossAmount     float64 `orm:"pretax_gross_amount" json:"pretax_gross_amount" desc:"原始金额,original cost"`
	DeductedByCoupons     float64 `orm:"deducted_by_coupons" json:"deducted_by_coupons" desc:"优惠券抵扣"`
	InvoiceDiscount       float64 `orm:"invoice_discount" json:"invoice_discount" desc:"优惠金额"`
	PretaxAmount          float64 `orm:"pretax_amount" json:"pretax_amount" desc:"应付金额"`
	PaymentAmount         float64 `orm:"payment_amount" json:"payment_amount" desc:"现金支付"`
	OutstandingAmount     float64 `orm:"outstanding_amount" json:"outstanding_amount" desc:"未结清金额"`
	Currency              string  `orm:"currency" json:"currency" desc:"币种 CNY／USD／JPY"`
	Status                string  `orm:"status" json:"status" desc:"支付状态:已支付（Cleared），未结清（Uncleared），未结算（Unsettled），免结算"`
	DeductedByPrepaidCard float64 `orm:"deducted_by_prepaid_card" json:"deducted_by_prepaid_card" desc:"代金券抵扣"`
	DeductedByCashCoupons float64 `orm:"deducted_by_cash_coupons" json:"deducted_by_cash_coupons" desc:"储值卡抵扣"`
	OwnerId               string  `orm:"owner_id" json:"owner_id" desc:"资源所有者的阿里云账号ID，通常为16位的数字。BID场景／分销场景OwnerID必须指定"`
	ResourceGroup         string  `orm:"resource_group" json:"resource_group" desc:"资源组"`
	AccountID             string  `orm:"accountID" json:"accountID" desc:"账号ID"`
	PipCode               string  `orm:"pip_code" json:"pip_code" desc:"产品Code"`
	InstanceConfig        string  `orm:"instance_config" json:"instance_config" desc:"实例详细配置"`
	InstanceSpec          string  `orm:"instance_spec" json:"instance_spec" desc:"实例规格"`
	InstanceId            string  `orm:"instance_id" json:"instance_id" desc:"实例ID"`
	InternetIp            string  `orm:"internet_ip" json:"internet_ip" desc:"公网IP"`
	IntranetIp            string  `orm:"intranet_ip" json:"intranet_ip" desc:"内网IP"`
	ListPrice             string  `orm:"list_price" json:"list_price" desc:"单价"`
	ListPriceUnit         string  `orm:"list_price_unit" json:"list_price_unit" desc:"单价单位"`
	NickName              string  `orm:"nick_name" json:"nick_name" desc:"实例昵称"`
	Region                string  `orm:"region" json:"region" desc:"地域"`
	ServicePeriod         string  `orm:"service_period" json:"service_period" desc:"服务时长"`
	ServicePeriodUnit     string  `orm:"service_period_unit" json:"service_period_unit" desc:"服务时长单位"`
	UsageNumber           string  `orm:"usage_number" json:"usage_number" desc:"用量"`
	UsageUnit             string  `orm:"usage_unit" json:"usage_unit" desc:"用量单位, 仅当isBillingItem为true时有效。"`
	Zone                  string  `orm:"zone" json:"zone" desc:"可用区"`
	Tags                  string  `orm:"tags" json:"tags" desc:"标签"`
	Deleted               int8    `orm:"deleted" json:"deleted" desc:"删除标志"`
	CreatedTime           string  `orm:"created_time" json:"created_time" desc:"创建时间"`
	UpdatedTime           string  `orm:"updated_time" json:"updated_time" desc:"更新时间"`
	TenantId              int64   `orm:"tenant_id" json:"tenant_id" desc:"租户ID"`
	CostUnit              string  `orm:"cost_unit" json:"cost_unit" desc:""`
	AccountName           string  `orm:"account_name" json:"account_name" desc:""`
	OwnerAccountName      string  `orm:"owner_account_name" json:"owner_account_name" desc:""`
	BillingDate           string  `orm:"billing_date" json:"billing_date" desc:""`
	RelatedRecordId       string  `orm:"related_record_id" json:"related_record_id" desc:""`
	BillingType           string  `orm:"billing_type" json:"billing_type" desc:""`
	BizType               string  `orm:"biz_type" json:"biz_type" desc:""`
	PaymentNo             string  `orm:"payment_no" json:"payment_no" desc:""`
	DeductedByCredit      float64 `orm:"deducted_by_credit" json:"deducted_by_credit" desc:""`
	PayType               string  `orm:"pay_type" json:"pay_type" desc:""`
	CloudAccountId        string  `orm:"cloud_account_id" json:"cloud_account_id" desc:""`
	CalculateTime         string  `orm:"calculate_time" json:"calculate_time" desc:""`
	UnionRGID             string  `orm:"union_rgid"  json:"union_rgid" desc:"rgunion主键"`
}

func main1() {
	saveWashingGroup := GetData()
	log.Println(len(saveWashingGroup))
	var tempInserts []*AliyunBillWashingQuery
	appendLens := 0

	for i := 0; i < len(saveWashingGroup); {
		temp := *saveWashingGroup[i]
		switch temp.ServicePeriodUnit {
		case "年":
			month, _ := strconv.Atoi(temp.ServicePeriod)
			moneys, months := subByYear(temp.PaymentTime, month, temp.PaymentAmount)
			appendLens += len(moneys)
			for j := range months {
				tempInserts = append(tempInserts, &AliyunBillWashingQuery{
					ItemUuid:              temp.ItemUuid,
					BillAccountId:         temp.BillAccountId,
					BillingCycle:          months[j],
					ProductCode:           temp.ProductCode,
					ProductName:           temp.ProductName,
					ProductType:           temp.ProductType,
					ProductDetail:         temp.ProductDetail,
					SubscriptionType:      temp.SubscriptionType,
					PaymentTime:           temp.PaymentTime,
					PretaxGrossAmount:     moneys[j],
					DeductedByCoupons:     temp.DeductedByCoupons,
					InvoiceDiscount:       temp.InvoiceDiscount,
					PretaxAmount:          moneys[j],
					PaymentAmount:         moneys[j],
					OutstandingAmount:     temp.OutstandingAmount,
					Currency:              temp.Currency,
					DeductedByPrepaidCard: temp.DeductedByPrepaidCard,
					DeductedByCashCoupons: temp.DeductedByCashCoupons,
					OwnerId:               temp.OwnerId,
					AccountID:             temp.AccountID,
					PipCode:               temp.PipCode,
					InstanceSpec:          temp.InstanceSpec,
					InstanceId:            temp.InstanceId,
					InternetIp:            temp.InternetIp,
					IntranetIp:            temp.IntranetIp,
					ListPrice:             temp.ListPrice,
					ListPriceUnit:         temp.ListPriceUnit,
					NickName:              temp.NickName,
					Region:                temp.Region,
					ServicePeriod:         temp.ServicePeriod,
					ServicePeriodUnit:     temp.ServicePeriodUnit,
					UsageUnit:             temp.UsageUnit,
					Zone:                  temp.Zone,
					Deleted:               temp.Deleted,
					TenantId:              temp.TenantId,
				})
			}
		case "月":
			month, _ := strconv.Atoi(temp.ServicePeriod)
			moneys, months := subByMonth(temp.PaymentTime, month, temp.PaymentAmount)
			appendLens += len(moneys)
			for j := range months {
				tempInserts = append(tempInserts, &AliyunBillWashingQuery{
					ItemUuid:              temp.ItemUuid,
					BillAccountId:         temp.BillAccountId,
					BillingCycle:          months[j],
					ProductCode:           temp.ProductCode,
					ProductName:           temp.ProductName,
					ProductType:           temp.ProductType,
					ProductDetail:         temp.ProductDetail,
					SubscriptionType:      temp.SubscriptionType,
					PaymentTime:           temp.PaymentTime,
					PretaxGrossAmount:     moneys[j],
					DeductedByCoupons:     temp.DeductedByCoupons,
					InvoiceDiscount:       temp.InvoiceDiscount,
					PretaxAmount:          moneys[j],
					PaymentAmount:         moneys[j],
					OutstandingAmount:     temp.OutstandingAmount,
					Currency:              temp.Currency,
					DeductedByPrepaidCard: temp.DeductedByPrepaidCard,
					DeductedByCashCoupons: temp.DeductedByCashCoupons,
					OwnerId:               temp.OwnerId,
					AccountID:             temp.AccountID,
					PipCode:               temp.PipCode,
					InstanceSpec:          temp.InstanceSpec,
					InstanceId:            temp.InstanceId,
					InternetIp:            temp.InternetIp,
					IntranetIp:            temp.IntranetIp,
					ListPrice:             temp.ListPrice,
					ListPriceUnit:         temp.ListPriceUnit,
					NickName:              temp.NickName,
					Region:                temp.Region,
					ServicePeriod:         temp.ServicePeriod,
					ServicePeriodUnit:     temp.ServicePeriodUnit,
					UsageUnit:             temp.UsageUnit,
					Zone:                  temp.Zone,
					Deleted:               temp.Deleted,
					TenantId:              temp.TenantId,
				})
			}
		default:
			i++
			continue
		}
		if i == 0 {
			saveWashingGroup = saveWashingGroup[1:]
		} else if i == len(saveWashingGroup)-1 {
			saveWashingGroup = saveWashingGroup[:i]
		} else {
			saveWashingGroup = append(saveWashingGroup[:i], saveWashingGroup[i+1:]...)
		}
	}
	saveWashingGroup = append(saveWashingGroup, tempInserts...)
	for _, query := range saveWashingGroup {
		if query.ServicePeriodUnit == "年" || query.ServicePeriodUnit == "月" && query.ServicePeriod == "3" && (query.BillingCycle == "2022-10" || query.BillingCycle == "2022-11") {
			log.Println(query.InstanceId, query.PaymentTime, query.BillingCycle, query.PaymentAmount)
		}
	}
	//moneys, months := subByMonth("2022-10-01", 2, 560)
	////moneys, months := subByyear("2021-01-25", 2, 1200)
	//for i := 0; i < len(moneys); i++ {
	//	log.Println("账期:", months[i], "金额:", moneys[i])
	//}

}

//day := remainday(1, 3, "2022-05-15")
////_, all := SubData("2022-08-14", "2022-05-15")
//m, _ := time.ParseInLocation("2006-01-02", "2022-05-15", time.Local)
//date := m.AddDate(0, 1, 0)
//fmt.Println("date: ", date)
////fmt.Println("day: ", all)
//fmt.Println("day: ", day)

//func AliInitCsvToDbBysubscription(accountId int64, enrollNum, curMonth, fileName string, tenantID, taskID int64, data [][]string, seed int, cpp, sub float64) {
//			if data[i][alitype.ServicePeriodUnit] == "年" {
//
//			} else if data[i][alitype.ServicePeriodUnit] == "月" {
//				if data[i][alitype.ServicePeriod] != "1" {
//					period, _ := strconv.Atoi(data[i][alitype.ServicePeriod])
//					m, _ := time.ParseInLocation("2006-01", curMonth, time.Local)
//					d, _ := time.ParseInLocation("2006-01-02", data[i][alitype.BillingDate], time.Local)
//					//总天数
//					date := m.AddDate(0, period, -1)
//					_, all := SubData(data[i][alitype.BillingDate], date)
//					//单价
//					dayPrice, _ := decimal.NewFromFloat(pretaxAmount).Div(decimal.NewFromFloat(float64(all))).Float64()
//					paymentPrice, _ := decimal.NewFromFloat(paymentAmount).Div(decimal.NewFromFloat(float64(all))).Float64()
//					outstandingPrice, _ := decimal.NewFromFloat(outstandingAmount).Div(decimal.NewFromFloat(float64(all))).Float64()
//
//					for i = 0; i < period; i++ {
//						day := RemainDay(i, period, data[i][alitype.BillingDate])
//						date := m.AddDate(0, i, 0)
//						idate := d.AddDate(0, i, 0)
//						divPretaxAmount, _ := decimal.NewFromFloat(dayPrice).Mul(decimal.NewFromFloat(float64(day))).Float64()
//						divPaymentAmount, _ := decimal.NewFromFloat(paymentPrice).Mul(decimal.NewFromFloat(float64(day))).Float64()
//						divOutstandingAmount, _ := decimal.NewFromFloat(outstandingPrice).Mul(decimal.NewFromFloat(float64(day))).Float64()
//						sMonth := time.Time.String(date)
//						save(uid, accountId, sMonth, data[i], pretaxGrossAmount, deductedByCoupons, invoiceDiscoun, divPretaxAmount, divPaymentAmount, divOutstandingAmount, deductedByPrepaidCard, deductedByCashCoupons, tenantID)
//						//月份在现在月份之后
//						if date.After(time.Now()) {
//							save(uid, accountId, sMonth, data[i], pretaxGrossAmount, deductedByCoupons, invoiceDiscoun, 0, 0, 0, deductedByPrepaidCard, deductedByCashCoupons, tenantID)
//						} else if idate.Before(time.Now()) {
//							fd := d.AddDate(0, 0, -idate.Day()+1)
//							_, a := SubData(fd, idate)
//							dPretaxAmount, _ := decimal.NewFromFloat(dayPrice).Mul(decimal.NewFromFloat(float64(a))).Float64()
//							dPaymentAmount, _ := decimal.NewFromFloat(paymentPrice).Mul(decimal.NewFromFloat(float64(a))).Float64()
//							dOutstandingAmount, _ := decimal.NewFromFloat(outstandingPrice).Mul(decimal.NewFromFloat(float64(a))).Float64()
//							save(uid, accountId, sMonth, data[i], pretaxGrossAmount, deductedByCoupons, invoiceDiscoun, dPretaxAmount, dPaymentAmount, dOutstandingAmount, deductedByPrepaidCard, deductedByCashCoupons, tenantID)
//						}
//
//					}
//
//				}
//			} else {
//				save(uid, accountId, curMonth, data[i], pretaxGrossAmount, deductedByCoupons, invoiceDiscoun, pretaxAmount, paymentAmount, outstandingAmount, deductedByPrepaidCard, deductedByCashCoupons, tenantID)
//			}
//
//			//}
//
//			idx++
//		}
//
//		if len(saveWashingGroup) > 0 {
//			wg.Add(1)
//			go func(saveWashingGroup []*billwashingtype.AliyunBillWashingQuery) {
//				plog.Info(overviewtype.TID, "TransMasterByMonth", "TransMasterByMonth_curMont:%s;length:%d \n---------------------TransMasterByMonth-------------addingToDb washing_query-------------------", curMonth, len(saveWashingGroup))
//				if len(saveWashingGroup) > 0 {
//					_, err = billwashingrdb.SaveBillWashing(saveWashingGroup, "t_inf_bill_aliyun_washing_query")
//					if err != nil {
//						plog.Error("", "aliyun-保存当前月份清洗数据失败", err)
//					}
//				}
//				wg.Done()
//			}(saveWashingGroup)
//		}
//
//	}
//	wg.Wait()
//	plog.Info(overviewtype.TID, "TransMasterByMonth", "TransMasterByMonth_curMont:%s;length:%d \n---------------------TransMasterByMonth-------------end-------------------", curMonth, length)
//}
func remainday(recent, end int, date string) (day int) {
	d, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	fmt.Println("today: ", d)
	firstday := d.AddDate(0, 0, -d.Day()+1)
	fmt.Println("firstday: ", firstday)
	finalday := d.AddDate(0, 1, -d.Day())
	fmt.Println("final day: ", finalday)
	switch recent {
	case 0: //第一个月
		sub := finalday.Sub(d)
		day = int(sub.Hours() / 24)
	case end - 1: //最后一个月
		sub := d.Sub(firstday)
		day = int(sub.Hours() / 24)
	default: //中间
		sub := finalday.Sub(firstday)
		day = int(sub.Hours() / 24)
	}
	return day
}

func GetData() []*AliyunBillWashingQuery {
	jsonFile, err := os.Open("C:\\Users\\Helen.Wang\\Desktop\\tool_attr.json")
	if err != nil {
		fmt.Println("error", err)
	}
	defer jsonFile.Close()
	//file, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		fmt.Println("read file failed", err)
	}
	//if len(file) == 0 {
	//	fmt.Println("unmarshal file failed", err)
	//
	//}
	var ali []*AliyunBillWashingQuery
	//var month int
	//err = json.Unmarshal(file, &ali)
	err = json.NewDecoder(jsonFile).Decode(&ali)
	if err != nil {
		fmt.Println("unmarshal file failed", err)
	}
	return ali
}

func subByMonth(startTime string, month int, money float64) (moneys []float64, months []string) {
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	monthNow := parse.Month()
	year := parse.Year()
	day := parse.Day()
	allDay := (parse.AddDate(0, month, 0).Sub(parse).Hours()) / 24
	oneDay := money / allDay
	if (time.Now().Sub(parse).Hours())/24 < allDay {
		monthForNow := time.Now().Format("2006-01")
		parseNow := parse.Format("2006-01")
		for {
			if monthForNow == parseNow {
				moneys = append(moneys, getMouthStartDay(time.Now())*oneDay)
				months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
				break
			} else {
				moneys = append(moneys, getMouthDay(parse)*oneDay)
				months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
				monthNow++
				if monthNow > 12 {
					year++
					monthNow %= 12
				}
				parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
				parse = parse.AddDate(0, 1, 0)
			}
			parseNow = parse.Format("2006-01")
		}
	} else {
		for i := 0; i < month; i++ {
			moneys = append(moneys, getMouthDay(parse)*oneDay)
			months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
			monthNow++
			if monthNow > 12 {
				year++
				monthNow %= 12
			}
			parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
			parse = parse.AddDate(0, 1, 0)
		}
		moneys = append(moneys, float64(day-1)*oneDay)
		months = append(months, strconv.Itoa(year)+"-"+returnStrMonth(int(monthNow)))
	}
	return
}
func subByyear(startTime string, year int, money float64) (moneys []float64, months []string) {
	allmonth := year * 12
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	monthNow := parse.Month()
	yearNow := parse.Year()
	day := parse.Day()
	allDay := (parse.AddDate(year, 0, 0).Sub(parse).Hours()) / 24
	oneDay := money / allDay
	if (time.Now().Sub(parse).Hours())/24 < allDay {
		monthForNow := time.Now().Format("2006-01")
		parseNow := parse.Format("2006-01")
		for {
			if monthForNow == parseNow {
				moneys = append(moneys, getMouthStartDay(time.Now())*oneDay)
				//moneys = append(moneys, time.Now().Sub(parse).Hours()/24*oneDay)
				months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
				break
			} else {
				moneys = append(moneys, getMouthDay(parse)*oneDay)
				//moneys = append(moneys, time.Now().Sub(parse).Hours()/24*oneDay)
				months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
				monthNow++
				if monthNow > 12 {
					yearNow++
					monthNow %= 12
				}
				parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
				parse = parse.AddDate(0, 1, 0)
			}
			parseNow = parse.Format("2006-01")
		}
	} else {
		for i := 0; i < allmonth; i++ {
			moneys = append(moneys, getMouthDay(parse)*oneDay)
			months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
			monthNow++
			if monthNow > 12 {
				yearNow++
				monthNow %= 12
			}
			parse, _ = time.Parse("2006-01-02", parse.Format("2006-01")+"-01")
			parse = parse.AddDate(0, 1, 0)
		}
		moneys = append(moneys, float64(day-1)*oneDay)
		months = append(months, strconv.Itoa(yearNow)+"-"+returnStrMonth(int(monthNow)))
	}
	return
}
func subByYear(startTime string, year int, money float64) (moneys []float64, years []string) {
	//allmonth := year * 12
	parse, err := time.Parse("2006-01-02", startTime)
	if err != nil {
		panic(err)
	}
	allDay := (parse.AddDate(year, 0, 0).Sub(parse).Hours()) / 24
	oneDay := money / allDay
	for year > 0 {
		if parse.AddDate(1, 0, 0).Sub(time.Now()) > 0 {
			moneys = append(moneys, time.Now().Sub(parse).Hours()/24*oneDay)
			years = append(years, parse.Format("2006-01"))
			break
		} else {
			if parse.Year()%4 == 0 {
				moneys = append(moneys, 366*oneDay)
			} else {
				moneys = append(moneys, 365*oneDay)
			}
			years = append(years, parse.Format("2006-01"))
			year--
			parse = parse.AddDate(1, 0, 0)
		}
	}
	return
}
func getMouthDay(t time.Time) (day float64) {
	parse, _ := time.Parse("2006-01-02", t.Format("2006-01")+"-01")
	return math.Ceil(parse.AddDate(0, 1, -1).Sub(t).Hours()/24) + 1
}

func getMouthStartDay(t time.Time) (day float64) {
	parse, _ := time.Parse("2006-01-02", t.Format("2006-01")+"-01")
	return math.Ceil(t.Sub(parse).Hours() / 24)
}

func returnStrMonth(i int) string {
	if i < 10 {
		return "0" + strconv.Itoa(i)
	} else {
		return strconv.Itoa(i)
	}
}

func RemainDay(recent, end int, date string) (day int) {
	d, _ := time.ParseInLocation("2006-01-02", date, time.Local)
	fmt.Println("today: ", d)
	firstday := d.AddDate(0, 0, -d.Day()+1)
	fmt.Println("firstday: ", firstday)
	finalday := d.AddDate(0, 1, -d.Day())
	fmt.Println("final day: ", finalday)
	switch recent {
	case 0: //第一个月
		sub := finalday.Sub(d)
		day = int(sub.Hours() / 24)
	case end - 1: //最后一个月
		sub := d.Sub(firstday)
		day = int(sub.Hours() / 24)
	default: //中间
		sub := finalday.Sub(firstday)
		day = int(sub.Hours() / 24)
	}
	return day
}

func SubData(data ...interface{}) (timeSub []int, all int) {
	var timeTable []time.Time
	for i := range data {
		switch data[i].(type) {
		case string:
			data[i], _ = time.Parse("2006-01-02", data[i].(string))
		}
		timeTable = append(timeTable, data[i].(time.Time))
		if i > 0 {
			timeSub = append(timeSub, int(math.Ceil(timeTable[i].Sub(timeTable[i-1]).Hours()))/24)
		}
	}
	return timeSub, int(math.Ceil(timeTable[len(timeTable)-1].Sub(timeTable[0]).Hours())) / 24
}
