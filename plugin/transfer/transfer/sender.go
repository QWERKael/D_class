package transfer

import (
	"errors"
	"io"
	"net"
	"os"
	"path/filepath"
)

func sendFile(path string, fileName string, conn net.Conn) error {
	defer conn.Close()
	fs, err := os.Open(path)
	defer fs.Close()
	if err != nil {
		log.Errorf("无法打开文件: ", err.Error())
		return err
	}
	buf := make([]byte, 1<<20)
	for {
		//  打开之后读取文件
		n, err := fs.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Infof("%s 已接收完毕", fileName)
				return nil
			}
			log.Errorf("无法读取文件: ", err.Error())
			return err
		}
		if n == 0 {
			log.Infof("%s 已接收完毕", fileName)
			return nil
		}

		//  发送文件
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Errorf("文件发送错误: ", err.Error())
			return err
		}
		log.Debugf("已发送 %d 字节", n)
	}
}

func Send(path string, sendTo string) error {
	//info, err := os.Stat(path)
	//if err != nil {
	//	log.Errorf("无法获取文件状态: ", err.Error())
	//	return err
	//}
	// 发送文件名
	conn, err := net.Dial("tcp", sendTo)
	if err != nil {
		log.Errorf("无法连接到目标服务器: ", err.Error())
		return err
	}
	defer conn.Close()
	//fileName := info.Name()
	fileName := filepath.Base(path)
	log.Debugf("请求发送文件 [%s]", fileName)
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		log.Errorf("发送文件名失败: ", err.Error())
		return err
	}
	// 接受到是不是ok
	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Errorf("未收到确认通知: ", err.Error())
		return err
	}
	if "ok" == string(buf[:n]) {
		log.Infof("开始发送文件 [%s]", fileName)
		err = sendFile(path, fileName, conn)
		if err != nil {
			log.Errorf("发送文件错误: ", err.Error())
			return err
		}
		return nil
	}
	return errors.New("收到未知的反馈: " + string(buf[:n]))
}

//func main() {
//	for {
//
//		fmt.Println("请输入一个全路径的文件,比如,D:\\a.jpg")
//		//  获取命令行参数
//		var path string
//		fmt.Scan(&path)
//		fmt.Printf("获取到文件名: %s", path)
//		// 获取文件名,
//		for _, p := range strings.Split(path, " ") {
//			go Send(p)
//		}
//
//		// 如果是ok,那么开启一个连接,发送文件
//	}
//}
