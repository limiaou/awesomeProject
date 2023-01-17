package main

import (
	"fmt"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"image/color"
)

func main() {
	qrc, err := qrcode.New("你好，我是玲娜贝儿，我正在森林探险。现在给我转49.9，待我今天肯德基疯狂星期四购买老北京鸡肉卷补充体力后，今晚，我来你家探险。\n📞贝儿向您发起了探险邀请\n\n🔘接受     丨     🔘想接受")
	//qrc, err := qrcode.New("http://10.255.201.207:32070/")
	if err != nil {
		fmt.Printf("could not generate QRCode: %v", err)
		return
	}

	w, err := standard.New("./qrcode.png", standard.WithHalftone("./mmm.jpg"), standard.WithQRWidth(21), standard.WithBgColor(color.White))
	if err != nil {
		fmt.Printf("standard.New failed: %v", err)
		return
	}

	// save file
	if err = qrc.Save(w); err != nil {
		fmt.Printf("could not save image: %v", err)
	}
}
