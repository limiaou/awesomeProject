package main

import (
	"github.com/tebeka/selenium"
	"log"
	"net/http/cookiejar"
)

const (
	seleniumPath = `C:\Program Files\Google\Chrome\Application\chromedriver.exe` //设置chromedriver在电脑磁盘的位置
	port         = 9515
)

func main() {
	ops := []selenium.ServiceOption{}
	service, _ := selenium.NewChromeDriverService(seleniumPath, port, ops...)
	defer service.Stop()
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	cookiejar.New(nil)
	wd, _ := selenium.NewRemote(caps, "http://127.0.0.1:9515/wd/hub")
	//defer wd.Quit()
	//time.Sleep(time.Second * 5)
	//wd.Get("https://login.taobao.com/")
	//time.Sleep(time.Second * 5)
	//wd.Get("https://cart.taobao.com/cart.htm")
	wd.Get("https://detail.tmall.com/item.htm?id=680718952144")
	//wd.Get("https://detail.tmall.com/item.htm?id=681586638199")
	//time.Sleep(time.Second * 5)
	//header := http.Header{}
	//header.Set("data", "[{\"shopId\":\"s_2382759536\",\"comboId\":0,\"shopActId\":0,\"cart\":[{\"quantity\":1,\"cartId\":\"2880795816239\",\"skuId\":\"4864346461422\",\"itemId\":\"655687317692\"}],\"operate\":[],\"type\":\"check\"}]")
	//wd.Get("https://buy.taobao.com/auction/order/confirm_order.htm?spm=a1z0d.6639537.0.0.undefined")
	//time.Sleep(time.Second * 5)
	we, _ := wd.FindElement(selenium.ByID, "J_LinkBuy") //找到搜索一下的id
	err := we.Click()
	log.Println(we, err) // 点击
	//time.Sleep(1 * time.Second)
	for {
		we1, err := wd.FindElement(selenium.ByClassName, "go-btn")
		if err != nil {
			log.Println(err)
			//time.Sleep(20 * time.Second)
		} else {
			we1.Click()
			//log.Println(we.TagName())
			break
		}
	}

	//for {
	//	we1, err := wd.FindElement(selenium.ByID, "J_Go")
	//	if err != nil {
	//		log.Println(err)
	//		time.Sleep(20 * time.Second)
	//	} else {
	//		we1.Click()
	//		time.Sleep(3 * time.Second)
	//		for {
	//			we1, err = wd.FindElement(selenium.ByClassName, "go-btn")
	//			if err != nil {
	//				log.Println(err)
	//				time.Sleep(20 * time.Second)
	//			} else {
	//				log.Println(we.TagName())
	//				break
	//			}
	//		}
	//		break
	//	}
	//
	//}

	//time.Sleep(time.Second * 100)
}
