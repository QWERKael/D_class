package prompt

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"context"
	"fmt"
	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var log = logger.Logger

//var suggestion = []prompt.Suggest{
//	{Text: "show ", Description: "展示命令"},
//	{Text: "net", Description: "查看网卡信息"},
//	{Text: "load", Description: "查看负载"},
//	{Text: "processlist", Description: "查看进程"},
//	{Text: "upload ", Description: "上传文件"},
//	{Text: "change-to ", Description: "改变连接到新的服务端"},
//	{Text: "restart", Description: "重启服务端"},
//	{Text: "exit", Description: "退出"},
//}

type Prompt struct {
	Addr      string
	Preload   string
	CC        pb.CommanderClient
	Ctx       context.Context
	Cancel    context.CancelFunc
	Runner    *prompt.Prompt
	PluginDir string
}

func (pmt *Prompt) ConnectToAddr(reAuth bool) {
	if reAuth {
		auth, err := utils.Simple()
		utils.CheckErrorPanic(err)
		if md, ok := metadata.FromOutgoingContext(pmt.Ctx); ok {
			md.Set("auth", auth)
			pmt.Ctx = metadata.NewOutgoingContext(pmt.Ctx, md)
		}
	}
	conn, err := grpc.Dial(pmt.Addr, grpc.WithInsecure())
	utils.CheckErrorPanic(errors.WithMessage(err, "连接到指定地址失败: "+pmt.Addr))
	client := pb.NewCommanderClient(conn)
	log.Infoln("连接到", pmt.Addr, "...")
	pmt.CC = client
	err = es.ESs["change-to"].Add(pmt.Addr, "IP地址", Cmd)
	if err != nil {
		log.Debugf("IP地址添加到快捷命令失败: %s", err.Error())
	}
	//suggestion = append(suggestion, prompt.Suggest{Text: pmt.Addr, Description: "IP地址"})
}

func (pmt *Prompt) executor(line string) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		if err == nil {
			return
		}
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()

	if strings.TrimSpace(line) == "" {
		return
	}

	if !strings.HasPrefix(line, "$set") && !strings.HasPrefix(line, "$unset") {
		line = pmt.Preload + " " + line
	}
	ccr, isLocal := Parse(line)

	//fmt.Printf("Preload: %s\n", pmt.Preload)
	//fmt.Printf("ccr: %#v\n", ccr)

	if strings.ToLower(ccr.Plugin) == "$set" {
		switch strings.ToLower(ccr.Cmd) {
		// 设置预加载命令, 如果需要重复使用相同的命令前缀, 可以设置此值, 避免重复输入
		case "preload":
			pmt.Preload = ccr.SubCmd[0]
		}
		return
	}

	if strings.ToLower(ccr.Plugin) == "$unset" {
		switch strings.ToLower(ccr.Cmd) {
		case "preload":
			pmt.Preload = ""
		}
		return
	}

	// 收到isLocal的标志时, 在本地运行插件
	if isLocal {
		var reverseFlag bool
		ccr.Flags, reverseFlag = utils.IfInFlagThenPop("reverse", ccr.Flags)
		pm := &tasks.PluginManager{
			PluginDir: pmt.PluginDir,
			PlugInfos: make(map[string]tasks.PlugInfo),
		}
		task := tasks.Task{Plugin: ccr.Plugin, Cmd: ccr.Cmd, SubCmd: ccr.SubCmd,
			Flags: ccr.Flags, Args: ccr.Args}

		//r, err := tasks.Exec(&task, pm)
		//utils.CheckErrorPanic(err)
		var (
			reply    *pb.CommonCmdReply
			nextTask *tasks.Task
			err      error
		)
		for {
			reply, nextTask, err = tasks.Exec(&task, pm, "", "")
			if nextTask == nil {
				break
			}
			task = *nextTask
			utils.CheckErrorPanic(err)
		}

		err = CommonCmdOutputter(reply, reverseFlag)
		utils.CheckErrorPanic(err)
		return
	}

	// 处理一些特殊命令
	switch ccr.Plugin {
	case "":
		return
	case "exit":
		pmt.Cancel()
		return
	case "upload":
		ccr.Type = pb.CommonCmdRequest_FILE_TRANSFER
	case "download":
		ccr.Type = pb.CommonCmdRequest_FILE_TRANSFER
	case "change-to":
		pmt.Addr = ccr.Cmd
		pmt.ConnectToAddr(true)
		return
	}

	switch {
	case ccr.Type == pb.CommonCmdRequest_COMMON_CMD:
		var (
			reverseFlag bool
			//asyncFlag   bool
		)
		ccr.Flags, reverseFlag = utils.IfInFlagThenPop("reverse", ccr.Flags)
		//ccr.Flags, asyncFlag = utils.IfInFlagThenPop("async", ccr.Flags)
		//if asyncFlag {
		//	ccr.Type = pb.CommonCmdRequest_ASYNC_TASK
		//}
		r, err := pmt.CC.CommonCmd(pmt.Ctx, &ccr)
		utils.CheckErrorPanic(err)
		// 对help命令做特殊处理, 加载其中的智能提示, 并返回帮助信息
		if ccr.Cmd == "help" {
			resEs := ParseYAML([]byte(r.ResultMsg))
			r.ResultMsg = resEs.Suggest.Description
			for k, v := range resEs.ESs {
				es.ESs[k] = v
			}
		}

		err = CommonCmdOutputter(r, reverseFlag)
		utils.CheckErrorPanic(err)
	case ccr.Type == pb.CommonCmdRequest_FILE_TRANSFER && ccr.Plugin == "upload":
		localPath := ccr.Cmd
		filePath, fileName := filepath.Split(ccr.Cmd)
		filePath = uploadPath(ccr.Flags)
		applyFi, err := pmt.CC.ApplyTransfer(pmt.Ctx,
			&pb.TransferInfo{
				Type:       pb.TransferInfo_Upload,
				State:      pb.TransferInfo_Apply,
				FileName:   fileName,
				FilePath:   filePath,
				TransferId: 0,
			})
		utils.CheckErrorPanic(err)
		uploadStream, err := pmt.CC.Upload(pmt.Ctx)
		utils.CheckErrorPanic(err)
		// 创建一个1M的buf
		defer uploadStream.CloseSend()
		file, err := os.Open(localPath)
		utils.CheckErrorPanic(err)
		//stat, err := file.Stat()
		//utils.CheckErrorPanic(err)
		//fmt.Printf("文件大小: %d\n", stat.Size())
		buf := make([]byte, 1<<20)
		writing := true
		for writing {
			//fmt.Printf("读取文件 %s ...\n", fileName)
			n, err := file.Read(buf[:])
			if err != nil {
				if err == io.EOF {
					fmt.Println("文件已发送")
					writing = false
					err = nil
					continue
				}
				utils.CheckErrorPanic(err)
			}
			err = uploadStream.Send(&pb.Chunks{
				TransferId: applyFi.TransferId,
				Size:       int64(n),
				Content:    buf[:n],
			})
		}
		recvFi, err := uploadStream.CloseAndRecv()
		utils.CheckErrorPanic(err)
		checkUploadResult(recvFi, localPath, fileName)
	}
	return
}

