package tasks

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"fmt"
	"strings"
)

var log = logger.Logger

type Task struct {
	Plugin  string
	Cmd     string
	SubCmd  []string
	Flags   []string
	Args    map[string]string
	Context map[string]string
}

// 将 Task struct 转换为命令语句
func (t *Task) ToString() string {
	//添加 Plugin Cmd SubCmd 字段
	s := fmt.Sprintf("%s %s %s", t.Plugin, t.Cmd, strings.Join(t.SubCmd, " "))
	//添加 Flags 字段
	if len(t.Flags) > 0 {
		s += " "
		s += "-" + strings.Join(t.Flags, " -")
	}
	//添加 Args 字段
	if len(t.Args) > 0 {
		s += " "
		for k, v := range t.Args {
			s += "--" + k + "=" + v
		}
	}
	return s
}

// 执行 Task
func Exec(t *Task, pm *PluginManager, forcePlugin string, forceCmd string) (*pb.CommonCmdReply, *Task, error) {
	var (
		reply    *pb.CommonCmdReply
		replyMsg string
		nextTask *Task
		err      error
		plugin   string
		cmd      string
	)

	plugin = t.Plugin
	cmd = t.Cmd

	if forcePlugin != "" {
		log.Debugln("强制使用插件", forcePlugin)
		plugin = forcePlugin
	}
	if forceCmd != "" {
		log.Debugln("强制使用命令", forceCmd)
		cmd = forceCmd
	}

	//pm := tm.PluginManger
	log.Debugln("加载插件", plugin)
	err = pm.LoadPlugin(plugin)
	if err != nil {
		return &pb.CommonCmdReply{ResultMsg: "加载插件失败: " + err.Error()}, nil, nil
	}
	log.Debugln("插件", plugin, "加载成功")
	p := pm.PlugInfos[plugin].Plugin
	log.Debugln("预处理函数名")
	funcName := strings.ToUpper(cmd[0:1]) + strings.ToLower(cmd[1:])
	log.Debugln("查找函数", funcName)
	f, err1 := p.Lookup(funcName)
	if err1 != nil {
		f, err = p.Lookup("DefaultFunc")
		if err != nil {
			return &pb.CommonCmdReply{ResultMsg: "查找不到指定函数[" + funcName + "]: " + err1.Error()}, nil, nil
		}
		log.Debugf("未能找到指定函数[%s], 已执行默认函数", funcName)
	}

	switch f.(type) {
	case func(task Task) (string, error):
		replyMsg, err = f.(func(task Task) (string, error))(*t)
		if err != nil {
			return &pb.CommonCmdReply{ResultMsg: "返回结果错误: " + err.Error(), Status: pb.CommonCmdReply_Err}, nil, nil
		}
		reply = &pb.CommonCmdReply{ResultMsg: replyMsg, Status: pb.CommonCmdReply_Ok}
	case func(task Task) (*pb.CommonCmdReply, error):
		reply, err = f.(func(task Task) (*pb.CommonCmdReply, error))(*t)
		if err != nil {
			return &pb.CommonCmdReply{ResultMsg: "返回结果错误: " + err.Error(), Status: pb.CommonCmdReply_Err}, nil, nil
		}
	case func(task Task) (string, *Task, error):
		replyMsg, nextTask, err = f.(func(task Task) (string, *Task, error))(*t)
		if err != nil {
			return &pb.CommonCmdReply{ResultMsg: "返回结果错误: " + err.Error(), Status: pb.CommonCmdReply_Err}, nil, nil
		}
		reply = &pb.CommonCmdReply{ResultMsg: replyMsg, Status: pb.CommonCmdReply_Ok}
	case func(task Task) (*pb.CommonCmdReply, *Task, error):
		reply, nextTask, err = f.(func(task Task) (*pb.CommonCmdReply, *Task, error))(*t)
		if err != nil {
			return &pb.CommonCmdReply{ResultMsg: "返回结果错误: " + err.Error(), Status: pb.CommonCmdReply_Err}, nil, nil
		}
	default:
		reply = &pb.CommonCmdReply{ResultMsg: "插件的返回格式暂不支持!"}
	}
	return reply, nextTask, nil
}
