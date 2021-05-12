package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
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
desc: "Config Plugin\n为命令提供默认参数的插件\n`
	return helpMsg, nil
}

type PluginSetting struct {
	Commands map[string]*CommandSetting
	KV       map[string]string
}

type CommandSetting struct {
	DefaultValues *DefaultValues
}

type DefaultValues struct {
	SubCmd []string
	Flags  []string
	Args   map[string]string
}

func Attach(task tasks.Task) (string, *tasks.Task, error) {
	filePath := task.Context["configPath"]
	if !path.IsAbs(filePath) {
		filePath = path.Join(task.Context["workDir"], filePath)
	}
	log.Debugf("使用配置文件[%s]补全缺省参数", filePath)
	defaults, err := getConfig(filePath)
	log.Debugf("使用配置文件[%s]补全缺省参数", filePath)
	if err != nil {
		return "", nil, err
	}
	log.Debugf("已获取全部缺省值信息")
	log.Debugf("加载Context默认值...")
	if context, ok := defaults["Context"]; ok {
		for ctName, ctInfo := range context.KV {
			if _, ok := task.Context[ctName]; !ok {
				task.Context[ctName] = ctInfo
				fmt.Printf("加载Context参数 %s 的默认值为: %#v\n", ctName, ctInfo)
			}
			fmt.Printf("Context参数 %s 的值已存在", ctName)
		}
	}
	log.Debugf("加载Context默认值完毕！")
	var defaultValues *DefaultValues
	if pluginSetting, ok := defaults[task.Plugin]; ok {
		log.Debugf("已获取插件[%s]的缺省值信息: %#v", task.Plugin, task.Cmd, defaults[task.Plugin])
		if commandSetting, ok := pluginSetting.Commands[task.Cmd]; ok {
			log.Debugf("已获取命令[%s]的缺省值信息: %#v", task.Plugin, task.Cmd, defaults[task.Cmd])
			defaultValues = commandSetting.DefaultValues
			log.Debugf("已获取[%s -> %s]的缺省值信息: %#v", task.Plugin, task.Cmd, defaultValues)
		}
	}
	if defaultValues != nil {
		log.Debugln("查找未定义值, 并进行补全")
		if task.SubCmd == nil {
			log.Debugln("补全未定义的SubCmd")
			task.SubCmd = defaultValues.SubCmd
		}
		if task.Flags == nil {
			log.Debugln("补全未定义的Flags")
			task.Flags = defaultValues.Flags
		}

		log.Debugf("获取缺省值Args: [%#v]", defaultValues.Args)
		if task.Args == nil {
			task.Args = make(map[string]string)
		}
		for key, val := range defaultValues.Args {
			if _, ok := task.Args[key]; !ok {
				log.Debugf("补全未定义的Arg: [%s:%s]", key, val)
				task.Args[key] = val
			}
		}
	}
	log.Debugln("补全完成!")
	return "", &task, nil
}

func getConfig(filepath string) (map[string]PluginSetting, error) {
	if ft, err := utils.CheckPath(filepath); err == nil {
		if ft == utils.File {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				return nil, errors.New("打开文件失败: " + err.Error())
			}
			config := make(map[string]PluginSetting)
			err = yaml.Unmarshal(b, config)
			if err != nil {
				return nil, errors.New("解析配置文件失败: " + err.Error())
			}
			return config, nil
		} else {
			return nil, errors.New("指定路径不是文件: " + filepath)
		}
	} else {
		return nil, errors.New("路径检测失败: " + err.Error())
	}
}

func main() {
	in := PluginSetting{
		Commands: map[string]*CommandSetting{
			"cmd1": {
				DefaultValues: &DefaultValues{
					SubCmd: []string{"s1", "s2"},
					Flags:  []string{"f1", "f2"},
					Args:   map[string]string{"k1": "v1", "k2": "v2"},
				},
			},
		},
	}
	ins := map[string]PluginSetting{
		"p1": in,
		"p2": in,
	}
	out, err := yaml.Marshal(ins)
	if err != nil {
		fmt.Println("错误: ", err.Error())
	} else {
		fmt.Println("配置: ")
		fmt.Printf("%s", out)
	}

	config, err := getConfig("target/config/d.yml")
	if err != nil {
		fmt.Println("错误: ", err.Error())
	}
	if config == nil {
		fmt.Println("错误: 配置为空")
	} else {
		for pluginName, pluginInfo := range config {
			fmt.Printf("读取插件 %s 的配置文件: %#v\n", pluginName, pluginInfo)
			for ctName, ctInfo := range pluginInfo.KV {
				fmt.Printf("读取Context %s 的配置文件: %#v\n", ctName, ctInfo)
				fmt.Printf("%#v\n", ctInfo)
			}
			for cmdName, cmdInfo := range pluginInfo.Commands {
				fmt.Printf("读取命令 %s 的配置文件: %#v\n", cmdName, cmdInfo)
				fmt.Printf("%#v\n", cmdInfo.DefaultValues)
			}
		}
	}
}
