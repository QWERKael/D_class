package prompt

import (
	"atk_D_class/utils"
	"github.com/c-bata/go-prompt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sort"
	"strings"
)

var es = ExtraSuggest{
	Type:    Root,
	Suggest: prompt.Suggest{Text: "root", Description: "根"},
	ESs: map[string]*ExtraSuggest{
		"cmd": {
			Type:    Plugin,
			Suggest: prompt.Suggest{Text: "cmd", Description: "执行command命令"},
			ESs: map[string]*ExtraSuggest{
				"help": {
					Type:    Cmd,
					Suggest: prompt.Suggest{Text: "help", Description: "查看帮助信息, 并加载智能提示"},
					ESs:     nil,
				},
			},
		},
		"show": {
			Type:    Plugin,
			Suggest: prompt.Suggest{Text: "show", Description: "展示服务器相关信息"},
			ESs: map[string]*ExtraSuggest{
				"help": {
					Type:    Cmd,
					Suggest: prompt.Suggest{Text: "help", Description: "查看帮助信息, 并加载智能提示"},
					ESs:     nil,
				},
			},
		},
		"async": {
			Type:    Plugin,
			Suggest: prompt.Suggest{Text: "async", Description: "执行异步任务"},
			ESs: map[string]*ExtraSuggest{
				"help": {
					Type:    Cmd,
					Suggest: prompt.Suggest{Text: "help", Description: "查看帮助信息, 并加载智能提示"},
					ESs:     nil,
				},
			},
		},
		"email": {
			Type:    Plugin,
			Suggest: prompt.Suggest{Text: "email", Description: "发送电子邮件"},
			ESs: map[string]*ExtraSuggest{
				"help": {
					Type:    Cmd,
					Suggest: prompt.Suggest{Text: "help", Description: "查看帮助信息, 并加载智能提示"},
					ESs:     nil,
				},
			},
		},
		"bee": {
			Type:    Plugin,
			Suggest: prompt.Suggest{Text: "bee", Description: "Honeycomb的客户端"},
			ESs: map[string]*ExtraSuggest{
				"help": {
					Type:    Cmd,
					Suggest: prompt.Suggest{Text: "help", Description: "查看帮助信息, 并加载智能提示"},
					ESs:     nil,
				},
			},
		},
		"upload": {
			Type:    Cmd,
			Suggest: prompt.Suggest{Text: "upload", Description: "上传文件"},
			ESs: map[string]*ExtraSuggest{
				"-plugin": {
					Type:    Flag,
					Suggest: prompt.Suggest{Text: "-plugin", Description: "上传到插件目录"},
					ESs:     nil,
				},
				"-script": {
					Type:    Flag,
					Suggest: prompt.Suggest{Text: "-script", Description: "上传到脚本目录"},
					ESs:     nil,
				},
				"-config": {
					Type:    Flag,
					Suggest: prompt.Suggest{Text: "-config", Description: "上传配置文件"},
					ESs:     nil,
				},
				"-update": {
					Type:    Flag,
					Suggest: prompt.Suggest{Text: "-update", Description: "更新主程序"},
					ESs:     nil,
				},
			},
		},
		"loc": {
			Type:    Prefix,
			Suggest: prompt.Suggest{Text: "loc", Description: "调用本地插件"},
			ESs:     nil,
		},
		"change-to": {
			Type:    Default,
			Suggest: prompt.Suggest{Text: "change-to", Description: "改变连接到新的服务端"},
			ESs:     nil,
		},
		"restart": {
			Type:    Default,
			Suggest: prompt.Suggest{Text: "restart", Description: "重启服务端"},
			ESs:     nil,
		},
		"exit": {
			Type:    Default,
			Suggest: prompt.Suggest{Text: "exit", Description: "退出"},
			ESs:     nil,
		},
	},
}

type ExtraSuggestType int32

const (
	Default  ExtraSuggestType = 0
	Root     ExtraSuggestType = 1
	Prefix   ExtraSuggestType = 2
	Plugin   ExtraSuggestType = 3
	Cmd      ExtraSuggestType = 4
	SubCmd   ExtraSuggestType = 5
	Flag     ExtraSuggestType = 6
	ArgKey   ExtraSuggestType = 7
	ArgValue ExtraSuggestType = 8
)

func (est *ExtraSuggestType) ToString() string {
	ests := []string{
		"Default",
		"Root",
		"Prefix",
		"Plugin",
		"Cmd",
		"SubCmd",
		"Flag",
		"ArgKey",
		"ArgValue",
	}
	return ests[*est]
}

func StrToEST(s string) ExtraSuggestType {
	s = strings.ToLower(s)
	ests := map[string]ExtraSuggestType{
		"default":  Default,
		"root":     Root,
		"prefix":   Prefix,
		"plugin":   Plugin,
		"cmd":      Cmd,
		"subcmd":   SubCmd,
		"flag":     Flag,
		"argkey":   ArgKey,
		"argvalue": ArgValue,
	}
	return ests[s]
}

type ExtraSuggest struct {
	Type    ExtraSuggestType
	Suggest prompt.Suggest
	ESs     map[string]*ExtraSuggest
}

func (es *ExtraSuggest) ToSuggest() []prompt.Suggest {
	suggest := make([]prompt.Suggest, 0)
	for _, e := range es.ESs {
		suggest = append(suggest, e.Suggest)
	}
	return suggest
}

