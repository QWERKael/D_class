package main

import (
	"atk_D_class/logger"
	"atk_D_class/prompt"
	"atk_D_class/utils"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/metadata"
	"path/filepath"
)

var log = logger.Logger

var workDir *string = flag.String("wd", "", "工作目录")
var name *string = flag.String("name", "guess", "what's your name?")
var address *string = flag.String("address", "127.0.0.1:8881", "server address")
var pmtCnf *string = flag.String("prompt", "prompt.yml", "自定义提示信息")
var user *string = flag.String("user", "", "用户名")
var password *string = flag.String("password", "", "密码")
var version *bool = flag.Bool("version", false, "显示编译时间")

////自定义token认证
//type CustomerTokenAuth struct {
//}
////获取元数据
//func (c CustomerTokenAuth) GetRequestMetadata(ctx context.KV,
//	uri...string) (map[string]string, error) {
//
//	return map[string]string{
//		"user":  *user,
//		"password": *password,
//	}, nil
//}
////是否开启传输安全 TLS
//func (c CustomerTokenAuth) RequireTransportSecurity() bool {
//	return false
//}

func main() {
	flag.Parse()

	if *version {
		fmt.Println(utils.Version())
		return
	}

	logger.LogLevel("debug")

	// 设置工作目录
	dir, err := filepath.Abs(*workDir)
	utils.CheckErrorPanic(err, "工作目录不正确")
	pluginDir := filepath.Join(dir, "plugin")

	auth, err := utils.MakeCred("ldap", *user, *password)
	if err != nil {
		auth = ""
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = metadata.NewOutgoingContext(ctx,
		metadata.Pairs(
			"name", *name,
			"auth", auth,
		),
	)
	pmt := prompt.Prompt{
		Addr:      *address,
		Preload:   "",
		Ctx:       ctx,
		Cancel:    cancel,
		PluginDir: pluginDir,
	}

	go func() {
		if len(*user) > 0 && len(*password) > 0 {
			pmt.ConnectToAddr(false)
		} else {
			pmt.ConnectToAddr(true)
		}
		pmt.Prepare()
		prompt.LoadPrompt(*pmtCnf)
		pmt.Runner.Run()

	}()

	<-ctx.Done()
	fmt.Println("Bye")
}
