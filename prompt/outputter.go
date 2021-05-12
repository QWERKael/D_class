package prompt

import (
	"atk_D_class/pb"
	"atk_D_class/utils"
	prettyTable "github.com/jedib0t/go-pretty/table"
	"os"
)

func CommonCmdOutputter(reply *pb.CommonCmdReply, isReverse bool) error {
	var err error
	if reply.ResultMsg != "" {
		println(reply.ResultMsg)
	}
	if reply.ResultTable != nil {
		t := reply.ResultTable
		if isReverse {
			err = printTableReversal(t)
		} else {
			err = printTable(t)
		}
		utils.CheckErrorPanic(err)
	}
	return nil
}

func printTable(table *pb.Table) error {
	header := fulfillRow(table.Header.Row)
	footer := fulfillRow(table.Footer.Row)
	body := fulfillTable(table.Body)
	t := prettyTable.NewWriter()
	t.SetOutputMirror(os.Stdout)
	if len(header) > 0 {
		t.AppendHeader(header)
	}
	t.AppendRows(body)
	if len(footer) > 0 {
		t.AppendFooter(footer)
	}
	t.SetStyle(prettyTable.StyleColoredBright)
	//t.SetStyle(prettyTable.StyleColoredCyanWhiteOnBlack)
	//t.SetColumnConfigs([]prettyTable.ColumnConfig{
	//	{
	//		WidthMin: 1,
	//		WidthMax: 1,
	//	},
	//})
	//t.SetAllowedRowLength(50)
	//t.SetStyle(prettyTable.StyleDouble)
	t.Render()
	return nil
}

func printTableReversal(table *pb.Table) error {
	header := prettyTable.Row{"Row ID", "Column Name", "Content"}
	footer := prettyTable.Row{}
	if len(table.Footer.Row) >= 2 {
		footer = fulfillRow(table.Footer.Row)[:2]
	} else if len(table.Footer.Row) == 1 {
		footer = fulfillRow(table.Footer.Row)[:1]
	}
	var body []prettyTable.Row
	for i, row := range table.Body {
		rowId := i + 1
		for j, colName := range table.Header.Row {
			r := prettyTable.Row{rowId, colName, row.Row[j]}
			body = append(body, r)
		}
	}
	t := prettyTable.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRows(body)
	t.AppendFooter(footer)
	t.SetStyle(prettyTable.StyleColoredBright)
	//t.SetStyle(prettyTable.StyleColoredCyanWhiteOnBlack)
	//t.SetColumnConfigs([]prettyTable.ColumnConfig{
	//	{
	//		WidthMin: 1,
	//		WidthMax: 1,
	//	},
	//})
	//t.SetAllowedRowLength(50)
	//t.SetStyle(prettyTable.StyleDouble)
	t.Render()
	return nil
}

func fulfillRow(s []string) prettyTable.Row {
	row := make([]interface{}, len(s))
	for i, cell := range s {
		row[i] = cell
	}
	return row
}

func fulfillTable(s []*pb.Row) []prettyTable.Row {
	table := make([]prettyTable.Row, len(s))
	for i, row := range s {
		table[i] = fulfillRow(row.Row)
	}
	return table
}

//func CommonCmdOutputter(reply *pb.CommonCmdReply) error {
//	if reply.ResultMsg != "" {
//		println(reply.ResultMsg)
//	}
//	if reply.ResultTable != nil {
//		t := reply.ResultTable
//		header := t.Header.Row
//		footer := t.Footer.Row
//		rows := t.Body
//		body := make([][]string, 0)
//		for _, row := range rows {
//			body = append(body, row.Row)
//		}
//		table := tablewriter.NewWriter(os.Stdout)
//		table.SetHeader(header)
//		log.Debugf("header已加载\n%#v", header)
//		table.SetFooter(footer)
//		log.Debugf("footer已加载\n%#v", footer)
//		table.SetBorder(true)
//		table.SetRowLine(true)
//		table.SetAutoMergeCells(true)
//
//		//table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold, tablewriter.BgGreenColor},
//		//	tablewriter.Colors{tablewriter.FgHiRedColor, tablewriter.Bold, tablewriter.BgBlackColor},
//		//	tablewriter.Colors{tablewriter.BgRedColor, tablewriter.FgWhiteColor},
//		//	tablewriter.Colors{tablewriter.BgCyanColor, tablewriter.FgWhiteColor})
//		//log.Debugf("header颜色已加载")
//		//table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
//		//	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
//		//	tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
//		//	tablewriter.Colors{tablewriter.Bold, tablewriter.FgBlackColor})
//		//log.Debugf("column颜色已加载")
//		//table.SetFooterColor(tablewriter.Colors{}, tablewriter.Colors{},
//		//	tablewriter.Colors{tablewriter.Bold},
//		//	tablewriter.Colors{tablewriter.FgHiRedColor})
//		log.Debugf("footer颜色已加载")
//		log.Debugf("颜色已全部加载")
//
//		table.Rich()
//		table.AppendBulk(body)
//		log.Debugf("数据已加载\n%#v", body)
//		table.SetReflowDuringAutoWrap(true)
//		table.Render()
//		log.Debugf("表格已输出")
//	}
//	return nil
//}

func ToTable(header []string, footer []string, body [][]string, limit int) *pb.CommonCmdReply {
	if limit == 0 {
		limit = len(body)
	}
	reply := &pb.CommonCmdReply{
		ResultTable: &pb.Table{
			Header: &pb.Row{},
			Footer: &pb.Row{},
			Body:   make([]*pb.Row, limit),
		},
		Status: pb.CommonCmdReply_Ok,
	}
	reply.ResultTable.Header.Row = header
	reply.ResultTable.Footer.Row = footer

	for i, _ := range reply.ResultTable.Body {
		reply.ResultTable.Body[i] = &pb.Row{Row: body[i]}
		log.Debugf("body[%d]: %#v", i, body[i])
	}
	log.Debugf("reply.ResultTable.Body: %#v", reply.ResultTable.Body)
	return reply
}
