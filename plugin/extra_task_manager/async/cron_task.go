package async

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

//type CronPersistence struct {
//	CronInfo   CronInfo            `json:"CronInfo"`
//	Type       TaskType            `json:"Type"`
//	State      TaskState           `json:"State"`
//	Request    pb.CommonCmdRequest `json:"Request"`
//	Reply      pb.CommonCmdReply   `json:"Reply"`
//	FinishTime time.Time           `json:"FinishTime"`
//	NotifyMsg  string              `json:"NotifyMsg"`
//}

func SaveCronToFile(ats []AsyncTask, savePath string) error {
	log.Debugf("持久化定时任务到: %s", savePath)
	b, err := json.MarshalIndent(ats, "", " ")
	if err != nil {
		return errors.New("序列化异步任务失败失败: " + err.Error())
	}
	var file *os.File
	file, err = os.OpenFile(savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return errors.New("打开定时任务持久化文件失败: " + err.Error())
	}
	defer file.Close()
	_, err = file.Write(b)
	if err != nil {
		return errors.New("打开定时任务持久化文件失败: " + err.Error())
	}
	//atJson := fmt.Sprintf("%s", b)
	//fmt.Printf("\nMarshal is \n%s\n", cpJson)
	return nil
}

func LoadCronFile(savePath string) ([]*AsyncTask, error) {
	log.Debugf("从 %s 持久化定时任务", savePath)
	b, err := ioutil.ReadFile(savePath)
	if err != nil {
		return nil, errors.New("打开定时任务持久化文件失败: " + err.Error())
	}
	var ats []*AsyncTask
	err = json.Unmarshal(b, &ats)
	if err != nil {
		return nil, errors.New("反序列化异步任务失败失败: " + err.Error())
	}
	return ats, nil
}
