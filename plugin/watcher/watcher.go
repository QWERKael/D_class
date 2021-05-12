package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/plugin/watcher/common"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"net"
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
desc: "Watcher Plugin\nWatcher插件, 用于接受Sentry返回的心跳数据\n
watcher run --listen='Ip:port' --timeout=N \n\t启动watcher, 指定监听的地址和超时时间\n"
yess:
  watcher:
    type: "Plugin"
    desc: "提供Water相关操作"
    yess:
      run:
        type: "Cmd"
        desc: "启动watcher"
        yess:
          "--listen":
            type: "ArgKey"
            desc: "用于监听指定的端口, 默认值为 [0.0.0.0:8880]"
          "--timeout":
            type: "ArgKey"
            desc: "接受心跳的超时时间, 默认值为 [3] 秒"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

type SentryInfo struct {
	Host string
	//Port  int
	Desc  string
	State *map[pb.Items]common.HealthState
}

type Service struct {
	timeout  int64
	Sentries map[string]*SentryInfo
}

var s *Service

func Info(task tasks.Task) (string, error) {
	log.Debugf("task信息: %#v", task)
	l := utils.Lines{}
	l.LineAppend(" %16s | %10s | %10s | %s", "HOST", "ITEM", "STATE", "DESC")
	for _, si := range s.Sentries {
		host := si.Host
		for k, v := range *(si.State) {
			l.LineAppend(" %16s | %10s | %10s | %20s",
				host, k.String(), common.HealthState_name[int32(v)], si.Desc)
		}
	}
	return l.String(), nil
}

func (s *Service) Register(ctx context.Context, req *pb.RegRequest) (*pb.RegReply, error) {
	var unicode string
	for {
		r := common.RandStringBytes(12)
		if _, ok := s.Sentries[r]; !ok {
			unicode = r
			break
		}
	}
	s.Sentries[req.LocalAddr] = &SentryInfo{
		Host:  req.LocalAddr,
		Desc:  req.Desc,
		State: &map[pb.Items]common.HealthState{req.Item: common.Init},
	}
	return &pb.RegReply{State: pb.RegReply_Agree, UniCode: unicode}, nil
}

func (s *Service) HeartBeat(stream pb.Watch_HeartBeatServer) error {
	d := time.Duration(s.timeout * int64(time.Second))
	t := time.NewTimer(d)

	var host string
	var item pb.Items
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		host = md.Get("host")[0]
		item = pb.Items(pb.Items_value[md.Get("item")[0]])
		log.Debugf("获取到 host 为 %s, item 为 %s", host, item.String())
	} else {
		return errors.New("未找到 host, item ")
	}

	state := s.Sentries[host].State
	go func() {
		for {
			ping, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					(*state)[item] = common.Die
					log.Debugln("连接已结束")
					return
				} else {
					(*state)[item] = common.Die
					return
				}
			}
			log.Debugf("收到ping信息: %#v", ping)
			if ping.PingState == pb.Ping_Active {
				(*state)[item] = common.Alive
				t.Reset(d)
			}
			pong := pb.Pong{Item: ping.Item, PongState: pb.Pong_Received}
			log.Debugf("发送pong信息: %#v", pong)
			err = stream.Send(&pong)
			utils.CheckErrorDebugLog(err, "发送pong消息出现错误: %s")
		}
	}()

	<-t.C
	(*state)[item] = common.Die
	log.Debugln("心跳连接已超时")
	return errors.New("心跳连接已超时")
}

func Run(task tasks.Task) (string, error) {
	s = &Service{Sentries: make(map[string]*SentryInfo)}
	listenAddr := utils.StringDefault(task.Args["listen"], "0.0.0.0:8880")
	s.timeout = utils.StringDefaultInt64(task.Args["timeout"], 3)
	//logger.LogLevel("debug")
	log.Debugf("task信息: %#v", task)
	log.Infoln("开始...")

	server := grpc.NewServer()
	lis, err := net.Listen("tcp", listenAddr)
	utils.CheckErrorPanic(err)
	log.Infoln("Listen on", listenAddr)
	log.Debugln("注册 commander server")
	pb.RegisterWatchServer(server, s)
	go func() {
		err := server.Serve(lis)
		utils.CheckErrorPanic(err)
	}()
	return "watcher监听中...", nil
}

func main() {
	logger.LogLevel("debug")
	s, err := Run(tasks.Task{})
	utils.CheckErrorPanic(err)
	fmt.Println(s)
	select {}
}
