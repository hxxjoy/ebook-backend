// system/core/plugin/interface.go
package plugin

import (
	"github.com/gofiber/fiber/v2"
)

// Plugin 插件接口
type Plugin interface {
	Init() error
	Start() error
	Stop() error
	Status() string
	Version() string
	RegisterAPIRoutes(app *fiber.App)
	RegisterModels() []interface{}
	Info() *PluginInfo
}

// PluginInfo 插件信息结构
type PluginInfo struct {
	Name         string            `json:"name"`
	Version      string            `json:"version"`
	Description  string            `json:"description,omitempty"`
	Author       string            `json:"author,omitempty"`
	License      string            `json:"license,omitempty"`
	Homepage     string            `json:"homepage,omitempty"`
	Repository   string            `json:"repository,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Dependencies map[string]string `json:"dependencies,omitempty"`
	Enabled      bool              `json:"enabled"`
}

// BasePlugin 基础插件实现
type BasePlugin struct {
	name    string
	version string
	config  *PluginConfig
	App     *fiber.App
	enabled bool
}

// Status implements Plugin.
func (p *BasePlugin) Status() string {
	panic("unimplemented")
}

// Version implements Plugin.
func (p *BasePlugin) Version() string {
	panic("unimplemented")
}

func (p *BasePlugin) Init() error {
	return nil
}

func (p *BasePlugin) Start() error {
	p.enabled = true
	return nil
}

func (p *BasePlugin) Stop() error {
	p.enabled = false
	return nil
}

func (p *BasePlugin) RegisterAPIRoutes(app *fiber.App) {
	// 基础实现为空
}

func (p *BasePlugin) RegisterModels() []interface{} {
	return nil
}

func (p *BasePlugin) Info() *PluginInfo {
	return &PluginInfo{
		Name:    p.name,
		Version: p.version,
		Enabled: p.enabled,
	}
}
