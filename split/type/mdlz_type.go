package typeMNLZ

type TempShare struct {
	AllShare        float64
	TotalMoney      float64
	TotalSubScript  float64
	TotalPayAsYouGo float64
}

func (t *TempShare) SaveMoney(overview *BillOverviewShare) *TempShare {
	switch true {
	case overview.AllShare > 0:
		t.AllShare = overview.AllShare
	case overview.SubscriptionCost > 0:
		t.TotalSubScript = overview.SubscriptionCost
		t.TotalMoney += overview.SubscriptionCost
	case overview.PayAsYouGoCost > 0:
		t.TotalPayAsYouGo = overview.PayAsYouGoCost
		t.TotalMoney += overview.PayAsYouGoCost
	}
	return t
}

type ByProduct struct {
	CloudRoomName string
	ProductName   string
	subscription  string
}

type AliyunBillSummaryDay struct {
	SummaryId      int64  `orm:"summary_id" json:"summary_id" desc:"汇总ID"`
	ProductId      int64  `orm:"product_id" json:"product_id" desc:"产品ID"`
	ResourceId     int64  `orm:"resource_id" json:"resource_id" desc:"资源ID"`
	AccountId      int64  `orm:"account_id" json:"account_id" desc:"云账号ID"`
	CloudAccountId string `orm:"cloud_account_id" json:"cloud_account_id" desc:"云账号ID"`

	Category    string `orm:"category" json:"category" desc:"类别"`
	Application string `orm:"application" json:"application" desc:"应用"`
	Environment string `orm:"environment" json:"environment" desc:"环境"`
	Department  string `orm:"department" json:"department" desc:"部门"`
	Owner       string `orm:"owner" json:"owner" desc:"负责人"`

	Business        string `orm:"business" json:"business" desc:"业务"`
	ProductSource   string `orm:"product_source" json:"product_source" desc:"产品"`
	CloudRoomId     string `orm:"cloud_room_id" json:"cloud_room_id" desc:"云单元（清洗）"`
	ResourceGroupId string `orm:"resource_group_id" json:"resource_group_id" desc:"资源分组id（清洗）"`
	Tags            string `orm:"tags" json:"tags" desc:"标签"`

	BillingCycle string `orm:"billing_cycle" json:"billing_cycle" desc:"计费周期"`
	BillingDate  string `orm:"billing_date" json:"billing_date" desc:"计费日期"`
	BillingWeek  string `orm:"billing_week" json:"billing_week" desc:"周账期"`
	ProductCode  string `orm:"product_code" json:"product_code" desc:"产品代码"`
	ProductName  string `orm:"product_name" json:"product_name" desc:"产品名称"`

	ProductType      string `orm:"product_type" json:"product_type" desc:"产品类型"`
	ProductDetail    string `orm:"product_detail" json:"product_detail" desc:"产品明细"`
	SubscriptionType string `orm:"subscription_type" json:"subscription_type" desc:"订阅类型 Subscription／PayAsYouGo"`
	PaymentTime      string `orm:"payment_time" json:"payment_time" desc:"订单支付时间"`
	ItemType         string `orm:"item_type" json:"item_type" desc:"计费项类型"`

	PretaxGrossAmount float64 `orm:"pretax_gross_amount" json:"pretax_gross_amount" desc:"原始金额,original cost"`
	DeductedByCoupons float64 `orm:"deducted_by_coupons" json:"deducted_by_coupons" desc:"优惠券抵扣"`
	InvoiceDiscount   float64 `orm:"invoice_discount" json:"invoice_discount" desc:"优惠金额"`
	PretaxAmount      float64 `orm:"pretax_amount" json:"pretax_amount" desc:"应付金额"`
	PaymentAmount     float64 `orm:"payment_amount" json:"payment_amount" desc:"现金支付"`

	OutstandingAmount     float64 `orm:"outstanding_amount" json:"outstanding_amount" desc:"未结清金额"`
	DeductedByPrepaidCard float64 `orm:"deducted_by_prepaid_card" json:"deducted_by_prepaid_card" desc:"代金券抵扣"`
	DeductedByCashCoupons float64 `orm:"deducted_by_cash_coupons" json:"deducted_by_cash_coupons" desc:"储值卡抵扣"`
	DeductedByCredit      float64 `orm:"deducted_by_credit" json:"deducted_by_credit" desc:"信用抵扣"`
	Currency              string  `orm:"currency" json:"currency" desc:"币种 CNY／USD／JPY"`

	Status        string `orm:"status" json:"status" desc:"支付状态:已支付（Cleared），未结清（Uncleared），未结算（Unsettled），免结算"`
	ResourceGroup string `orm:"resource_group" json:"resource_group" desc:"资源组"`
	PipCode       string `orm:"pip_code" json:"pip_code" desc:"产品Code"`
	InstanceId    string `orm:"instance_id" json:"instance_id" desc:"实例ID"`
	NickName      string `orm:"nick_name" json:"nick_name" desc:"实例昵称"`

	Region            string `orm:"region" json:"region" desc:"地域"`
	ServicePeriod     string `orm:"service_period" json:"service_period" desc:"服务时长"`
	ServicePeriodUnit string `orm:"service_period_unit" json:"service_period_unit" desc:"服务时长单位"`
	UsageNumber       string `orm:"usage_number" json:"usage_number" desc:"用量"`
	UsageUnit         string `orm:"usage_unit" json:"usage_unit" desc:"用量单位, 仅当isBillingItem为true时有效"`

	Zone        string `orm:"zone" json:"zone" desc:"可用区"`
	CostUnit    string `orm:"cost_unit" json:"cost_unit" desc:""`
	BillingType string `orm:"billing_type" json:"billing_type" desc:""`
	BizType     string `orm:"biz_type" json:"biz_type" desc:""`
	PaymentNo   string `orm:"payment_no" json:"payment_no" desc:""`

	PayType     string `orm:"pay_type" json:"pay_type" desc:""`
	CreatedBy   string `orm:"created_by" json:"created_by" desc:"创建人"`
	CreatedTime string `orm:"created_time" json:"created_time" desc:"创建时间"`
	UpdatedBy   string `orm:"updated_by" json:"updated_by" desc:"更新人"`
	UpdatedTime string `orm:"updated_time" json:"updated_time" desc:"更新时间"`

	Enabled          int8   `orm:"enabled" json:"enabled" desc:"有效标志 1-无效,0-有效"`
	Deleted          int8   `orm:"deleted" json:"deleted" desc:"删除标志 1-删除,0-正常"`
	TenantId         int64  `orm:"tenant_id" json:"tenant_id" desc:"租户ID"`
	RecordId         string `orm:"record_id" json:"record_id" desc:"订单号/账单号"`
	CloudAccountName string `orm:"cloud_account_name" json:"cloud_account_name" desc:"云账号名称"`

	ResourceGroupName string `orm:"resource_group_name" json:"resource_group_name" desc:"资源组名称"`
	CloudRoomName     string `orm:"cloud_room_name" json:"cloud_room_name" desc:"云资源名称"`
	ProductTypeCnName string `orm:"product_type_cn_name" json:"product_type_cn_name" desc:"资源类型名称"`
	GroupCnName       string `orm:"group_cn_name" json:"group_cn_name" desc:"资源分类名称"`
	ProductTypeEnName string `orm:"product_type_en_name" json:"product_type_en_name" desc:"资源类型英文名称"`
	GroupEnName       string `orm:"group_en_name" json:"group_en_name" desc:"资源分类英文名称"`

	UnionRGID string `orm:"union_rgid"  json:"union_rgid" desc:"rgunion主键"`
	AppId     int64  `orm:"app_id" json:"app_id" desc:"应用ID"`
}

