package prompt

import (
	"atk_D_class/pb"
	"atk_D_class/utils"
	"fmt"
)

func uploadPath(flags []string) string {
	switch {
	case utils.IsInFlag("plugin", flags):
		log.Debugln("上传插件...")
		return "plugin"
	case utils.IsInFlag("update", flags):
		log.Debugln("上传更新...")
		return ""
	case utils.IsInFlag("script", flags):
		log.Debugln("上传脚本...")
		return "scripts"
	case utils.IsInFlag("config", flags):
		log.Debugln("上传配置文件...")
		return "config"
	default:
		log.Debugln("上传文件...")
		return "files"
	}
}

// 校验传输结果
func checkUploadResult(recvFi *pb.TransferInfo, localPath string, fileName string) {
	if recvFi.State == pb.TransferInfo_Complete {
		md5Str := utils.SumMd5FromFile(localPath)
		if recvFi.Md5 == md5Str {
			fmt.Printf("文件%s传输成功\n", fileName)
		} else {
			fmt.Printf("文件%s传输失败, MD5不匹配\n本地文件md5校验值: %s\n远程文件md5校验值: %s\n",
				fileName, md5Str, recvFi.Md5)
		}
	} else if recvFi.State == pb.TransferInfo_Error {
		fmt.Printf("文件传输失败, 错误: %s", recvFi.ErrorMsg)
	} else {
		fmt.Printf("文件传输失败, 未知错误")
	}
}
