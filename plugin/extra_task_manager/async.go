package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/plugin/extra_task_manager/async"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	//"github.com/robfig/cron/v3"
	"gopkg.in/robfig/cron.v3"
	"path/filepath"
	"strconv"
	"strings"
)

var c *cron.Cron

var log = logger.Logger
var taskManager *async.TaskManager
var CC pb.CommanderClient
var ctx context.Context
var cancel context.CancelFunc
var workDir string

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
desc: "Async Plugin\n支持异步操作的插件\n
async run \n\t启动异步任务管理器\n
async [plugin cmd ...] \n\t在要执行的命令前添加 async 前缀, 将命令作为异步命令执行, 并返回异步命令的Id, 根据Id查看异步命令的执行情况\n
async [plugin cmd ...] --cron='* * * * * *' \n\t将异步任务作为定时任务执行, 根据cron规则 [秒 分 时 日 月 周] 定时执行\n
async [plugin cmd ...] --trigger-by=[id|name] --trigger-by-status=[success] --triggers-number=[1] \n\t将异步任务作为触发器执行, 根据trigger-by指定触发该任务
所依赖的任务Id或者别名, 根据trigger-by-status指定触发该任务所依赖的状态, 根据triggers-number指定触发该任务所重复的次数(0表示无限重复, 默认是1)\n
async state [-detail]\n\t查看当前任务列表, 使用-detail查看任务详情\n
async get [Id]\n\t获取异步任务执行的结果\n
async pop [Id]\n\t获取异步任务执行的结果, 并从列表中删除该任务\n
async cron [Id] -del\n\t查看定时任务的详细信息, 可以使用-del删除定时任务\n
async cron [save|load]\n\t持久化定时任务/从持久化文件加载定时任务\n
async trigger [Id] -del\n\t查看触发器任务的详细信息, 可以使用-del删除触发器任务\n"
yess:
  async:
    type: "Prefix"
    desc: "提供异步任务相关操作"
    yess:
      run:
        type: "Cmd"
        desc: "启动异步任务管理器"
      --cron=:
        type: "ArgKey"
        desc: "将当前任务作为定时任务执行, 并指定定时任务的定时策略"
      --trigger-by=:
        type: "ArgKey"
        desc: "将当前任务作为触发器任务执行, 并指定触发器任务所依赖的任务Id或者别名"
        yess:
          --trigger-by-status=:
            type: "ArgKey"
            desc: "指定触发器任务所依赖的任务状态"
          --triggers-number=:
            type: "ArgKey"
            desc: "指定触发器任务所重复的任务次数(0表示无限重复)"
      state:
        type: "Cmd"
        desc: "查看当前任务列表"
        yess:
          -detail:
            type: "Flag"
            desc: "显示任务详情"
      get:
        type: "Cmd"
        desc: "获取异步任务执行的结果"
      pop:
        type: "Cmd"
        desc: "获取异步任务执行的结果, 并从列表中删除该任务"
      cron:
        type: "Cmd"
        desc: "查看定时任务的详细信息"
        yess:
          "-del":
            type: "Flag"
            desc: "删除定时任务"
          "-save":
            type: "Flag"
            desc: "持久化定时任务"
          "-load":
            type: "Flag"
            desc: "从持久化文件加载定时任务"
      trigger:
        type: "Cmd"
        desc: "查看触发器任务的详细信息"
        yess:
          "-del":
            type: "Flag"
            desc: "删除触发器任务"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

func Run(task tasks.Task) (string, error) {
	workDir = task.Context["workDir"]
	var capacity int
	if v, ok := task.Args["cap"]; ok {
		capacity = utils.StringDefaultInt(v, 100)
		delete(task.Args, "cap")
	}
	port := strings.Split(task.Context["listenAddr"], ":")[1]
	addr := fmt.Sprintf("127.0.0.1:%s", port)
	ctx, cancel = context.WithCancel(context.Background())

	// 为ctx添加用户认证
	auth, err := utils.MakeCred("simple", "admin", "admin")
	if err != nil {
		auth = ""
	}
	ctx = metadata.NewOutgoingContext(ctx,
		metadata.Pairs(
			"auth", auth,
		),
	)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	utils.CheckErrorPanic(errors.WithMessage(err, "建立代理连接到指定地址失败: "+addr))
	client := pb.NewCommanderClient(conn)
	log.Infoln("建立代理连接到", addr, "...")
	CC = client
	taskManager = async.MakeTaskManager(capacity, CC, ctx, cancel)
	// 依次运行异步任务队列
	go taskManager.AsyncRunner()
	// 监听Notify消息, 运行相应的触发器
	go taskManager.TriggerRunner()
	// 启动定时任务管理器
	log.Debugf("启动定时任务管理器")
	c = cron.New(cron.WithSeconds())
	taskManager.Cron = c
	c.Start()

	return "异步任务队列已开启", nil
}

