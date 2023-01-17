// @Copyright(c)，Shanghai Lianwei Technology Co., Ltd.，All Rights Resevered.:
// @Package: main
// @Author: Emory Du
// @Description: ${TODO}
// @Date: 2022-08-09 22:19:13
package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
)

// 联蔚科技运维中心告警邮件发送
func main() {
	//message, err := send.ReadExcelMsg()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//send.BatchSendEmail(message)
	//r := gomail.Rename("100086")
	m := gomail.NewMessage()
	m.SetHeader("From", "jingyiwang8@163.com")
	m.SetHeader("To", "1287137114@qq.com")
	//m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	d := gomail.NewDialer("smtp.163.com", 25, "jingyiwang8@163.com", "xuejian821")
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
	//m.Attach("", r)
	m.Attach("yyq.jpg")
	m.Attach("foo.txt", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write([]byte("Content of foo.txt"))
		fmt.Println(w)
		return err
	}))
}
