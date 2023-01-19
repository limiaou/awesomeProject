package gettoken

import (
	"awesomeProject/gettoken/structs"
	"crypto/tls"
	"encoding/json"
	tools "github.com/iEvan-lhr/exciting-tool"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetTokenByDefault() string {
	//TODO 改为从consul读取配置信息
	addr := "https://cmp.cn-lghnh.com:9800/backend/service-s-authentication/v1/usermgmt/login"
	userName := "super"
	password := "Connext@1qaz@WSX"
	js := "{\"user_name\":\"{{Name}}\",\"passwd\":\"{{Password}}\"}"
	js = strings.Replace(js, "{{Name}}", userName, -1)
	js = strings.Replace(js, "{{Password}}", password, -1)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	request, err := http.NewRequest("POST", addr, strings.NewReader(js))
	request.Header = HeaderPublic()
	if err != nil {
		panic(err)
	}
	response := tools.ReturnValue(client.Do(request)).(*http.Response)
	all := tools.ReturnValue(io.ReadAll(response.Body)).([]byte)
	pubRes := structs.PublicResponse{}
	tools.Error(json.Unmarshal(all, &pubRes))
	time.Sleep(1 * time.Second)
	return pubRes.Data.Token
}

func HeaderPublic() http.Header {
	header := http.Header{}
	header.Set("Accept", "*/*")
	header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	header.Set("Connection", "keep-alive")
	header.Set("Content-Type", "application/json")
	header.Set("User-Agent", "PostmanRuntime/7.28.4")
	return header
}