func DefaultFunc(task tasks.Task) (string, error) {
	var (
		err         error
		s           string
		asyncTaskId int
		aka         string
		ok          bool
	)
	var msg = ""
	if taskManager == nil {
		s, err = Run(task)
		if err != nil {
			return "创建异步任务队列失败: " + err.Error(), nil
		}
		msg += s + "\n"
	}

	// 获取别名参数
	if aka, ok = task.Args["aka"]; ok {
		log.Debugf("任务有别名: %s", aka)
		delete(task.Args, "aka")
		// 别名不能为纯数字, 以和id区分开
		if number, err := strconv.Atoi(aka); err == nil {
			log.Errorf("aka的值不能为纯数字: %d", number)
			return "", errors.New(fmt.Sprintf("aka的值不能为纯数字: %d", number))
		}
	}

	if cronStr, ok := task.Args["cron"]; ok {
		// 执行定时任务
		log.Debugf("添加定时任务")
		delete(task.Args, "cron")

		at := async.BuildAsyncTask(&task, async.TCron, aka)
		at.CronInfo = async.CronInfo{
			CronId:       0,
			CronSchedule: cronStr,
			CronState:    async.CronUP,
		}
		//asyncTaskId = taskManager.StateList.Insert(&at)
		//// 将触发器加入到异步任务队列, 此处不是执行队列
		//taskManager.AsyncTaskList[at.AsyncTaskId] = &at

		//at.CronInfo.CronId, err = c.AddFunc(cronStr, func() {
		//	log.Debugf("唤醒定时任务 [%d](定时器编号 [%d])", at.AsyncTaskId, at.CronInfo.CronId)
		//	taskManager.PrepareList <- &at
		//})
		//if err != nil {
		//	at.CronInfo.CronState = async.CronDown
		//	msg += fmt.Sprintf("定时任务添加失败: %s", err.Error())
		//	return msg, nil
		//}
		asyncTaskId, err := taskManager.AddAsyncTask(&at)
		if err != nil {
			msg += fmt.Sprintf("定时任务添加失败: %s", err.Error())
			return msg, nil
		}

		//taskManager.AsyncTaskList[asyncTaskId].CronId = int(CronId)
		msg += fmt.Sprintf("定时任务已添加, 任务编号 [%d], 定时器编号 [%d]\n", asyncTaskId, at.CronInfo.CronId)

		if err = saveCron(); err != nil {
			log.Warnf("持久化定时任务失败: %s", err.Error())
			msg += fmt.Sprintf("持久化定时任务失败: %s", err.Error())
		}
	} else if triggerBy, ok := task.Args["trigger-by"]; ok {
		// 执行触发器任务
		triggeredById := taskManager.GetAsyncTaskId(triggerBy)
		//triggeredById, err := strconv.Atoi(triggerBy)
		//if err != nil {
		//	msg += fmt.Sprintf("触发器添加失败: %s", err.Error())
		//	return msg, nil
		//}
		triggersNumber := utils.StringDefaultInt(task.Args["triggers-number"], 1)
		triggerByState := async.Success
		if triggerByStatus, ok := task.Args["trigger-by-status"]; ok {
			triggerByState = async.StrToTaskState(triggerByStatus)
			if triggerByState == async.Unknown {
				msg += fmt.Sprintf("触发器添加失败: %s", "未知的触发状态")
				return msg, nil
			}
		}
		log.Debugf("添加触发器")
		delete(task.Args, "trigger-by-status")
		delete(task.Args, "triggers-number")

		at := async.BuildAsyncTask(&task, async.TTrigger, aka)
		at.TriggerInfo = async.TriggerInfo{
			TriggerId:      0,
			TriggerById:    triggeredById,
			TriggerByName:  at.Aka,
			TriggerByState: triggerByState,
			TriggerState:   async.TriggerUP,
			TriggersNumber: triggersNumber,
			RemainTimes:    triggersNumber,
		}
		asyncTaskId = taskManager.StateList.Insert(&at)
		at.TriggerInfo.TriggerId = asyncTaskId
		// 将别名和id对应起来
		if at.Aka != "" {
			taskManager.AkaList[at.Aka] = at.AsyncTaskId
		}
		// 将触发器加入到异步任务队列, 此处不是执行队列
		taskManager.AsyncTaskList[at.AsyncTaskId] = &at
		// 将触发器任务Id加入到依赖任务的触发列表里
		if _, ok = taskManager.NotifyList[triggeredById]; ok {
			taskManager.NotifyList[triggeredById][asyncTaskId] = struct{}{}
		} else {
			taskManager.NotifyList[triggeredById] = make(map[int]struct{})
			taskManager.NotifyList[triggeredById][asyncTaskId] = struct{}{}
		}
		log.Debugf("触发器已添加, 正在监听任务 [%d]", triggeredById)
		msg += fmt.Sprintf("触发器已添加, 任务编号 [%d], 触发器编号 [%d]", asyncTaskId, at.TriggerInfo.TriggerId)

	} else {
		// 执行异步任务
		log.Debugf("添加异步任务")
		at := async.BuildAsyncTask(&task, async.TAsync, aka)
		asyncTaskId := taskManager.StateList.Insert(&at)
		// 将别名和id对应起来
		if at.Aka != "" {
			taskManager.AkaList[at.Aka] = at.AsyncTaskId
		}
		// 将触发器加入到异步任务队列, 此处不是执行队列
		taskManager.AsyncTaskList[at.AsyncTaskId] = &at
		taskManager.PrepareList <- &at
		msg += fmt.Sprintf("异步任务已添加, 任务编号 [%d]", asyncTaskId)
	}

	return msg, nil
}

