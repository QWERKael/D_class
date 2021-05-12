package async

import (
	"atk_D_class/pb"
	"atk_D_class/tasks"
	"fmt"
	"testing"
)

func TestSaveCronToFile(t *testing.T) {
	task := tasks.Task{
		Plugin:  "async",
		Cmd:     "cmd",
		SubCmd:  []string{"ls"},
		Flags:   nil,
		Args:    nil,
		Context: nil,
	}
	at := BuildAsyncTask(&task, TCron, "cron_test")
	at.CronInfo = CronInfo{
		CronId:       0,
		CronSchedule: "* * * * * *",
		CronState:    CronUP,
	}
	at.Reply = &pb.CommonCmdReply{
		ResultMsg: "结果",
		ResultTable: &pb.Table{
			Header: &pb.Row{
				Row: []string{"header1", "header2"},
			},
			Body: []*pb.Row{
				{Row: []string{"body11", "body12"},
				},
				{Row: []string{"body21", "body22"},
				},
			},
			Footer: &pb.Row{
				Row: []string{"footer1", "footer1"},
			},
		},
	}
	//cp := CronPersistence{
	//	CronInfo: at.CronInfo,
	//	Type:     at.Type,
	//	State:    at.State,
	//	Request:    *at.Request,
	//	Reply:      *at.Reply,
	//	FinishTime: at.FinishTime,
	//	NotifyMsg:  at.NotifyMsg,
	//}
	err := SaveCronToFile(at, "cron.save")
	if err != nil {
		fmt.Printf("%#v", err)
	}
	var yat *AsyncTask
	yat, err = LoadCronFile("cron.save")
	if err != nil {
		fmt.Printf("%#v", err)
	}
	fmt.Printf("\nUnmarshal is \n%#v\n", *yat)
}
