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
	token := gettoken.GetTokenByDefault()
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("POST",
		"https://cmp.cn-lghnh.com:9800/backend/service-b-cmdb/v2/cmdbresourcev2/listresource",
		strings.NewReader(`{
    "text_query_condi": {
        "fuzzy_search": "Premium SSD Managed Disks - P10 - LRS - Disk - CN East 2-Azure"
    },
    "select_query_condi": {
        "allocate_state": [
            "0"
        ],
        "is_invalid": [
            "0"
        ]
    },
    "page_size": 550,
    "page_num": 1
}`))
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
		SaveList: []structs.SaveList{{EnvironmentId: 395460795647040, AppId: 394099990945792}},
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