func State(task tasks.Task) (*pb.CommonCmdReply, error) {
	detailFlag := false
	if taskManager == nil {
		return nil, errors.New("任务管理器尚未启动")
	}
	detailFlag = utils.IsInFlag("detail", task.Flags)
	reply := taskManager.State(detailFlag)
	return reply, nil
}

func saveCron() error {
	log.Debugf("SAVE定时任务...")
	cronSaveList := make([]async.AsyncTask, 0)
	for _, v := range taskManager.AsyncTaskList {
		cronSaveList = append(cronSaveList, *v)
	}
	err := async.SaveCronToFile(cronSaveList, filepath.Join(workDir, "cron.save"))
	if err != nil {
		return err
	}
	return nil
}

func Cron(task tasks.Task) (*pb.CommonCmdReply, error) {
	log.Debugf("任务: %#v", task)
	delFlag := utils.IsInFlag("del", task.Flags)
	saveFlag := utils.IsInFlag("save", task.Flags)
	loadFlag := utils.IsInFlag("load", task.Flags)
	if len(task.SubCmd) < 1 {
		if saveFlag {
			err := saveCron()
			if err != nil {
				return nil, errors.New(fmt.Sprintf("定时任务save失败: %s", err.Error()))
			}
			return &pb.CommonCmdReply{
				ResultMsg: fmt.Sprintf("定时任务save成功!"),
				Status:    pb.CommonCmdReply_Ok,
			}, nil
		} else if loadFlag {
			ats, err := async.LoadCronFile(filepath.Join(workDir, "cron.save"))
			if err != nil {
				return nil, errors.New(fmt.Sprintf("定时任务load失败: %s", err.Error()))
			}
			for _, at := range ats {
				asyncTaskId, err := taskManager.AddAsyncTask(at)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("定时任务[%d]加载失败: %s", asyncTaskId, err.Error()))
				}
				fmt.Printf("\nUnmarshal is \n%#v\n", *at.Request)
			}

			return &pb.CommonCmdReply{
				ResultMsg: fmt.Sprintf("定时任务load成功!"),
				Status:    pb.CommonCmdReply_Ok,
			}, nil
		}
	}
	//asyncTaskId := utils.StringDefaultInt(task.SubCmd[0], -1)
	asyncTaskId := taskManager.GetAsyncTaskId(task.SubCmd[0])
	if at, ok := taskManager.AsyncTaskList[asyncTaskId]; ok {
		if at.Type != async.TCron {
			return nil, errors.New("该任务不是定时任务")
		}

		// 删除定时任务
		if delFlag {
			// 将定时任务状态置为Down
			at.CronInfo.CronState = async.CronDown
			// 从定时器中删除
			c.Remove(at.CronInfo.CronId)
			// 从异步任务队列和状态队列中删除
			delete(taskManager.AsyncTaskList, asyncTaskId)
			delete(taskManager.StateList.List, asyncTaskId)

			return &pb.CommonCmdReply{
				ResultMsg: fmt.Sprintf("定时任务 [%d](定时器编号 [%d]) 已删除", asyncTaskId, at.CronInfo.CronId),
				Status:    pb.CommonCmdReply_Ok,
			}, nil
		}

		// 返回定时任务信息
		header := []string{"定时任务属性", "参数"}
		body := [][]string{
			{"任务编号", strconv.Itoa(asyncTaskId)},
			{"任务别名", at.Aka},
			{"定时器编号", strconv.Itoa(int(at.CronInfo.CronId))},
			{"执行计划", at.CronInfo.CronSchedule},
			{"运行状态", at.CronInfo.CronState.ToString()},
		}
		reply := prompt.ToTable(header, []string{}, body, 0)
		return reply, nil
	}
	return nil, errors.New("未找到该任务")
}

