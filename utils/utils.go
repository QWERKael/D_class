package utils

import (
	"atk_D_class/logger"
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

var log = logger.Logger

type Lines struct {
	StrBuf bytes.Buffer
}

func CheckErrorPanic(err error, errMsg ...string) {
	if err != nil {
		err = errors.Wrap(err, "我们似乎遇到了一些状况...")
		for _, e := range errMsg {
			fmt.Println(e)
		}
		logger.Logger.Debugf("%+v\n", err)
		panic(err.Error())
	}
}

func CheckErrorDebugLog(err error, errMsg string) {
	if err != nil {
		log.Debugf(errMsg, err.Error())
	}
}

func (l *Lines) LineAppend(format string, args ...interface{}) {
	if len(args) == 0 {
		l.StrBuf.WriteString(fmt.Sprintln(format))
	} else {
		l.StrBuf.WriteString(fmt.Sprintf(format, args...))
		l.StrBuf.WriteString(fmt.Sprintln())
	}
}

func (l *Lines) String() string {
	return l.StrBuf.String()
}

func SumMd5FromFile(fileName string) string {
	file, err := os.Open(fileName)
	CheckErrorPanic(err)
	m := md5.New()
	_, err = io.Copy(m, file)
	CheckErrorPanic(err)
	Md5Str := fmt.Sprintf("%x", m.Sum(nil))
	return Md5Str
}

// 以保护方式运行函数
func ProtectRun(entry func()) {
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		err := recover()
		switch err.(type) {
		case runtime.Error: // 运行时错误
			fmt.Println("出现运行时错误: ", err)
		default: // 非运行时错误
			fmt.Println("出现错误: ", err)
		}
	}()
	entry()
}

// 判断某项标志是否在标志列表里
func IsInFlag(item string, items []string) bool {
	log.Debugf("判断 %#v 是否在 %#v 中\n", item, items)
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}

// 判断某项标志是否在标志列表里, 如果在标志列表里则删除并返回true
func IfInFlagThenPop(item string, items []string) ([]string, bool) {
	log.Debugf("判断 %#v 是否在 %#v 中\n", item, items)
	isIncluded := false
	newItems := make([]string, 0)
	for _, i := range items {
		if i == item {
			isIncluded = true
		} else {
			newItems = append(newItems, i)
		}
	}
	return newItems, isIncluded
}

// 如果字符串不为空则返回原值, 否则返回默认字符串
func StringDefault(s string, d string) string {
	if s == "" {
		return d
	} else {
		return s
	}
}

// 如果字符串不为空则将原值转换为int返回, 否则返回默认int
func StringDefaultInt(s string, d int) int {
	log.Debugf("获取原始字符串: %s", s)
	if s == "" {
		return d
	} else {
		i, err := strconv.Atoi(s)
		CheckErrorPanic(err)
		return i
	}
}

// 如果字符串不为空则将原值转换为int返回, 否则返回默认int
func StringDefaultInt64(s string, d int64) int64 {
	if s == "" {
		return d
	} else {
		i64, err := strconv.ParseInt(s, 10, 64)
		CheckErrorPanic(err)
		return i64
	}
}

// 如果字符串不为空则将原值转换为float64返回, 否则返回默认float64
func StringDefaultFloat64(s string, d float64) float64 {
	if s == "" {
		return d
	} else {
		f64, err := strconv.ParseFloat(s, 64)
		CheckErrorPanic(err)
		return f64
	}
}

func GetLocalIp(connectTo string) string {
	conn, _ := net.Dial("tcp", connectTo)
	defer conn.Close()
	addr := conn.LocalAddr().String()
	return strings.Split(addr, ":")[0]
}
