// system/core/plugin/asset_manager.go
package plugin

import (
	"errors"
	"os"
	"path/filepath"
)

// AssetManager 资源管理器
type AssetManager struct {
	publicDir string
	devMode   bool
}

// NewAssetManager 创建资源管理器
func NewAssetManager() *AssetManager {
	return &AssetManager{
		publicDir: "public/plugins",
		devMode:   os.Getenv("APP_ENV") == "development",
	}
}

// ProcessPluginAssets 处理插件资源
func (am *AssetManager) ProcessPluginAssets(config *PluginConfig) error {
	if config == nil {
		return errors.New("plugin config is nil")
	}

	// 创建插件资源目录
	pluginPublicDir := filepath.Join(am.publicDir, config.Name)
	if err := os.MkdirAll(pluginPublicDir, 0755); err != nil {
		return err
	}

	if am.devMode {
		return am.setupDevEnvironment(config)
	}
	return am.copyProductionAssets(config)
}

// 开发环境设置
func (am *AssetManager) setupDevEnvironment(config *PluginConfig) error {
	// 创建软链接到前端源码目录
	sourcePath := filepath.Join("plugins", config.Name, "frontend")
	targetPath := filepath.Join(am.publicDir, config.Name)

	// 删除已存在的链接
	os.Remove(targetPath)

	return os.Symlink(sourcePath, targetPath)
}

// 复制生产环境资源
func (am *AssetManager) copyProductionAssets(config *PluginConfig) error {
	sourcePath := filepath.Join("plugins", config.Name, "frontend/dist")
	targetPath := filepath.Join(am.publicDir, config.Name)

	return copyDir(sourcePath, targetPath)
}

// 复制目录及其内容
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 计算目标路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		// 复制文件
		return copyFile(path, targetPath)
	})
}

// 复制单个文件
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

// CleanPluginAssets 清理插件资源
func (am *AssetManager) CleanPluginAssets(pluginName string) error {
	return os.RemoveAll(filepath.Join(am.publicDir, pluginName))
}

// GetPluginAssetPath 获取插件资源路径
func (am *AssetManager) GetPluginAssetPath(pluginName string) string {
	return filepath.Join(am.publicDir, pluginName)
}
