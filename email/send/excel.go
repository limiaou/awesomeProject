package send

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
)

// S 发送消息结构体
type S struct {
	Supervisor      string `json:"supervisor"`
	ApplicationID   string `json:"application_id"`
	ApplicationName string `json:"application_name"`
	UpgradeTime     string `json:"upgrade_time"`
	SupervisorEmail string `json:"supervisor_email"`
}

// ColNumMap 需要获取的列映射
var ColNumMap = map[string]int{
	"Supervisor":      6,
	"ApplicationID":   0,
	"ApplicationName": 1,
	"UpgradeTime":     3,
	"SupervisorEmail": 8,
}

// ReadExcelMsg 读取Excel获取发送信息
func ReadExcelMsg() ([]*S, error) {
	batchSendInfos := make([]*S, 0)
	file, err := excelize.OpenFile("test.xlsx")
	if err != nil {
		log.Println("ReadExcelMsg: ", err)
		return batchSendInfos, err
	}
	sheetName := file.GetSheetName(1)
	rows := file.GetRows(sheetName)
	for i := 1; i < len(rows); i++ {
		message := &S{}
		for index, cel := range rows[i] {
			if index == ColNumMap["Supervisor"] {
				message.Supervisor = cel
			}
			if index == ColNumMap["ApplicationID"] {
				message.ApplicationID = cel
			}
			if index == ColNumMap["ApplicationName"] {
				message.ApplicationName = cel
			}
			if index == ColNumMap["UpgradeTime"] {
				message.UpgradeTime = cel
			}
			if index == ColNumMap["SupervisorEmail"] {
				message.SupervisorEmail = cel
			}
		}
		batchSendInfos = append(batchSendInfos, message)
	}
	return batchSendInfos, nil
}
