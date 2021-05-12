package main

import (
	"atk_D_class/tasks"
	"fmt"
)

type T struct {
	Text string
}

var t1 = &T{Text: "开始"}

func New() tasks.PlugInterface {
	return &T{Text: "初始化"}
}

func Set(s string) string {
	t := t1
	t.Text = s
	fmt.Printf("结构体t的结构为: %#v", t)
	return t.Text
}

func Get() string {
	t := t1
	fmt.Printf("结构体t的结构为: %#v", t)
	return t.Text
}

func (t *T) Help() string {
	fmt.Println("使用方法")
	return "使用方法"
}

func main() {
	fmt.Println("测试")
}
