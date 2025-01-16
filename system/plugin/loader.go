// system/core/plugin/loader.go
package plugin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// PluginConfig 插件配置结构
type PluginConfig struct {
	Name         string `json:"name"`
	Version      string `json:"version"`
	Description  string `json:"description"`
	Author       string `json:"author"`
	License      string `json:"license"`
	Enable       string `json:"enable"`
	Dependencies struct {
		Plugins  []string          `json:"plugins"`
		Frontend map[string]string `json:"frontend"`
		Backend  map[string]string `json:"backend"`
	} `json:"dependencies"`
	Frontend struct {
		Entry  string   `json:"entry"`
		Assets []string `json:"assets"`
	} `json:"frontend"`
	Backend struct {
		Main   string   `json:"main"`
		Models []string `json:"models"`
	} `json:"backend"`
}

// loadPluginConfig 加载插件配置
func (pm *PluginManager) loadPluginConfig(name string) (*PluginConfig, error) {
	// 构建插件配置文件路径
	configPath := filepath.Join("plugins", name, "plugin.json")

	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, errors.New("failed to read plugin config: " + err.Error())
	}

	// 解析配置文件
	var config PluginConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, errors.New("failed to parse plugin config: " + err.Error())
	}

	fmt.Println(config, "=====================")
	// 验证必要字段
	if config.Name == "" {
		return nil, errors.New("plugin name is required")
	}

	return &config, nil
}

// loadBackendPlugin 加载后端插件
func (pm *PluginManager) loadBackendPlugin(config *PluginConfig) (Plugin, error) {
	// 检查插件是否已加载
	if _, exists := pm.plugins[config.Name]; exists {
		return nil, errors.New("plugin already loaded")
	}

	// 检查依赖
	if err := pm.checkDependencies(config); err != nil {
		return nil, err
	}

	// 创建插件实例
	plugin := &BasePlugin{
		name:    config.Name,
		version: config.Version,
		config:  config,
		App:     pm.app,
	}

	// 初始化插件
	if err := plugin.Init(); err != nil {
		return nil, errors.New("failed to initialize plugin: " + err.Error())
	}

	return plugin, nil
}

// checkDependencies 检查插件依赖
func (pm *PluginManager) checkDependencies(config *PluginConfig) error {
	for _, dep := range config.Dependencies.Plugins {
		if _, exists := pm.plugins[dep]; !exists {
			return errors.New("required plugin not loaded: " + dep)
		}
	}
	return nil
}
