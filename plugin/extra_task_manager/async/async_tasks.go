package async

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"context"
	"fmt"
	"gopkg.in/robfig/cron.v3"
	//"github.com/robfig/cron/v3"
	"strconv"
	"strings"
	"sync"
	"time"
)

var log = logger.Logger

type TaskState int32

const (
	Unknown TaskState = 0
	Prepare TaskState = 1
	Running TaskState = 2
	Success TaskState = 3
	Fail    TaskState = 4
)

func (ts TaskState) ToString() string {
	return []string{
		"Unknown",
		"Prepare",
		"Running",
		"Success",
		"Fail",
	}[ts]
}

func StrToTaskState(s string) TaskState {
	m := map[string]TaskState{
		"unknown": Unknown,
		"prepare": Prepare,
		"running": Running,
		"success": Success,
		"fail":    Fail,
	}
	if state, ok := m[strings.ToLower(s)]; ok {
		return state
	}
	return Unknown
}

type CronState int32

const (
	CronDown CronState = 0
	CronUP   CronState = 1
)

func (cs CronState) ToString() string {
	return []string{
		"CronDown",
		"CronUP",
	}[cs]
}

type TriggerState int32

const (
	TriggerDown TriggerState = 0
	TriggerUP   TriggerState = 1
)

func (ts TriggerState) ToString() string {
	return []string{
		"TriggerDown",
		"TriggerUP",
	}[ts]
}

type TaskType int32

const (
	TAsync   TaskType = 0
	TCron    TaskType = 1
	TTrigger TaskType = 2
)

func (tt TaskType) ToString() string {
	return []string{
		"TAsync",
		"TCron",
		"TTrigger",
	}[tt]
}

type CronInfo struct {
	CronId       cron.EntryID `json:"-"`
	CronSchedule string       `json:"CronSchedule"`
	CronState    CronState    `json:"CronState"`
}

type TriggerInfo struct {
	TriggerId     int
	TriggerById   int
	TriggerByName string
	//TriggerByNotify *chan int
	TriggerByState TaskState
	TriggerState   TriggerState
	TriggersNumber int
	RemainTimes    int
}

type AsyncTask struct {
	AsyncTaskId int                  `json:"-"`
	CronInfo    CronInfo             `json:"CronInfo"`
	TriggerInfo TriggerInfo          `json:"-"`
	Type        TaskType             `json:"Type"`
	State       TaskState            `json:"State"`
	Request     *pb.CommonCmdRequest `json:"Request"`
	Reply       *pb.CommonCmdReply   `json:"Reply"`
	FinishTime  time.Time            `json:"FinishTime"`
	Aka         string               `json:"Aka"`
	NotifyMsg   string               `json:"NotifyMsg"`
}

func Request2String(req *pb.CommonCmdRequest) string {
	s := req.Plugin
	if req.Cmd != "" {
		s += " " + req.Cmd
	}
	if len(req.SubCmd) > 0 {
		s += " " + strings.Join(req.SubCmd, " ")
	}
	if len(req.Flags) > 0 {
		s += " "
		s += "-" + strings.Join(req.Flags, " -")
	}
	if len(req.Args) > 0 {
		s += " "
		for k, v := range req.Args {
			s += "--" + k + "=" + v
		}
	}
	return s
}

func BuildAsyncTask(task *tasks.Task, tt TaskType, aka string) AsyncTask {
	//notify := make(chan Notify, 5)
	var (
		cmd    string
		subcmd []string
	)

	switch len(task.SubCmd) {
	case 0:
		cmd = "DefaultFunc"
		subcmd = nil
	case 1:
		cmd = task.SubCmd[0]
		subcmd = nil
	default:
		cmd = task.SubCmd[0]
		subcmd = task.SubCmd[1:]
	}

	at := AsyncTask{
		AsyncTaskId: 0,
		Type:        tt,
		State:       Unknown,
		Request: &pb.CommonCmdRequest{
			Type:   0,
			Plugin: task.Cmd,
			Cmd:    cmd,
			SubCmd: subcmd,
			Flags:  task.Flags,
			Args:   task.Args,
		},
		Reply:     nil,
		Aka:       aka,
		NotifyMsg: "",
	}
	return at
}

type StateList struct {
	NextId int
	List   map[int]TaskState
	Mu     sync.Mutex
}

func (sl *StateList) Insert(at *AsyncTask) int {
	var id int
	sl.Mu.Lock()
	id = sl.NextId
	sl.NextId++
	sl.Mu.Unlock()
	at.AsyncTaskId = id
	at.State = Prepare
	sl.List[id] = Prepare
	return id
}

func (sl *StateList) ToTable() *pb.CommonCmdReply {
	header := []string{"Id", "State"}
	body := make([][]string, 0)
	for k, v := range sl.List {
		body = append(body, []string{strconv.Itoa(k), v.ToString()})
	}
	footer := []string{"??????", fmt.Sprintf("%d ???", len(body))}
	reply := prompt.ToTable(header, footer, body, 0)
	return reply
}

type Notify struct {
	NotifyFromId int
	TaskState    TaskState
	Msg          string
}

// ????????????????????????:1.???????????? 2.???????????? 3.????????????
// ?????????????????????????????????
type TaskManager struct {
	AsyncTaskList map[int]*AsyncTask
	PrepareList   chan *AsyncTask
	StateList     StateList
	AkaList       map[string]int
	Notify        chan Notify
	NotifyList    map[int]map[int]struct{}
	CC            pb.CommanderClient
	Ctx           context.Context
	Cancel        context.CancelFunc
	Cron          *cron.Cron
}

