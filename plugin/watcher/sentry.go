package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"strings"
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
desc: "Sentry Plugin\nSentry插件, 用于向Watcher发送的心跳数据\n
sentry run --connect-to='Ip:port' --sec=Inv --timeout=N --desc=[string] \n\t
启动sentry, 指定watcher的地址, 发送心跳的间隔时间, 超时时间 和对sentry的描述\n"
yess:
  sentry:
    type: "Plugin"
    desc: "提供Sentry相关操作"
    yess:
      run:
        type: "Cmd"
        desc: "启动sentry"
        yess:
          "--connect-to":
            type: "ArgKey"
            desc: "指定watcher地址, 默认值为 [127.0.0.1:8880]"
          "--sec":
            type: "ArgKey"
            desc: "发送心跳的间隔时间, 默认值为 [1] 秒"
          "--timeout":
            type: "ArgKey"
            desc: "发送心跳的超时时间, 默认值为 [3] 秒"
          "--desc":
            type: "ArgKey"
            desc: "对sentry的描述, 方便watcher对该sentry做备注"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

type Sentry struct {
	Addr      string
	LocalAddr string
	CC        pb.WatchClient
	Ctx       context.Context
	Cancel    context.CancelFunc
}

var sentry *Sentry

func (sentry *Sentry) ConnectToAddr() {
	conn, err := grpc.Dial(sentry.Addr, grpc.WithInsecure())
	utils.CheckErrorPanic(errors.WithMessage(err, "连接到指定地址失败: "+sentry.Addr))
	client := pb.NewWatchClient(conn)
	log.Infoln("连接到", sentry.Addr, "...")
	sentry.CC = client
}

func (sentry *Sentry) RegSentry(servAddr string, desc string) error {
	item := pb.Items_MainServ
	reply, err := sentry.CC.Register(sentry.Ctx,
		&pb.RegRequest{Item: item, LocalAddr: servAddr})
	utils.CheckErrorPanic(err)
	if reply.State == pb.RegReply_Agree {
		sentry.Ctx = metadata.NewOutgoingContext(sentry.Ctx,
			metadata.Pairs("host", servAddr, "item", item.String()))
		return nil
	} else {
		return errors.New("注册 Sentry 失败!")
	}
}

func Run(task tasks.Task) (string, error) {
	addr := utils.StringDefault(task.Args["connect-to"], "127.0.0.1:8880")
	sec := utils.StringDefaultInt64(task.Args["sec"], 1)
	timeout := utils.StringDefaultInt64(task.Args["timeout"], 3)
	desc := utils.StringDefault(task.Args["desc"], "")

	log.Debugf("向 %s 报告心跳", addr)
	ctx, cancel := context.WithCancel(context.Background())
	//ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("name", *name))
	sentry = &Sentry{Addr: addr, Ctx: ctx, Cancel: cancel}
	sentry.ConnectToAddr()

	// 注册unicode
	port := strings.Split(task.Context["listenAddr"], ":")[1]
	servAddr := utils.GetLocalIp(sentry.Addr) + ":" + port
	err := sentry.RegSentry(servAddr, desc)
	utils.CheckErrorPanic(err)

	stream, err := sentry.CC.HeartBeat(sentry.Ctx)
	utils.CheckErrorPanic(err)
	ds := time.Duration(sec * int64(time.Second))
	go func() {
		for {
			ping := pb.Ping{Item: pb.Items_MainServ, PingState: pb.Ping_Active}
			log.Debugf("发送ping消息: %#v", ping)
			err := stream.Send(&ping)
			if err != nil {
				if err == io.EOF {
					log.Debugln("连接已结束")
					break
				} else {
					log.Debugf("发送ping信息出现错误: %#v", err.Error())
					break
				}
			}
			time.Sleep(ds)
		}
		return
	}()
	go func() {
		d := time.Duration(timeout * int64(time.Second))
		t := time.NewTimer(d)
		go func() {
			for {
				pong, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						log.Debugln("连接已结束")
						break
					} else {
						log.Debugf("接受pong消息出现错误: %#v", err.Error())
						break
					}
				} else {
					t.Reset(d)
				}
				log.Debugf("收到pong信息: %#v", pong)
			}
			return
		}()
		<-t.C
		log.Debugln("心跳连接已超时")
		return
	}()
	return "执行成功", nil
}

func main() {
	s, err := Run(tasks.Task{})
	utils.CheckErrorPanic(err)
	fmt.Println(s)
}
