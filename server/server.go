package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var log = logger.Logger

var workDir *string = flag.String("wd", "", "工作目录")
var listenAddr *string = flag.String("listen", "127.0.0.1:8881", "监听地址")
var logLevel *string = flag.String("lvl", "debug", "日志级别")
var capacity *int = flag.Int("cap", 100, "最大异步任务队列数")

//var unauth *bool = flag.Bool("unauth", false, "不进行权限验证")
//var auth *string = flag.String("auth", "Simple", "认证方式")
var preFilter *string = flag.String("pre-filter", "config,attach;auth,simple", "前置过滤器")
var configPath *string = flag.String("config", "config/d.yml", "配置文件地址（需要config前置过滤器）")
var nonFilter *bool = flag.Bool("non-filter", false, "不进行前置过滤")
var version *bool = flag.Bool("version", false, "显示编译时间")

type Service struct {
	bootFlag chan int
	WorkDir  string
	//TaskManager *recycle.TaskManager
	PluginManager *tasks.PluginManager
	FileManager   *tasks.FileManager
}

func (s *Service) CommonCmd(ctx context.Context, CommonCmdRequest *pb.CommonCmdRequest) (*pb.CommonCmdReply, error) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()

	md, _ := metadata.FromIncomingContext(ctx)

	if CommonCmdRequest.Plugin == "restart" {
		log.Debugln("已收到启动命令")
		s.bootFlag <- 1
		log.Debugln("已发送重启标志")
		return &pb.CommonCmdReply{ResultMsg: "重启服务..."}, nil
	}

	c := make(map[string]string)
	c["listenAddr"] = *listenAddr
	c["workDir"] = s.WorkDir
	c["configPath"] = *configPath
	c["auth"] = md.Get("auth")[0]

	task := tasks.Task{Plugin: CommonCmdRequest.Plugin, Cmd: CommonCmdRequest.Cmd, SubCmd: CommonCmdRequest.SubCmd,
		Flags: CommonCmdRequest.Flags, Args: CommonCmdRequest.Args, Context: c}

	//if CommonCmdRequest.Type == pb.CommonCmdRequest_ASYNC_TASK {
	//	id := s.TaskManager.AddAsyncTask(&task)
	//	return &pb.CommonCmdReply{ResultMsg: fmt.Sprintf("异步任务已添加, 任务编号 [%d]", id)}, nil
	//}
	//
	//if task.Plugin == "async" {
	//	switch task.Cmd {
	//	case "state":
	//		reply := s.TaskManager.StateList.ToTable()
	//		return reply, nil
	//	case "get":
	//		id := utils.StringDefaultInt(task.SubCmd[0], -1)
	//		reply := s.TaskManager.AsyncTaskList[id].Result
	//		return reply, nil
	//	}
	//}
	log.Debugf("执行task: %#v", task)
	var (
		reply          *pb.CommonCmdReply
		nextTask       *tasks.Task
		err            error
		preFilterItems [][]string
	)

	//前置过滤器
	if !*nonFilter {
		log.Debugln("进行前置过滤")

		for _, item := range strings.Split(*preFilter, ";") {
			preFilterItems = append(preFilterItems, strings.Split(item, ","))
		}

		for _, item := range preFilterItems {
			if reply, nextTask, err = tasks.Exec(&task, s.PluginManager, item[0], item[1]); reply == nil || reply.Status == pb.CommonCmdReply_Err {
				log.Debugln("前置过滤器执行失败: %s, %s", item[0], item[1])
				return reply, nil
			}
			if nextTask == nil {
				break
			}
			task = *nextTask
			utils.CheckErrorPanic(err)
		}
	}

	for {
		reply, nextTask, err = tasks.Exec(&task, s.PluginManager, "", "")
		if nextTask == nil {
			break
		}
		task = *nextTask
		utils.CheckErrorPanic(err)
	}
	return reply, nil
	//return &pb.CommonCmdReply{ResultMsg: rst}, nil
}