func (tm *TaskManager) State(detailFlag bool) *pb.CommonCmdReply {
	header := []string{"Id", "AKA", "Type", "State", "Finish Time"}
	if detailFlag {
		header = append(header, "Detail")
	}
	body := make([][]string, 0)
	for k, v := range tm.AsyncTaskList {
		row := []string{
			strconv.Itoa(k),
			v.Aka,
			v.Type.ToString(),
			v.State.ToString(),
			v.FinishTime.Format("2006-01-02 15:04:05")}
		if detailFlag {
			row = append(row, Request2String(v.Request))
		}
		body = append(body, row)
	}
	footer := []string{"??????", fmt.Sprintf("%d ???", len(body)), "", ""}
	reply := prompt.ToTable(header, footer, body, 0)
	return reply
}

func MakeTaskManager(capacity int, cc pb.CommanderClient, ctx context.Context, cancel context.CancelFunc) *TaskManager {
	return &TaskManager{
		AsyncTaskList: make(map[int]*AsyncTask),
		PrepareList:   make(chan *AsyncTask, capacity),
		StateList: StateList{
			NextId: 0,
			List:   make(map[int]TaskState),
		},
		AkaList:    make(map[string]int),
		Notify:     make(chan Notify, 10),
		NotifyList: make(map[int]map[int]struct{}),
		CC:         cc,
		Ctx:        ctx,
		Cancel:     cancel,
	}
}

func (tm *TaskManager) AsyncRunner() {
	log.Debugf("???????????????????????????...")
	for at := range tm.PrepareList {
		go func(at AsyncTask) {
			tm.AsyncTaskList[at.AsyncTaskId] = &at
			log.Debugf("??????????????????: %#v", at.Request)
			var err error
			at.State = Running
			at.Reply, err = tm.CC.CommonCmd(tm.Ctx, at.Request)
			at.FinishTime = time.Now()
			if err != nil {
				at.State = Fail
				at.Reply = &pb.CommonCmdReply{ResultMsg: err.Error()}
				log.Debugf("AsyncRunner: ??????????????????! %s", err.Error())
			} else if at.Reply.Status == pb.CommonCmdReply_Err {
				at.State = Fail
				log.Debugf("AsyncRunner: ??????????????????! %s", at.Reply.ResultMsg)
			} else {
				at.State = Success
				log.Debugf("AsyncRunner: ??????????????????!")
			}
			log.Debugf("????????????????????? [%d], ???????????????...", at.AsyncTaskId)

			tm.Notify <- Notify{
				NotifyFromId: at.AsyncTaskId,
				TaskState:    at.State,
				Msg:          at.NotifyMsg,
			}
			log.Debugf("???????????????????????????")
			return
		}(*at)
		tm.StateList.List[at.AsyncTaskId] = Running
	}
}

func (tm *TaskManager) TriggerRunner() {
	log.Debugf("??????????????????????????????...")
	for notify := range tm.Notify {
		notifyToList := tm.NotifyList[notify.NotifyFromId]
		for notifyToId := range notifyToList {
			at := tm.AsyncTaskList[notifyToId]
			// ???????????????????????????????????????, ??????
			if at.Type != TTrigger ||
				// ????????????????????????????????????, ??????
				at.TriggerInfo.TriggerState == TriggerDown {
				continue
			} else if at.TriggerInfo.TriggersNumber > 0 &&
				at.TriggerInfo.RemainTimes < 1 {
				// ????????????????????????????????????????????????, ???????????????????????????
				at.TriggerInfo.TriggerState = TriggerDown
				continue
			} else if at.TriggerInfo.TriggerByState != notify.TaskState {
				// ???????????????????????????????????????????????????????????????, ??????
				continue
			}
			// ????????????????????????????????????
			tm.PrepareList <- at
			log.Debugf("??????????????? [%d] ?????????????????????", notifyToId)
			if at.TriggerInfo.TriggersNumber > 0 {
				at.TriggerInfo.RemainTimes--
				if at.TriggerInfo.RemainTimes < 1 {
					// ????????????????????????????????????????????????, ???????????????????????????
					at.TriggerInfo.TriggerState = TriggerDown
				}
			}
		}
	}
}

func (tm *TaskManager) AddAsyncTask(at *AsyncTask) (int, error) {
	var err error
	asyncTaskId := tm.StateList.Insert(at)

	// ????????????id????????????
	if at.Aka != "" {
		tm.AkaList[at.Aka] = at.AsyncTaskId
	}

	// ???????????????????????????????????????, ????????????????????????
	tm.AsyncTaskList[at.AsyncTaskId] = at
	if at.Type == TCron {
		at.CronInfo.CronId, err = tm.Cron.AddFunc(at.CronInfo.CronSchedule, func() {
			log.Debugf("?????????????????? [%d](??????????????? [%d])", at.AsyncTaskId, at.CronInfo.CronId)
			tm.PrepareList <- at
		})
		if err != nil {
			at.CronInfo.CronState = CronDown
			return 0, err
		}
	}
	return asyncTaskId, nil
}

func (tm *TaskManager) GetAsyncTaskId(name string) int {
	if name == "" {
		return -1
	} else if id, err := strconv.Atoi(name); err == nil {
		return id
	} else {
		if id, ok := tm.AkaList[name]; ok {
			return id
		} else {
			return -1
		}
	}
}
