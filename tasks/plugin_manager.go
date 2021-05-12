package tasks

import (
	"path/filepath"
	"plugin"
	"strings"
)

type PluginState int32

const (
	Unavailable PluginState = 0
	Available   PluginState = 1
	Replaced    PluginState = 2
	Removed     PluginState = 3
)

type PlugInfo struct {
	Name   string
	Path   string
	State  PluginState
	Plugin *plugin.Plugin
}

type PluginManager struct {
	PluginDir string
	PlugInfos map[string]PlugInfo
}

func (pm *PluginManager) LoadPlugin(pluginName string) error {
	if pi, ok := pm.PlugInfos[pluginName]; ok {
		if pi.State == Available {
			log.Debugln("插件已存在")
			return nil
		}
	}
	pluginPath := filepath.Join(pm.PluginDir, strings.ToLower(pluginName)+".so")
	log.Debugln("加载插件", pluginPath, "中...")
	p, err := plugin.Open(pluginPath)
	if err != nil {
		return err
	}
	pluginInfo := PlugInfo{Name: pluginName, Path: pluginPath, State: Available, Plugin: p}
	pm.PlugInfos[pluginName] = pluginInfo
	log.Debugln("插件已被插入管理器")
	return nil
}
