package transfer

import (
	"atk_D_class/logger"
	"io"
	"net"
	"os"
	"path/filepath"
)

var log = logger.Logger

func revFile(recvDir string, fileName string, conn net.Conn) error {
	defer conn.Close()
	recvPath := filepath.Join(recvDir, fileName)
	fs, err := os.Create(recvPath)
	defer fs.Close()
	if err != nil {
		log.Errorf("创建文件失败: %s", err.Error())
		return err
	}

	// 拿到数据
	buf := make([]byte, 1<<20)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Infof("%s 已接收完毕", fileName)
				return nil
			}
			log.Errorf("接收数据失败: %s", err.Error())
			return err
		}
		if n == 0 {
			log.Infof("%s 已接收完毕", fileName)
			return nil
		}
		_, err = fs.Write(buf[:n])
		if err != nil {
			log.Errorf("写入文件失败: %s", err.Error())
			return err
		}
		log.Debugf("已接受到 %d 字节", n)
	}
}
func SetReceiver(recvServ net.Listener, recvDir string) {
	// 接受文件名
	for {
		conn, err := recvServ.Accept()
		defer conn.Close()
		if err != nil {
			log.Errorf("接受连接出错: %s", err.Error())
			return
		}
		go func(recvDir string) {
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				log.Errorf("获取文件名失败: %s", err.Error())
				return
			}
			// 拿到了文件的名字
			fileName := string(buf[:n])
			log.Infof("接收到文件名: %s", fileName)
			// 返回ok
			_, err = conn.Write([]byte("ok"))
			if err != nil {
				log.Errorf("发送文件名确认信息失败: %s", err.Error())
				return
			}
			// 接收文件,
			err = revFile(recvDir, fileName, conn)
			if err != nil {
				log.Errorf("接收文件失败: %s", err.Error())
				return
			}
		}(recvDir)
	}
}
