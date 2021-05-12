package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/plugin/transfer/transfer"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"net"
	"path/filepath"
	"strings"
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
desc: "Transfer Plugin\n文件传输插件\n
transfer run \n\t启动文件传输监听端口\n
transfer send [file] to [ip:port] \n\t将文件发送到指定的Ip端口\n"
yess:
  transfer:
    type: "Plugin"
    desc: "提供文件传输相关操作"
    yess:
      run:
        type: "Cmd"
        desc: "启动文件传输监听端口"
      send:
        type: "Cmd"
        desc: "将文件发送到指定的Ip端口"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

func Run(task tasks.Task) (string, error) {
	var msg string
	listenAddr := utils.StringDefault(task.Args["listen"], "0.0.0.0:8883")
	recvDir := utils.StringDefault(task.Args["dir"], "files")
	recvDir, err := filepath.Abs(recvDir)
	if err != nil {
		return "工作目录不正确", nil
	}
	recvServ, err := net.Listen("tcp", listenAddr)
	if err != nil {
		msg = "监听端口出错: " + err.Error()
		log.Errorf(msg)
		return msg, nil
	}
	log.Infof("文件传输端口 %s 正在监听中...", listenAddr)
	go transfer.SetReceiver(recvServ, recvDir)
	return "接收器已就绪", nil
}

func Send(task tasks.Task) (string, error) {
	log.Debugf("task信息: %#v", task)
	var err error
	if len(task.SubCmd) < 3 || strings.ToLower(task.SubCmd[1]) != "to" {
		return "命令格式错误", nil
	}
	srcFilePath, err := filepath.Abs(task.SubCmd[0])
	if err != nil {
		return "文件路径不合法", nil
	}
	sendTo := task.SubCmd[2]
	if ft, err := utils.CheckPath(srcFilePath); err != nil {
		return "指定的路径不合法", nil
	} else if ft != utils.File {
		return "指定的路径不是文件", nil
	}
	log.Debugf("将文件 [%s] 发送到 [%s]", srcFilePath, sendTo)
	err = transfer.Send(srcFilePath, sendTo)
	if err != nil {
		return "文件发送失败" + err.Error(), nil
	}
	return "文件发送成功", nil
}