func (es *ExtraSuggest) Add(text string, desc string, t ExtraSuggestType) error {
	if _, ok := es.ESs[text]; ok {
		return errors.New("提示符节点 " + text + " 已存在")
	}
	if es.ESs == nil {
		es.ESs = make(map[string]*ExtraSuggest, 0)
	}
	es.ESs[text] = &ExtraSuggest{
		Type:    t,
		Suggest: prompt.Suggest{Text: text, Description: desc},
		ESs:     nil,
	}
	return nil
}

type YetExtraSuggest struct {
	Type string                     `yaml:"type"`
	Text string                     `yaml:"text"`
	Desc string                     `yaml:"desc"`
	ESs  map[string]YetExtraSuggest `yaml:"yess"`
}

func (yes *YetExtraSuggest) ToES() *ExtraSuggest {
	es := &ExtraSuggest{
		Type: StrToEST(yes.Type),
		Suggest: prompt.Suggest{
			Text:        yes.Text,
			Description: yes.Desc,
		},
		ESs: nil,
	}
	if yes.ESs != nil {
		es.ESs = make(map[string]*ExtraSuggest)
		for yText, y := range yes.ESs {
			if y.Text == "" {
				y.Text = yText
			}
			es.ESs[yText] = y.ToES()
		}
	}
	return es
}

// 从配置文件加载 智能提示
func LoadPrompt(path string) {
	if ft, err := utils.CheckPath(path); err == nil {
		if ft == utils.File {
			b, err := ioutil.ReadFile(path)
			utils.CheckErrorPanic(err)
			//result := make(map[string]string)
			resEs := ParseYAML(b)
			for k, v := range resEs.ESs {
				es.ESs[k] = v
			}

			//promptSet := result["prompt"]
			//for k, v := range result {
			//	log.Debugf(`加载快捷命令 "%s": "%s"`, k, v)
			//	err := es.Add(k, v, Default)
			//	if err != nil {
			//		log.Debugf("加载快捷命令失败: %s", err.Error())
			//	}
			//	//suggestion = append(suggestion, prompt.Suggest{Text: k, Description: v})
			//}
		}
	} else {
		utils.CheckErrorPanic(err)
	}
}

// 将 []byte 解析成 *ExtraSuggest
func ParseYAML(b []byte) *ExtraSuggest {
	var result YetExtraSuggest
	err := yaml.Unmarshal(b, &result)
	utils.CheckErrorPanic(err)
	resEs := result.ToES()
	return resEs
}

func completer(in prompt.Document) []prompt.Suggest {
	args := strings.Split(strings.TrimLeft(in.TextBeforeCursor(), " "), " ")
	if in.TextBeforeCursor() == "" || len(args) < 2 {
		suggestion := es.ToSuggest()
		sortedSug := sortedSuggests{suggestion}
		sortedSug.Sort()
		suggestion = sortedSug.suggests
		return prompt.FilterHasPrefix(suggestion, in.GetWordBeforeCursorWithSpace(), true)
	}

	args = args[:len(args)-1]

	suggestion := make([]prompt.Suggest, 0)
	ess := &ExtraSuggest{}
	*ess = es
	e := &ExtraSuggest{}
	cmdSug := ExtraSuggest{}
	subCmdSug := ExtraSuggest{}
	var ok bool
	for i, arg := range args {
		//fmt.Printf("\n%#v\n", arg)
		if e, ok = ess.ESs[arg]; ok {
			if e.Type != Flag && e.Type != ArgKey && e.Type != ArgValue {
				ess = ess.ESs[arg]
				if e.Type == Cmd {
					cmdSug = *ess
				}
				if e.Type == SubCmd {
					subCmdSug = *ess
				}
				suggestion = ess.ToSuggest()

				if e.Type == Prefix {
					suggestion = append(suggestion, es.ToSuggest()...)
				}

				//fmt.Println()
				//fmt.Println(ess.Type.ToString())
				continue
			} else if e.Type == ArgKey {
				suggestion = ess.ESs[arg].ToSuggest()
				//fmt.Println(ess.ESs[arg].Type.ToString())
				if i+1 == len(args) {
					break
				}
			}
		} else {
			//fmt.Printf("\n没有找到 %s 在 %#v 中\n", arg, ess.Suggest)
		}
		suggestion = append(cmdSug.ToSuggest(), subCmdSug.ToSuggest()...)
		//fmt.Printf("SFA--%s", ess.Type.ToString())
	}

	sortedSug := sortedSuggests{suggestion}
	sortedSug.Sort()
	suggestion = sortedSug.suggests
	//fmt.Printf("\nsuggest 是: %#v\n", suggestion)
	return prompt.FilterContains(suggestion, in.GetWordBeforeCursor(), true)
}

// 定义 []prompt.Suggest 的排序
type sortedSuggests struct {
	suggests []prompt.Suggest
}

func (ss sortedSuggests) Len() int {
	return len(ss.suggests)
}
func (ss sortedSuggests) Less(i, j int) bool {
	return ss.suggests[i].Text < ss.suggests[j].Text
}
func (ss sortedSuggests) Swap(i, j int) {
	ss.suggests[i], ss.suggests[j] = ss.suggests[j], ss.suggests[i]
}
func (ss sortedSuggests) Sort() {
	sort.Sort(ss)
}
