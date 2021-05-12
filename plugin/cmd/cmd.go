package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"errors"
	"fmt"
	"os/exec"
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
desc: "Cmd Plugin\n提供执行服务器各项命令的插件\n
cmd ls \n\t查看当前目录信息\n
cmd command [命令行指令] \n\t执行命令行指令\n
cmd script [脚本地址] \n\t执行shell脚本\n"
yess:
  cmd:
    type: "Plugin"
    desc: "执行服务器命令"
    yess:
      ls:
        type: "Cmd"
        desc: "查看当前目录信息"
      command:
        type: "Cmd"
        desc: "执行命令行指令"
      script:
        type: "Cmd"
        desc: "执行shell脚本"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

func Ls(task tasks.Task) (string, error) {
	log.Debugf("task信息: %#v", task)
	var flag []string
	if task.Flags != nil {
		for _, f := range task.Flags {
			flag = append(flag, "-"+f)
		}
	} else {
		flag = []string{"-lah"}
	}
	args := append(flag, task.SubCmd...)
	cmd := exec.Command("ls", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Command(task tasks.Task) (string, error) {
	out, err := exec.Command("bash", "-c", task.SubCmd[0]).Output()
	if err != nil {
		log.Errorf("执行命令失败: %s", task.SubCmd[0])
		return "", errors.New(fmt.Sprintf("执行命令失败: %s", task.SubCmd[0]))
	}
	return string(out), nil
}

func Script(task tasks.Task) (string, error) {
	cmd := exec.Command("sh", task.SubCmd[0])
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Errorf("执行脚本失败: %s", task.SubCmd[0])
		return "", errors.New(fmt.Sprintf("执行脚本 [%s] 失败: \n%s", task.SubCmd[0], out))
	}
	return string(out), nil
}

//func Run(task tasks.Task) (string, error) {
//	log.Debugf("task信息: %#v", task)
//	var (
//		flag []string
//		args []string
//	)
//
//	for _, f := range task.Flags {
//		flag = append(flag, "-"+f)
//	}
//	for k, v := range task.Args {
//		args = append(args, fmt.Sprintf("--%s=%s", k, v))
//	}
//
//	need := append(flag, args...)
//	need = append(task.SubCmd[1:], need...)
//	cmd := exec.Command(task.SubCmd[0], need...)
//	out, err := cmd.CombinedOutput()
//	if err != nil {
//		return "", err
//	}
//	return string(out), nil
//}

func main() {
	//cmd := exec.Command("ls", "plugin -lah")

	//cmd := exec.Command("sh", "scripts/test1.sh")
	//if err := cmd.Start(); err != nil {
	//	log.Fatalf("cmd.Start: %v", err)
	//}
	//
	//if err := cmd.Wait(); err != nil {
	//	if exiterr, ok := err.(*exec.ExitError); ok {
	//		// The program has exited with an exit code != 0
	//
	//		// This works on both Unix and Windows. Although package
	//		// syscall is generally platform dependent, WaitStatus is
	//		// defined for both Unix and Windows and in both cases has
	//		// an ExitStatus() method with the same signature.
	//		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
	//			log.Printf("Exit Status: %d", status.ExitStatus())
	//		}
	//	} else {
	//		log.Fatalf("cmd.Wait: %v", err)
	//	}
	//}

	cmd := exec.Command("sh", "scripts/test1.sh")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("err: ", err)
	}
	fmt.Printf("rst: %s\n", out)
	fmt.Println("code: ", cmd.ProcessState.ExitCode())

	//fmt.Println("aaa")
}