func (s *Service) ApplyTransfer(ctx context.Context, TransferInfo *pb.TransferInfo) (*pb.TransferInfo, error) {
	fi := tasks.Unpack(TransferInfo)
	//absFilePath, err := filepath.Abs(fi.FilePath)
	//if err != nil {
	//	fi.TransferState = pb.TransferInfo_Error
	//	fi.ErrorMsg = ""
	//	ti := fi.Assemble()
	//	return &ti, nil
	//}
	//filepath.Rel(s.WorkDir, fi.FilePath)
	fi = s.FileManager.Add(fi)
	ti := fi.Assemble()
	return &ti, nil
}

func (s *Service) Upload(uploadStream pb.Commander_UploadServer) error {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("runtime error:", err)
		default: // 非运行时错误
			fmt.Println("error:", err)
		}
	}()

	chunk, err := uploadStream.Recv()
	if err != nil {
		if err == io.EOF {
			ti := pb.TransferInfo{State: pb.TransferInfo_Error, ErrorMsg: "没有接收到任何数据"}
			err = uploadStream.SendAndClose(&ti)
			utils.CheckErrorPanic(err, "没有接收到任何数据\n")
			return nil
		}
		utils.CheckErrorPanic(err)
	}
	tid := chunk.TransferId
	fi := s.FileManager.FileInfos[tid]
	absPath := filepath.Join(s.WorkDir, fi.FilePath, fi.FileName)
	err = utils.CheckAndCreateDir(filepath.Dir(absPath))
	utils.CheckErrorPanic(err)
	log.Infoln("写入文件: ", absPath)
	file, err := os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	err = s.FileManager.BuildFile(file, chunk)
	utils.CheckErrorPanic(err)
	for {
		chunk, err := uploadStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			utils.CheckErrorPanic(err)
		}
		err = s.FileManager.BuildFile(file, chunk)
		utils.CheckErrorPanic(err)
	}
	err = file.Close()
	utils.CheckErrorPanic(err)
	fi.Md5 = utils.SumMd5FromFile(absPath)
	fi.TransferState = pb.TransferInfo_Complete
	ti := fi.Assemble()
	err = uploadStream.SendAndClose(&ti)
	utils.CheckErrorPanic(err)
	return nil
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(utils.Version())
		return
	}
	logger.LogLevel(*logLevel)
	dir, err := filepath.Abs(*workDir)
	utils.CheckErrorPanic(err, "工作目录不正确")
	pluginDir := filepath.Join(dir, "plugin")
	//connect_pool = &ConnectPool{}
	log.Infoln("开始...")
	log.Infoln(utils.Version())
	log.Debugln("运行PID: ", os.Getpid())

	//设置退出时执行的cmd命令, 用于重启server端
	cmd := exec.Command(`echo "服务已退出"`)
	server := grpc.NewServer()

	s := &Service{
		WorkDir:  dir,
		bootFlag: make(chan int),
		//TaskManager: recycle.MakeTaskManager(*capacity, pluginDir),
		PluginManager: &tasks.PluginManager{
			PluginDir: pluginDir,
			PlugInfos: make(map[string]tasks.PlugInfo),
		},
		FileManager: &tasks.FileManager{
			FileInfos: []tasks.FileInfo{},
		},
	}


	// 启动任务管理器
	//go func() {
	//	s.TaskManager.AsyncRunner()
	//}()

	// 监听重启
	go func() {
		bootFlag := <-s.bootFlag
		log.Debugln("获取启动标志: ", bootFlag)
		switch bootFlag {
		case 1:
			server.GracefulStop()
			// 设置重启命令
			cmd = exec.Command(os.Args[0], os.Args[1:]...)
		}
		bootFlag = 0
	}()

	lis, err := net.Listen("tcp", *listenAddr)
	utils.CheckErrorPanic(err)
	log.Infoln("Listen on", *listenAddr)
	log.Debugln("注册 commander server")
	pb.RegisterCommanderServer(server, s)

	reflection.Register(server)

	err = server.Serve(lis)
	if err == grpc.ErrServerStopped {
		log.Debugln("服务已关闭")
	} else {
		utils.CheckErrorPanic(err)
	}

	// 执行重启(或其他)命令
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	utils.CheckErrorPanic(err)
}
