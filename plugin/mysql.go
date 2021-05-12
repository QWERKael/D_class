package main

import (
	"atk_D_class/logger"
	"atk_D_class/pb"
	"atk_D_class/prompt"
	"atk_D_class/tasks"
	"atk_D_class/utils"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

var log = logger.Logger

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
desc: "MySQL Plugin\n提供MySQL相关操作的插件\n
mysql state \n\t查看当前连接的相关信息\n
mysql connect \n\t连接到数据库, 默认连接到服务端本地的3306\n
mysql disconnect \n\t断开当前连接\n
mysql Q [SQL] \n\t查询 SQL 命令\n
mysql show \n\t展示 MySQL 的相关统计信息\n
mysql show host conn\n\t统计各服务器的连接数\n
mysql show max mem\n\t计算最大内存使用量\n
mysql show read rate\n\计算Buffer Pool的读磁盘命中率\n
mysql show table usage [Table Name]\n\t查看指定库中的表使用情况\n"
yess:
  mysql:
    type: "Plugin"
    desc: "提供MySQL相关操作"
    yess:
      state:
        type: "Cmd"
        desc: "查看当前连接的相关信息"
      connect:
        type: "Cmd"
        desc: "连接数据库"
        yess:
          "--host=":
            type: "ArgKey"
            desc: "连接到指定数据库IP"
          "--port=":
            type: "ArgKey"
            desc: "连接到指定数据库端口"
          "--user=":
            type: "ArgKey"
            desc: "使用指定用户连接"
          "--password=":
            type: "ArgKey"
            desc: "指定用户的密码"
          "--network=":
            type: "ArgKey"
            desc: "使用指定网络协议连接"
          "--database=":
            type: "ArgKey"
            desc: "连接到指定数据库"
      disconnect:
        type: "Cmd"
        desc: "断开当前连接"
      q:
        type: "Cmd"
        desc: "查询SQL命令"
      show:
        type: "Cmd"
        desc: "查看数据库相关信息"
        yess:
          "max mem":
            type: "SubCmd"
            desc: "计算最大内存使用量"
          "read rate":
            type: "SubCmd"
            desc: "计算Buffer Pool的读磁盘命中率"
          "host conn":
            type: "SubCmd"
            desc: "统计各服务器的连接数"
          "table usage ":
            type: "SubCmd"
            desc: "查看指定库中的表使用情况"
      version:
        type: "Cmd"
        desc: "查看版本信息"
      help:
        type: "Cmd"
        desc: "查看帮助信息, 并加载智能提示"`
	return helpMsg, nil
}

//数据库连接信息
const (
	USERNAME = "admin"
	PASSWORD = "123456"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "test"
)

var MysqlConnInfo = MySQLConnInfo{}
var MysqlConn = MySQLConn{}

type MySQLConnInfo struct {
	UserName string
	Password string
	Network  string
	Host     string
	Port     int
	Database string
}

func InitConn(task tasks.Task) error {
	var err error
	mci := MySQLConnInfo{}
	mci.Host = utils.StringDefault(task.Args["host"], "127.0.0.1")
	mci.Port = utils.StringDefaultInt(task.Args["port"], 3306)
	mci.UserName = utils.StringDefault(task.Args["user"], "admin")
	mci.Password = utils.StringDefault(task.Args["password"], "123456")
	mci.Network = utils.StringDefault(task.Args["network"], "tcp")
	mci.Database = utils.StringDefault(task.Args["database"], "")
	ms := mci.String()
	if MysqlConn.Db == nil {
		MysqlConnInfo = mci
		MysqlConn.ConnStr = ms
		err = MysqlConn.Connect()
		if err != nil {
			return err
		}
		return nil
	} else if MysqlConn.ConnStr == ms {
		return nil
	} else if utils.IsInFlag("reconnect", task.Flags) {
		MysqlConnInfo = mci
		MysqlConn.ConnStr = ms
		err = MysqlConn.Disconnect()
		if err != nil {
			return err
		}
		err = MysqlConn.Connect()
		if err != nil {
			return err
		}
		return nil
	} else {
		log.Debugf("已有链接: %s\n请求连接: %s", MysqlConn.ConnStr, ms)
		return errors.New("连接已存在, 重连请使用 -reconnect")
	}
}

func (i *MySQLConnInfo) String() string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s", i.UserName, i.Password, i.Network, i.Host, i.Port, i.Database)
}

type MySQLConn struct {
	ConnStr string
	Db      *sql.DB
}

func (m *MySQLConn) Connect() error {
	var err error
	m.Db, err = sql.Open("mysql", m.ConnStr)
	log.Debugf("建立连接字符串: %s", m.ConnStr)
	return err
}

func (m *MySQLConn) Disconnect() error {
	err := m.Db.Close()
	return err
}

func (m *MySQLConn) Query(sql string) ([]string, [][]string) {
	var (
		header []string
		body   [][]string
	)

	rows, err := m.Db.Query(sql)
	if err != nil {
		log.Debugf("Failed to run query", err)
		utils.CheckErrorPanic(err)
	}

	cols, err := rows.Columns()
	if err != nil {
		log.Debugf("Failed to get columns", err)
		utils.CheckErrorPanic(err)
	}
	header = cols

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		result := make([]string, len(cols))
		err = rows.Scan(dest...)
		if err != nil {
			log.Debugf("Failed to scan row", err)
			utils.CheckErrorPanic(err)
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "Null"
			} else {
				result[i] = string(raw)
			}
		}

		//fmt.Printf("%#v\n", result)
		body = append(body, result)
	}
	return header, body
}

func State(task tasks.Task) (*pb.CommonCmdReply, error) {
	body := [][]string{
		[]string{"UserName", MysqlConnInfo.UserName},
		[]string{"Password", strings.Repeat("*", len(MysqlConnInfo.Password))},
		[]string{"Network", MysqlConnInfo.Network},
		[]string{"Host", MysqlConnInfo.Host},
		[]string{"Port", strconv.Itoa(MysqlConnInfo.Port)},
		[]string{"Database", MysqlConnInfo.Database},
	}

	reply := prompt.ToTable([]string{}, []string{}, body, 0)
	return reply, nil
}

func Connect(task tasks.Task) (*pb.CommonCmdReply, error) {
	err := InitConn(task)
	if err != nil {
		return nil, errors.New("连接出错: " + err.Error())
	}

	reply, _ := State(task)
	reply.ResultMsg = "连接已建立成功"
	return reply, nil
}

func Disconnect(task tasks.Task) (string, error) {
	if MysqlConn.Db == nil {
		return "连接未建立", nil
	} else {
		err := MysqlConn.Disconnect()
		if err != nil {
			return "", err
		}
		return "连接已断开", nil
	}
}

func query(task tasks.Task) (*pb.CommonCmdReply, error) {
	if MysqlConn.Db == nil {
		err := InitConn(task)
		if err != nil {
			return nil, err
		}
	}
	sqlStr := task.SubCmd[0]
	header, body := MysqlConn.Query(sqlStr)

	footer := make([]string, len(header))
	if len(header) > 1 {
		footer[0] = "总计"
		footer[1] = fmt.Sprintf("%d 行", len(body))
	} else {
		footer[0] = fmt.Sprintf("总计: %d 行", len(body))
	}
	reply := prompt.ToTable(header, footer, body, 0)
	return reply, nil
}

func Q(task tasks.Task) (*pb.CommonCmdReply, error) {
	return query(task)
}

func Show(task tasks.Task) (*pb.CommonCmdReply, error) {
	if MysqlConn.Db == nil {
		err := InitConn(task)
		if err != nil {
			return nil, err
		}
	}
	var ccr *pb.CommonCmdReply
	switch strings.ToLower(strings.Join(task.SubCmd, " ")) {
	case "max mem":
		return showMaxMemUsage(task.Args)
	case "read rate":
		return showBufferPoolReadRate()
	case "host conn":
		return showHostConn()
	default:
		if strings.ToLower(strings.Join(task.SubCmd[:2], " ")) == "table usage" {
			return showTableUsage(task.SubCmd[2])
		}
	}
	return ccr, nil
}

func showHostConn() (*pb.CommonCmdReply, error) {
	log.Infof("统计各服务器的连接情况")
	var header []string
	var body [][]string
	header, body = MysqlConn.Query("SELECT `USER`,LEFT(`host`,LOCATE(':',`host`)-1) as HOST,COUNT(*) as COUNT, SUM(IF(`state`<>'Sleep',1,0)) AS 'Count Without Sleep' FROM information_schema.`PROCESSLIST` GROUP BY `USER`,LEFT(`host`,LOCATE(':',`host`)-1) ORDER BY count(*) desc")
	totalCount := 0
	noSleepCount := 0
	for _, row := range body {
		totalCount += utils.StringDefaultInt(row[2], 0)
		noSleepCount += utils.StringDefaultInt(row[3], 0)
	}
	footer := make([]string, len(header))
	footer[0] = "Total Connections(No Sleep Connections)"
	footer[1] = fmt.Sprintf("%d (%d)", totalCount, noSleepCount)
	reply := prompt.ToTable(header, footer, body, 0)
	return reply, nil
}

func showTableUsage(dbName string) (*pb.CommonCmdReply, error) {
	log.Infof("查看库 %s 中的表使用情况", dbName)
	var header []string
	var body [][]string
	header, body = MysqlConn.Query(fmt.Sprintf("select table_name,(data_length + index_length) as table_mb,table_rows from information_schema.tables where table_schema='%s'", dbName))
	tbMb := 0.0
	tbRows := 0
	for _, row := range body {
		tbMb += utils.StringDefaultFloat64(row[1], 0)
		tbRows += utils.StringDefaultInt(row[2], 0)
	}
	footer := make([]string, len(header))
	footer[0] = "Total Size(Total Rows)"
	footer[1] = fmt.Sprintf("%f MB (%d Rows)", tbMb/1024/1024, tbRows)
	reply := prompt.ToTable(header, footer, body, 0)
	return reply, nil
}

func getVariablesLike(key string, isGlobal bool) ([]string, [][]string, map[string]string) {
	var header []string
	var body [][]string
	var variables = make(map[string]string)
	if isGlobal {
		header, body = MysqlConn.Query("SHOW GLOBAL VARIABLES LIKE '%" + key + "%'")
	} else {
		header, body = MysqlConn.Query("SHOW VARIABLES LIKE '%" + key + "%'")
	}
	for _, raw := range body {
		variables[raw[0]] = raw[1]
	}
	return header, body, variables
}

func getStatusLike(key string, isGlobal bool) ([]string, [][]string, map[string]string) {
	var header []string
	var body [][]string
	var status = make(map[string]string)
	if isGlobal {
		header, body = MysqlConn.Query("SHOW GLOBAL STATUS LIKE '%" + key + "%'")
	} else {
		header, body = MysqlConn.Query("SHOW STATUS LIKE '%" + key + "%'")
	}
	for _, raw := range body {
		status[raw[0]] = raw[1]
	}
	return header, body, status
}

// 计算最大内存使用量
func showMaxMemUsage(args map[string]string) (*pb.CommonCmdReply, error) {
	log.Infof("计算最大内存使用量")
	//sqlStr := task.SubCmd[0]
	header, _, variables := getVariablesLike("", true)

	// 通过args估算修改参数后的最大内存使用量
	for key, val := range args {
		variables[key] = val
	}

	maxMem := utils.StringDefaultFloat64(variables["key_buffer_size"], 0) +
		utils.StringDefaultFloat64(variables["query_cache_size"], 0) +
		utils.StringDefaultFloat64(variables["tmp_table_size"], 0) +
		utils.StringDefaultFloat64(variables["innodb_buffer_pool_size"], 0) +
		utils.StringDefaultFloat64(variables["innodb_additional_mem_pool_size"], 0) +
		utils.StringDefaultFloat64(variables["innodb_log_buffer_size"], 0) +
		utils.StringDefaultFloat64(variables["max_connections"], 0)*(utils.StringDefaultFloat64(variables["read_buffer_size"], 0)+
			utils.StringDefaultFloat64(variables["read_rnd_buffer_size"], 0)+
			utils.StringDefaultFloat64(variables["sort_buffer_size"], 0)+
			utils.StringDefaultFloat64(variables["join_buffer_size"], 0)+
			utils.StringDefaultFloat64(variables["binlog_cache_size"], 0)+
			utils.StringDefaultFloat64(variables["thread_stack"], 0))

	body := [][]string{
		{"key_buffer_size", variables["key_buffer_size"]},
		{"query_cache_size", variables["query_cache_size"]},
		{"tmp_table_size", variables["tmp_table_size"]},
		{"innodb_buffer_pool_size", variables["innodb_buffer_pool_size"]},
		{"innodb_additional_mem_pool_size", variables["innodb_additional_mem_pool_size"]},
		{"innodb_log_buffer_size", variables["innodb_log_buffer_size"]},
		{"max_connections", variables["max_connections"]},
		{"read_buffer_size", variables["read_buffer_size"]},
		{"read_rnd_buffer_size", variables["read_rnd_buffer_size"]},
		{"sort_buffer_size", variables["sort_buffer_size"]},
		{"join_buffer_size", variables["join_buffer_size"]},
		{"binlog_cache_size", variables["binlog_cache_size"]},
		{"thread_stack", variables["thread_stack"]},
	}

	footer := make([]string, len(header))
	footer[0] = "Max Mem Usage"
	footer[1] = fmt.Sprint(maxMem/1024/1024/1024, " GB")
	reply := prompt.ToTable(header, footer, body, 0)
	return reply, nil
}

// 计算Buffer Pool的读磁盘命中率
func showBufferPoolReadRate() (*pb.CommonCmdReply, error) {
	log.Infof("计算Buffer Pool的读磁盘命中率")
	//sqlStr := task.SubCmd[0]
	header, body, status := getStatusLike("Innodb_buffer_pool_read_", true)
	rate := (utils.StringDefaultFloat64(status["Innodb_buffer_pool_read_requests"], 0) - utils.StringDefaultFloat64(status["Innodb_buffer_pool_reads"], 0)) / utils.StringDefaultFloat64(status["Innodb_buffer_pool_read_requests"], 0)

	footer := make([]string, len(header))
	footer[0] = "Buffer Pool Read Rate"
	footer[1] = fmt.Sprint(rate*100, "%")
	reply := prompt.ToTable(header, footer, body, 0)
	return reply, nil
}

func main() {
	conn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	db, err := sql.Open("mysql", conn)
	defer db.Close()

	if err != nil {
		fmt.Println("Failed to connect", err)
		return
	}

	rows, err := db.Query(`SELECT * FROM information_schema.PROCESSLIST WHERE COMMAND <> 'Sleep' ORDER BY TIME DESC;`)
	if err != nil {
		fmt.Println("Failed to run query", err)
		return
	}

	cols, err := rows.Columns()
	if err != nil {
		fmt.Println("Failed to get columns", err)
		return
	}
	fmt.Printf("%#v\n", cols)

	// Result is your slice string.
	rawResult := make([][]byte, len(cols))
	result := make([]string, len(cols))

	dest := make([]interface{}, len(cols)) // A temporary interface{} slice
	for i, _ := range rawResult {
		dest[i] = &rawResult[i] // Put pointers to each string in the interface slice
	}

	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			fmt.Println("Failed to scan row", err)
			return
		}

		for i, raw := range rawResult {
			if raw == nil {
				result[i] = "\\N"
			} else {
				result[i] = string(raw)
			}
		}

		fmt.Printf("%#v\n", result)
	}
}
