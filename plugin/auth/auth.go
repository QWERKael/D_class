package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"errors"
	"fmt"
	ldap "github.com/go-ldap/ldap/v3"
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
desc: "Auth Plugin\n简单的认证插件, 支持LDAP形式的权限认证, 可以将该插件设置到前置过滤器中, 以实现登陆验证, 如: --pre-filter='auth,simple'"`
	return helpMsg, nil
}

func Simple(task tasks.Task) (string, *tasks.Task, error) {
	log.Debugln("开始Simple权限认证")
	var ok bool
	authString := task.Context["auth"]
	ldapURL := task.Context["ldapURL"]
	auth, err := utils.GetCred(authString)
	if err != nil {
		return "", nil, err
	}
	switch auth.Type {
	case "simple":
		ok = simpleAuth(auth)
	case "ldap":
		if ldapURL == "" {
			ldapURL = "ldap://ipaddress:port"
		}
		ok = ldapAuth(auth, ldapURL)
	default:
		return "", nil, errors.New("非法的认证方式")
	}
	if !ok {
		log.Debugln("认证失败")
		return "", nil, errors.New("认证失败")
	}
	return task.Cmd, &task, nil
}

func simpleAuth(cred *utils.Cred) bool {
	if cred.Type == "simple" && cred.User == "admin" && cred.Password == "admin" {
		return true
	}
	return false
}

func ldapAuth(cred *utils.Cred, ldapURL string) bool {
	fmt.Println("ldapURL: ", ldapURL)
	userName := cred.User
	fmt.Println("userName: ", userName)
	password := cred.Password
	fmt.Println("password: ", password)

	conn, err := ldap.DialURL(ldapURL)
	if err != nil {
		log.Errorln("连接LDAP地址失败: ", err)
		return false
	}
	defer conn.Close()
	err = conn.Bind(userName, password)
	if err != nil {
		log.Errorln("用户名密码错误: ", err)
		return false
	}
	return true
}
