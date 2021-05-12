package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"crypto/tls"
	"strings"

	//"atk_D_class/utils"
	"errors"
	"gopkg.in/gomail.v2"
)

var log = logger.Logger

var (
	BuildVersion string
	BuildTime    string
	BuildName    string
)

func Version(task tasks.Task) (*pb.CommonCmdReply, error) {
	body := [][]string{
		{"BuildName", BuildName},
		{"BuildVersion", BuildVersion},
		{"BuildTime", BuildTime},
	}
	reply := prompt.ToTable([]string{}, []string{}, body, 0)
	return reply, nil
}

func Help(task tasks.Task) (string, error) {
	helpMsg := `type: "Root"
text: "root"
desc: "Email Plugin\n电子邮件插件\n
email send 邮件标题 邮件内容 \n\t发送电子邮件\n"
yess:
  email:
    type: "Plugin"
    desc: "提供电子邮件相关操作"
    yess:
      send:
        type: "Cmd"
        desc: "发送电子邮件"
        yess:
          "--host=":
            type: "ArgKey"
            desc: "指定邮件服务器地址"
          "--port=":
            type: "ArgKey"
            desc: "指定邮件服务器端口"
          "--username=":
            type: "ArgKey"
            desc: "指定邮件服务器登陆用户名"
          "--password=":
            type: "ArgKey"
            desc: "指定邮件服务器登陆密码"
          "--from=":
            type: "ArgKey"
            desc: "指定邮件发件人"
          "--to=":
            type: "ArgKey"
            desc: "指定邮件收件人"`
	return helpMsg, nil
}

func Send(task tasks.Task) (string, error) {
	host := utils.StringDefault(task.Args["host"], "")
	port := utils.StringDefaultInt(task.Args["port"], 25)
	username := utils.StringDefault(task.Args["username"], "")
	password := utils.StringDefault(task.Args["password"], "")

	from := utils.StringDefault(task.Args["from"], "")
	to := utils.StringDefault(task.Args["to"], "")
	subject := utils.StringDefault(task.SubCmd[0], "")
	body := utils.StringDefault(strings.Join(task.SubCmd[1:], "\n"), "")

	log.Debugf("[%s]用户准备通过服务器[%s:%d(%s)]发送标题为[%s]的邮件到[%s]", from, host, port, username, subject, to)

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		return "", errors.New("发送失败: " + err.Error())
	}
	return "发送成功", nil
}
