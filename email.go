// email

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/smtp"
	"os"
	"runtime/debug"
	"strings"

	//	"github.com/widuu/goini"
)

type cfgmail struct {
	Username   string
	Password   string
	Smtphost   string
	To         string
	University string
}

type cfg struct {
	Name, Text string
}

func init() {

}

// send email
func DOSendMail(to, subject, content string) {
	//从json文件中读取发送邮件服务器配置信息
	cfgjson := getConf()
	var cfg cfgmail
	dec := json.NewDecoder(strings.NewReader(cfgjson))
	for {
		if err := dec.Decode(&cfg); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			debug.PrintStack()
		}
	}

	username := cfg.Username
	password := cfg.Password
	host := cfg.Smtphost
	university := cfg.University
	body := `
    <html>
    <body>
    <h3>
    ` + university + content + `
    </h3>
    </body>
    </html>
    `
	err := SendMail(username, password, host, to, subject, body, "html")
	if err != nil {

		log.Println("send mail error!", err)
		debug.PrintStack()

	} else {
		log.Println("send mail success!")
	}

}

/*
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	//tel version  sendemail  only return nil

	conf := goini.SetConfig("/sendto/DL365/after/conf/conf.ini")
	email := conf.GetValue("version", "email")
	if email == "false" {
		return nil
	}

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
*/
func getConf() string {
	//filename := "/sendto/captor/conf/email.json"

	filename := "/Users/zhangyingchuang/Desktop/captor/conf/email.json"
	file, err := os.Open(filename)

	defer file.Close()
	if err != nil {
		fmt.Println("read conf file error")
		log.Fatal(err)
	}

	buf := make([]byte, 512)
	var str1 string
	for {
		n, _ := file.Read(buf)
		if 0 == n {
			break
		}
		//os.Stdout.Write(buf[:n])
		str := string(buf[:n])

		str1 = str1 + str
	}
	return str1
}

//ye
func SendMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	fmt.Println(msg)
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
