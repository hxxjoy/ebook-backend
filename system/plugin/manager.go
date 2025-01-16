// system/core/plugin/manager.go
package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
)

type PluginManager struct {
	app          *fiber.App
	plugins      map[string]Plugin
	events       *EventManager
	assetManager *AssetManager
}

// NewPluginManager 创建插件管理器
func NewPluginManager(app *fiber.App) *PluginManager {
	return &PluginManager{
		app:          app,
		plugins:      make(map[string]Plugin),
		events:       NewEventManager(),
		assetManager: NewAssetManager(),
	}
}

// LoadPlugin 加载插件
func (pm *PluginManager) LoadPlugin(name string) error {
	// 加载插件配置
	config, err := pm.loadPluginConfig(name)
	if err != nil {
		return err
	}

	// 加载后端插件
	backendPlugin, err := pm.loadBackendPlugin(config)
	if err != nil {
		return err
	}

	// 处理前端资源
	if err := pm.assetManager.ProcessPluginAssets(config); err != nil {
		return err
	}

	fmt.Println("---------- LoadPlugin ------------------")

	// 调用插件的 Init 方法 - 添加这段代码
	// 类型断言为 Plugin 接口
	if p, ok := backendPlugin.(Plugin); ok {
		fmt.Printf("Calling Init for plugin: %s\n", name)
		if err := p.Init(); err != nil {
			return fmt.Errorf("failed to initialize plugin %s: %v", name, err)
		} else {
			fmt.Println(" Init OK ")
		}
	} else {
		return fmt.Errorf("plugin %s does not implement Init method", name)
	}

	// 注册API路由
	backendPlugin.RegisterAPIRoutes(pm.app)

	// 注册数据模型
	pm.registerModels(backendPlugin.RegisterModels())

	// 存储插件实例
	pm.plugins[name] = backendPlugin

	return nil
}

// registerModels 注册数据模型
func (pm *PluginManager) registerModels(models []interface{}) error {
	// 这里实现数据模型注册逻辑
	// 例如使用 GORM 注册模型
	return nil
}

// GetPlugin 获取插件实例
func (pm *PluginManager) GetPlugin(name string) (Plugin, bool) {
	plugin, exists := pm.plugins[name]
	return plugin, exists
}

// EnablePlugin 启用插件
func (pm *PluginManager) EnablePlugin(name string) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return errors.New("plugin not found")
	}
	return plugin.Start()
}

// DisablePlugin 禁用插件
func (pm *PluginManager) DisablePlugin(name string) error {
	plugin, exists := pm.plugins[name]
	if !exists {
		return errors.New("plugin not found")
	}
	return plugin.Stop()
}

// LoadPlugins 加载所有插件
func (pm *PluginManager) LoadPlugins() error {
	// 获取插件目录下的所有插件
	plugins, err := ioutil.ReadDir("plugins")
	if err != nil {
		return err
	}

	// 遍历加载每个插件
	for _, p := range plugins {
		if p.IsDir() {
			if err := pm.LoadPlugin(p.Name()); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListPlugins 返回所有已加载的插件列表
func (pm *PluginManager) ListPlugins(c *fiber.Ctx) error {
	pluginList := make([]map[string]interface{}, 0)

	for name, plugin := range pm.plugins {
		pluginInfo := map[string]interface{}{
			"name":    name,
			"status":  plugin.Status(),  // 假设 Plugin 接口有 Status() 方法
			"version": plugin.Version(), // 假设 Plugin 接口有 Version() 方法
			// 其他需要的插件信息
		}
		pluginList = append(pluginList, pluginInfo)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   pluginList,
	})
}

// InstallPlugin 安装插件
func (pm *PluginManager) InstallPlugin(c *fiber.Ctx) error {
	// 获取要安装的插件信息
	var pluginInfo struct {
		Name string `json:"name"`
		// 其他需要的安装参数
	}

	if err := c.BodyParser(&pluginInfo); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// 安装插件的逻辑
	if err := pm.LoadPlugin(pluginInfo.Name); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Plugin installed successfully",
	})
}

// UninstallPlugin 卸载插件
func (pm *PluginManager) UninstallPlugin(c *fiber.Ctx) error {
	name := c.Params("name")

	plugin, exists := pm.plugins[name]
	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"status":  "error",
			"message": "Plugin not found",
		})
	}

	// 停止插件
	if err := plugin.Stop(); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// 删除插件
	delete(pm.plugins, name)

	// 这里可能还需要清理插件文件等资源

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Plugin uninstalled successfully",
	})
}
