// all
package main

import (
//"fmt"
//	"log"
)

const (
	//邮件配置信息
	emailuser = "m15210229082@163.com"
	password  = "bmxeudeememsjroh"
	emailhost = "smtp.163.com:25"
	//to := "yeertai@hotmail.com;79968176@qq.com"
	to = "876230915@qq.com"
)

func main() {
	dealDatabase()

	//	if checkNetFile() {
	//		fmt.Println("网络文件目录存在！")
	//		fmt.Println("开始备份数据库")
	//		dealDatabase()
	//	} else {
	//		// send email
	//		fmt.Println("网络文件目录不存在！")
	//		fmt.Println("发送邮件")

	//		subject := "网络文件目录不存在！"
	//		body := "windowns 目录没有挂在到linux机器上"
	//		sendEmail(subject, body)
	//	}

}

//发送邮件
func sendEmail(subject, body string) {

	mailtype := "html"

	SendMail(emailuser, password, emailhost, to, subject, body, mailtype)
}
