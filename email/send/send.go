// Package send
// @Copyright(c)，Shanghai Lianwei Technology Co., Ltd.，All Rights Resevered.:
// @Author: Emory Du
// @Description: ${TODO}
// @Date: 2022-08-09 22:22:02
package send

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"strconv"
)

// Q
// 1. 邮件收件人是多个/从哪个地方获取? (./从当前目录下的Excel中固定的列中获取)
// 2. 邮件抄送人是多个/是固定的: 从外部文件获取
// 3. 邮件的内容从哪里获取? (./当前目录下的Excel中获取)

func BatchSendEmail(message []*S) {

	mailConn := map[string]string{
		"user": "tiger.liu@lianwei.com.cn",
		"pass": "Lfh@0207cn", // qq邮箱smtp服务授权码
		"host": "smtp.mxhichina.com",
		"port": "465",
	}
	for _, msg := range message {
		m := gomail.NewMessage()
		m.SetHeader("From", mailConn["user"])
		m.SetHeader("To", msg.SupervisorEmail)
		//m.SetHeader("Cc", "1504454838@qq.com")
		m.SetHeader("Subject", "联蔚运维中心告警提醒")
		template := "Dear: %s<br/>&nbsp;&nbsp;&nbsp;&nbsp;应用 %s 《%s》本周测试环境，服务器系统升级更新，届时会重启服务器。操作时间大约在%s。请知悉。谢谢！<br/><br/>Best Regards!<br/>联蔚盘云-运维中心<br/>" +
			"<img src='' />" + "<br/>5F,Building 2, JuXin Park, No. 188 PingFu Road | 200231 Shanghai | P.R.China<br/><br/>" +
			"<a>www.lianwei.com.cn<a/>"

		message := fmt.Sprintf(template, msg.Supervisor, msg.ApplicationID, msg.ApplicationName, msg.UpgradeTime)
		m.SetBody("text/html", message)
		port, err := strconv.Atoi(mailConn["port"])
		if err != nil {
			fmt.Println(err)
		}
		d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])
		if err := d.DialAndSend(m); err != nil {
			fmt.Println(err)
		}
	}
}

func learn(message *S) {
	r := gomail.Rename("100086")
	m := gomail.NewMessage()
	m.Attach("", r)
	m.Attach("foo.txt", gomail.SetCopyFunc(func(w io.Writer) error {
		_, err := w.Write([]byte("Content of foo.txt"))
		return err
	}))
}
