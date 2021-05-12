package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	hc "atk_D_class/pb/honeycomb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
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
desc: "Bee Plugin\nBee插件, Honeycomb的客户端，一个简单的K-V存储模块\n
bee connect --addr='Ip:port' \n\t
连接指定的Queen\n
bee use [dbName] \n\t
跳转到指定的库\n
bee set [key] [value] \n\t
插入或修改指定的值\n
bee get [key] \n\t
获取指定的值\n
bee del [key] \n\t
删除指定的值\n
bee keys \n\t
列出当前库的全部key\n
bee showdb \n\t
列出所有的库\n
bee createdb [dbName] \n\t
创建指定的库\n
bee status \n\t
查看连接状态\n"
yess:
  bee:
    type: "Plugin"
    desc: "连接到指定Honeycomb进行相关操作"
    yess:
      connect:
        type: "Cmd"
        desc: "连接指定的Queen"
        yess:
          "--addr":
            type: "ArgKey"
            desc: "指定Queen的地址, 默认值为 [::1]:10000"
      use:
        type: "Cmd"
        desc: "跳转到指定的库"
      set:
        type: "Cmd"
        desc: "插入或修改指定的值"
      get:
        type: "Cmd"
        desc: "获取指定的值"
      del:
        type: "Cmd"
        desc: "删除指定的值"
      keys:
        type: "Cmd"
        desc: "列出当前库的全部key"
      showdb:
        type: "Cmd"
        desc: "列出所有的库"
      createdb:
        type: "Cmd"
        desc: "创建指定的库"
      status:
        type: "Cmd"
        desc: "查看连接状态"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

type Bee struct {
	Addr    string
	UsingDB string
	CC      hc.CommanderClient
	Ctx     context.Context
	Cancel  context.CancelFunc
}

var bee *Bee

func defaultBee() *Bee {
	return &Bee{
		Addr:    "[::1]:10000",
		UsingDB: "db0",
		CC:      nil,
		Ctx:     nil,
		Cancel:  nil,
	}
}

func Status(task tasks.Task) (string, error) {
	return fmt.Sprintf("bee status:\n %#v\n", bee), nil
}

func (bee *Bee) connect() error {
	// Set up a connection to the server.
	conn, err := grpc.Dial(bee.Addr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	//defer conn.Close()
	c := hc.NewCommanderClient(conn)

	// Contact the server and print out its response.
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()

	bee.CC = c
	bee.Ctx = ctx
	bee.Cancel = cancel
	return nil
}

func tryConnect() error {
	var err error
	if bee == nil {
		bee = defaultBee()
	}
	if bee.CC == nil {
		err = bee.connect()
	}
	if err != nil {
		return err
	}
	return nil
}

func (bee *Bee) disconnect() error {
	if bee != nil && bee.CC != nil {
		bee.Cancel()
	}
	return nil
}

func Connect(task tasks.Task) (string, error) {
	addr := utils.StringDefault(task.Args["addr"], "[::1]:10000")
	if bee == nil {
		bee = defaultBee()
		bee.Addr = addr
	} else {
		_ = bee.disconnect()
	}
	_ = bee.connect()
	return "bee connected", nil
}

func commandCall(req *hc.CmdRequest) (*hc.CmdReply, error) {
	err := tryConnect()
	if err != nil {
		return nil, err
	}
	rep, err := bee.CC.CmdCall(bee.Ctx, req)
	if err != nil {
		return nil, err
	}
	if rep.GetStatus() == hc.CmdReply_Err {
		return nil, errors.New(rep.Message)
	}
	return rep, nil
}

func Use(task tasks.Task) (string, error) {
	dbName := task.SubCmd[0]
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_ShowDB,
		Args:    [][]byte{},
	})
	if err != nil {
		return "", err
	}
	for _, result := range rep.Results {
		resName := string(result)
		if resName == dbName {
			bee.UsingDB = resName
			return "已切换到 " + resName, nil
		}
	}
	return "未找到 " + dbName, nil
}

func Showdb(task tasks.Task) (*pb.CommonCmdReply, error) {
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_ShowDB,
		Args:    [][]byte{},
	})
	if err != nil {
		return nil, err
	}

	header := []string{"DB List"}
	body := [][]string{nil}
	for _, result := range rep.Results {
		body = append(body, []string{string(result)})
	}
	reply := prompt.ToTable(header, []string{fmt.Sprintf("总计: %d 行", len(rep.Results))}, body, 0)
	return reply, nil
}

func Createdb(task tasks.Task) (string, error) {
	dbName := task.SubCmd[0]
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_CreateDB,
		Args:    [][]byte{[]byte(dbName)},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", string(rep.Results[0])), nil
}

func Keys(task tasks.Task) (*pb.CommonCmdReply, error) {
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_Keys,
		Args:    [][]byte{},
	})
	if err != nil {
		return nil, err
	}

	header := []string{"Keys List"}
	body := [][]string{nil}
	for _, result := range rep.Results {
		body = append(body, []string{string(result)})
	}
	reply := prompt.ToTable(header, []string{fmt.Sprintf("总计: %d 行", len(rep.Results))}, body, 0)
	return reply, nil
}

func Set(task tasks.Task) (string, error) {
	key := task.SubCmd[0]
	val := task.SubCmd[1]
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_Set,
		Args:    [][]byte{[]byte(key), []byte(val)},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", rep.Message), nil
}

func Get(task tasks.Task) (string, error) {
	key := task.SubCmd[0]
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_Get,
		Args:    [][]byte{[]byte(key)},
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", string(rep.Results[0])), nil
}

func Del(task tasks.Task) (string, error) {
	keys := [][]byte{nil}
	for _, key := range task.SubCmd {
		keys = append(keys, []byte(key))
	}
	rep, err := commandCall(&hc.CmdRequest{
		UsingDB: bee.UsingDB,
		Cmd:     hc.CmdRequest_Delete,
		Args:    keys,
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", rep.Message), nil
}

func main() {
	bee = defaultBee()
	defer bee.disconnect()
	_ = bee.connect()
	status, _ := Status(tasks.Task{})
	println(status)

	rst, err := Showdb(tasks.Task{})
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	fmt.Printf("show db:\n%s", rst)

	//time.Sleep(4*time.Second)

	//rst, err = Set(tasks.Task{
	//	SubCmd: []string{"key1", "value1"},
	//})
	//if err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	//fmt.Println("set:", rst)
	//
	//rst, err = Get(tasks.Task{
	//	SubCmd: []string{"key1"},
	//})
	//if err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	//fmt.Println("get:", rst)
	//
	//rst, err = Set(tasks.Task{
	//	SubCmd: []string{"key2", "value2"},
	//})
	//if err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	//fmt.Println("set:", rst)
	//
	//rst, err = Keys(tasks.Task{})
	//if err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	//fmt.Println("keys:", rst)
	//
	//rst, err = Del(tasks.Task{
	//	SubCmd: []string{"key1", "key2"},
	//})
	//if err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	//fmt.Println("delete:", rst)
	//
	//rst, err = Keys(tasks.Task{})
	//if err != nil {
	//	log.Fatalf("failed to serve: %v", err)
	//}
	//fmt.Println("keys:", rst)
}
