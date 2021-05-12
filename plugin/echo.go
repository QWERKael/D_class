package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"fmt"
	"time"
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
desc: "Echo Plugin\n简单的echo插件, 常用来进行各种测试\n
echo [string] \n\t将输入的字符串原封不动的返回回来\n
echo now \n\t返回现在的服务器时间\n
echo sleep [sec] \n\t休眠 sec 秒, 常用于测试各种异步操作\n"
yess:
  echo:
    type: "Plugin"
    desc: "提供echo相关操作"
    yess:
      now:
        type: "Cmd"
        desc: "返回现在的服务器时间"
      sleep:
        type: "Cmd"
        desc: "休眠指定的时间, 单位: 秒"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

func DefaultFunc(task tasks.Task) (string, error) {
	log.Debugf("接收到信息: %s", task.Cmd)
	return task.Cmd, nil
}

func Now(task tasks.Task) (string, error) {
	return time.Now().String(), nil
}

func Sleep(task tasks.Task) (string, error) {
	sec := utils.StringDefaultInt64(task.SubCmd[0], 10)
	d := time.Duration(sec * int64(time.Second))
	time.Sleep(d)
	return fmt.Sprintf("已休眠 %d 秒", sec), nil
}

func Parse(task tasks.Task) (*pb.CommonCmdReply, error) {
	ccr, _ := prompt.Parse(task.SubCmd[0])
	task = tasks.Task{Plugin: ccr.Plugin, Cmd: ccr.Cmd, SubCmd: ccr.SubCmd,
		Flags: ccr.Flags, Args: ccr.Args}
	body := [][]string{
		{"Plugin", task.Plugin},
		{"Cmd", task.Cmd},
	}
	for _, sub := range task.SubCmd {
		body = append(body, []string{"SubCmd", sub})
	}
	for _, flag := range task.Flags {
		body = append(body, []string{"-" + flag, "true"})
	}
	for k, v := range task.Args {
		body = append(body, []string{"--" + k, v})
	}
	reply := prompt.ToTable([]string{"Item", "Value"}, []string{fmt.Sprintf("总计: %d 行", len(body))}, body, 0)
	return reply, nil
}

func Echoloop(task tasks.Task) (string, *tasks.Task, error) {
	s := task.SubCmd[0]
	log.Debugf("接收到信息: %s", s)
	if len(s) <= 1 {
		return s, nil, nil
	}
	task.SubCmd[0] = s[1:]
	return s, &task, nil
}