func (pmt *Prompt) livePrefix() (string, bool) {
	return fmt.Sprintf("[%s] %s >>> ", pmt.Addr, pmt.Preload), true
}

func (pmt *Prompt) Prepare() {
	pmt.Runner = prompt.New(
		pmt.executor,
		completer,
		prompt.OptionPrefix(fmt.Sprintf("[%s] %s >>> ", pmt.Addr, pmt.Preload)),
		prompt.OptionLivePrefix(pmt.livePrefix),
		prompt.OptionTitle("D级人员委派指南"),
		prompt.OptionPrefixTextColor(prompt.Yellow),
		prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),
	)
}

// 将命令行解析为CommonCmdRequest格式
func Parse(line string) (pb.CommonCmdRequest, bool) {
	isLocal := false
	ccr := pb.CommonCmdRequest{
		Type:   0,
		Plugin: "",
		Cmd:    "",
		SubCmd: nil,
		Flags:  nil,
		Args:   make(map[string]string),
	}
	fields, err := shlex.Split(line)
	utils.CheckErrorPanic(err)

	if len(fields) < 1 {
		return ccr, false
	}

	// 支持client端加载插件
	if strings.ToLower(fields[0]) == "loc" {
		isLocal = true
		fields = fields[1:]
	}

	// 解析plugin
	ccr.Plugin = fields[0]
	if len(fields) > 1 {
		// 解析cmd
		ccr.Cmd = fields[1]

		fields = fields[2:]
		l := len(fields)
		// 将arg参数等号前后有空格的情况合并为一个field
		for i := 0; i < l; i++ {
			if fields[i] == "" {
				continue
			}
			if fields[i] == "=" {
				// "--key = value"的情况
				fields[i-1] = fields[i-1] + fields[i] + fields[i+1]
				fields[i] = ""
				fields[i+1] = ""
				continue
			} else if strings.HasPrefix(fields[i], "=") {
				// "--key =value"的情况
				fields[i-1] = fields[i-1] + fields[i]
				fields[i] = ""
				continue
			} else if fields[i][len(fields[i])-1:] == "=" {
				// "--key= value"的情况
				fields[i] = fields[i] + fields[i+1]
				fields[i+1] = ""
				continue
			}
		}
	}

	//解析subcmd, flag, args
	for _, field := range fields {
		if field == "" {
			continue
		}
		if strings.HasPrefix(field, "-") {
			if strings.HasPrefix(field, "--") {
				//fmt.Println("这是一个arg")
				kv := strings.SplitN(field[2:], "=", 2)
				//fmt.Printf("%#v\n", kv)
				ccr.Args[kv[0]] = kv[1]
			} else {
				ccr.Flags = append(ccr.Flags, field[1:])
				//fmt.Println("这是一个flag:\t", )
			}
		} else {
			ccr.SubCmd = append(ccr.SubCmd, field)
		}
	}
	//for i, suggestion := range fields {
	//	if i == 0 {
	//		ccr.Plugin = suggestion
	//	} else {
	//
	//	}
	//}
	return ccr, isLocal
}
