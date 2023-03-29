package main

import (
	"awesomeProject/autosave/structs"
	"awesomeProject/gettoken"
	"crypto/tls"
	"encoding/json"
	tools "github.com/iEvan-lhr/exciting-tool"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetResourceID() {
	//token := gettoken.GetTokenByDefault()
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjAsImlhdCI6MTY3OTg5NzAzOSwianRpIjoiMTY3OTg5NzAzOTEiLCJ1c2VyX2luZm8iOnsiVXNlcklEIjoxLCJVc2VyVHlwZSI6ImFkbWluIiwiVXNlck5hbWUiOiJzdXBlciIsIk5pY2tOYW1lIjoic3VwZXIiLCJMb2dpblRpbWUiOiIiLCJUZW5hbnRJRCI6MH19.CPf1WxP12754JTtoU6Yaf1mlmXNWK9n0lsIxWWP-DBQ"
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("POST",
		//"https://cmp.cn-lghnh.com:9800/backend/service-b-cmdb/v2/cmdbresourcev2/listresource",
		"https://cmp.pernod-ricard-china.com/backend/service-b-cmdb/v2/cmdbresourcev2/listresource",
		strings.NewReader(`{
		   "text_query_condi": {
        	"account_id": "3"
    	},
		"select_query_condi": {
			"cloud_type": [
				"2"
			],
			"allocate_state": [
				"0"
			],
			"is_invalid": [
				"0"
			]
		   },
		   "page_size": 170,
		   "page_num": 1
		}`))
	//		strings.NewReader(`{
	//    "text_query_condi": {
	//        "account_id": "2"
	//"?fuzzy_search": "Premium SSD Managed Disks - P10 - LRS - Disk - CN East 2-Azure"
	//    },
	//    "select_query_condi": {
	//"allocate_state": [
	//"0"
	//],
	//        "cloud_type": [
	//            "2"
	//        ],
	//        "is_invalid": [
	//            "0"
	//        ]
	//    },
	//    "page_size": 400,
	//    "page_num": 1
	//}`))
	request.Header = gettoken.HeaderPublic()
	request.Header.Set("Authorization", token)
	if err != nil {
		panic(err)
	}
	response := tools.ReturnValue(client.Do(request)).(*http.Response)
	all := tools.ReturnValue(io.ReadAll(response.Body)).([]byte)
	resp := &structs.ResponseResource{}
	tools.Error(json.Unmarshal(all, resp))
	//uat
	save := structs.Save{
		//SaveList: []structs.SaveList{{EnvironmentId: 458450667091008, AppId: 458449754726400}}, //aws landing zone
		SaveList: []structs.SaveList{{EnvironmentId: 469856822846528, AppId: 469780427239424}}, //Magento
		//SaveList: []structs.SaveList{{EnvironmentId: 469856822846528, AppId: 469780427239424}},
	}
	//prd
	//save := structs.Save{
	//	SaveList: []structs.SaveList{{EnvironmentId: 394101379982400, AppId: 394099990945792}},
	//}
	for i := range resp.Data.ResourceList {
		save.SaveList[0].ResourceId = append(save.SaveList[0].ResourceId, resp.Data.ResourceList[i].ResourceId)
	}
	marshal := tools.ReturnValue(json.Marshal(&save)).([]byte)
	tools.Error(os.WriteFile("data1.json", marshal, 0666))
}
