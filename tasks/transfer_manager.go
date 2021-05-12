package tasks

import (
	"atk_D_class/pb"
	"os"
	"sync"
)

type FileInfo struct {
	TransferType  pb.TransferInfo_TransferType
	TransferState pb.TransferInfo_TransferState
	FileName      string
	FilePath      string
	File          *os.File
	TransferId    int32
	ErrorMsg      string
	Md5           string
}

type FileManager struct {
	FileInfos []FileInfo
}

func (f *FileInfo) Assemble() pb.TransferInfo {
	return pb.TransferInfo{Type: f.TransferType, State: f.TransferState, FileName: f.FileName,
		FilePath: f.FilePath, TransferId: f.TransferId, ErrorMsg: f.ErrorMsg, Md5: f.Md5}
}

func Unpack(ti *pb.TransferInfo) FileInfo {
	return FileInfo{TransferType: ti.Type, TransferState: ti.State, FileName: ti.FileName,
		FilePath: ti.FilePath, TransferId: ti.TransferId, ErrorMsg: ti.ErrorMsg, Md5: ti.Md5}
}

func (fm *FileManager) Add(fi FileInfo) FileInfo {
	if fi.TransferId == 0 {
		var mut sync.Mutex
		mut.Lock()
		log.Debugf("%#v", fm)
		fi.TransferId = int32(len(fm.FileInfos))
		fm.FileInfos = append(fm.FileInfos, fi)
		mut.Unlock()
	}
	return fi
}

func (fm *FileManager) BuildFile(fp *os.File, chunk *pb.Chunks) error {
	n, err := fp.Write(chunk.Content[:chunk.Size])
	if err != nil {
		return err
	}
	log.Debugf("写入数据 %d 字节\n", n)
	return nil
}