func Trigger(task tasks.Task) (*pb.CommonCmdReply, error) {
	log.Debugf("任务: %#v", task)
	delFlag := utils.IsInFlag("del", task.Flags)
	//asyncTaskId := utils.StringDefaultInt(task.SubCmd[0], -1)
	asyncTaskId := taskManager.GetAsyncTaskId(task.SubCmd[0])
	if at, ok := taskManager.AsyncTaskList[asyncTaskId]; ok {
		if at.Type != async.TTrigger {
			return nil, errors.New("该任务不是触发器任务")
		}

		triggeredById := at.TriggerInfo.TriggerById
		// 删除触发器任务
		if delFlag {
			// 将触发器状态置为Down
			at.TriggerInfo.TriggerState = async.TriggerDown
			// 从通知列表的对应关系中删除
			delete(taskManager.NotifyList[triggeredById], asyncTaskId)
			delete(taskManager.NotifyList, asyncTaskId)
			// 从异步任务队列和状态队列中删除
			delete(taskManager.AsyncTaskList, asyncTaskId)
			delete(taskManager.StateList.List, asyncTaskId)

			return &pb.CommonCmdReply{
				ResultMsg: fmt.Sprintf("触发器任务 [%d] 已删除", asyncTaskId),
				Status:    pb.CommonCmdReply_Ok,
			}, nil
		}

		// 返回定时任务信息
		header := []string{"触发器任务属性", "参数"}
		body := [][]string{
			{"任务编号", strconv.Itoa(asyncTaskId)},
			{"任务别名", at.Aka},
			{"任务状态", at.TriggerInfo.TriggerState.ToString()},
			{"依赖任务编号", strconv.Itoa(triggeredById)},
			{"依赖任务别名", at.TriggerInfo.TriggerByName},
			{"依赖任务状态", at.TriggerInfo.TriggerByState.ToString()},
			{"执行次数", strconv.Itoa(at.TriggerInfo.TriggersNumber)},
			{"剩余执行次数", strconv.Itoa(at.TriggerInfo.RemainTimes)},
		}
		reply := prompt.ToTable(header, []string{}, body, 0)
		return reply, nil
	}
	return nil, errors.New("未找到该任务")
}

func Get(task tasks.Task) (*pb.CommonCmdReply, error) {
	id := taskManager.GetAsyncTaskId(task.SubCmd[0])
	if taskManager.AsyncTaskList[id].State != async.Success {
		return &pb.CommonCmdReply{ResultMsg: "任务尚未就绪, 请等待...", Status: pb.CommonCmdReply_Ok}, nil
	}
	reply := taskManager.AsyncTaskList[id].Reply
	return reply, nil
}

func Pop(task tasks.Task) (*pb.CommonCmdReply, error) {
	//id := utils.StringDefaultInt(task.SubCmd[0], -1)
	id := taskManager.GetAsyncTaskId(task.SubCmd[0])
	if taskManager.AsyncTaskList[id].State != async.Success {
		return &pb.CommonCmdReply{ResultMsg: "任务尚未就绪, 请等待...", Status: pb.CommonCmdReply_Ok}, nil
	}
	reply := taskManager.AsyncTaskList[id].Reply
	if taskManager.AsyncTaskList[id].Type != async.TAsync {
		reply.ResultMsg += "\n指定任务为特殊任务, 无法删除该任务"
		reply.Status = pb.CommonCmdReply_Err
		return reply, nil
	}
	delete(taskManager.AsyncTaskList, id)
	delete(taskManager.StateList.List, id)
	return reply, nil
}