type TInfBillOverviewSummaryMonthItem struct {
	Uuid             int64   `json:"uuid"`
	AccountId        int64   `json:"account_id"`
	AccountName      string  `json:"account_name"`
	BillingCycle     string  `json:"billing_cycle"`
	TotalCost        float64 `json:"total_cost"`
	AppId            int64   `json:"app_id"`
	AppName          string  `json:"app_name"`
	CloudRoomId      int64   `json:"cloud_room_id"`
	CloudRoomName    string  `json:"cloud_room_name"`
	SubscriptionType int64   `json:"subscription_type"`
	CloudType        int64   `json:"cloud_type"`
	BillingType      int64   `json:"billing_type"`
}

type ShareByProduct struct {
	ShareMoney       float64
	SubscriptionType string
	ProductName      string
}

type BillOverviewShare struct {
	Uuid         int64  `orm:"uuid" json:"uuid" desc:"汇总ID"`
	AccountId    int64  `orm:"account_id" json:"account_id" desc:"云账号ID"`
	AccountName  string `orm:"account_name" json:"account_name" desc:"云账号名称"`
	BillingCycle string `orm:"billing_cycle" json:"billing_cycle" desc:"计费周期-月"`
	AppId        int64  `orm:"app_id" json:"app_id" desc:"应用id"`

	AppName           string  `orm:"app_name" json:"app_name" desc:"应用名称"`
	CloudRoomId       int64   `orm:"cloud_room_id" json:"cloud_room_id" desc:"云单元id"`
	CloudRoomName     string  `orm:"cloud_room_name" json:"cloud_room_name" desc:"云单元名称"`
	TotalCost         float64 `orm:"total_cost" json:"total_cost" desc:"费用"`
	UpdatedTime       string  `orm:"updated_time" json:"updated_time" desc:"更新时间"`
	SubscriptionType  int64   `orm:"subscription_type" json:"subscription_type" desc:"订阅类型：0：包年包月 1:按量付费"`
	CloudType         int64   `orm:"cloud_type" json:"cloud_type" desc:"云类型: 0:aliyun 1:azure 2:aws"`
	BillingType       int64   `orm:"billing_type" json:"billing_type" desc:"费用类型：0:账单 1:cpp"`
	TenantId          int64   `orm:"tenant_id" json:"tenant_id" desc:"租户ID"`
	ProductTypeCnName string  `orm:"product_type_cn_name" json:"product_type_cn_name" desc:"资源类型名称"`
	GroupCnName       string  `orm:"group_cn_name" json:"group_cn_name" desc:"资源分类名称"`
	ProductTypeEnName string  `orm:"product_type_en_name" json:"product_type_en_name" desc:"资源类型英文名称"`
	GroupEnName       string  `orm:"group_en_name" json:"group_en_name" desc:"资源分类英文名称"`

	SubscriptionCost float64 `orm:"subscription_cost" json:"subscription_cost" desc:"订阅类型：0：包年包月"`
	PayAsYouGoCost   float64 `orm:"pay_as_you_go_cost" json:"pay_as_you_go_cost" desc:"订阅类型：1:按量付费"`
	AllShare         float64 `orm:"all_share" json:"all_share" desc:"订阅类型：2:共享资源"`
	AllSubscription  float64 `orm:"all_subscription_cost" json:"all_subscription_cost" desc:"包年包月（分摊后）"`
	AllPayAsYouGo    float64 `orm:"all_pay_as_you_go_cost" json:"all_pay_as_you_go_cost" desc:"按量付费（分摊后）"`
}

const BEGINNING_OF_THE_YEAR = "2022-01-01"
