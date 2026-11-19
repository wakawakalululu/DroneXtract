package deps

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ANG13T/DroneXtract/assets"
)

// AppDir 返回程序所在目录，所有配置与输出都以它为基准。
func AppDir() string {
	return assets.Dir()
}

// ConfigFile 返回当前生效的配置文件路径（程序目录下的 .env）。
func ConfigFile() string {
	return filepath.Join(AppDir(), ".env")
}

// ConfigPaths 返回配置文件的候选路径，按优先级从高到低：
//
//  1. $DRONEXTRACT_ENV            显式指定的文件
//  2. <程序目录>/.env              主配置（默认写入位置）
//  3. <工作目录>/.env              工作目录覆盖（便于一次会话临时改动）
//
// 内嵌默认值作为最终兜底，保证任何情况下都有可用参数。
func ConfigPaths() []string {
	var paths []string

	if p := strings.TrimSpace(os.Getenv("DRONEXTRACT_ENV")); p != "" {
		paths = append(paths, p)
	}

	appDir := AppDir()
	paths = append(paths, filepath.Join(appDir, ".env"))

	if wd, err := os.Getwd(); err == nil {
		wdEnv := filepath.Join(wd, ".env")
		if wdEnv != paths[len(paths)-1] {
			paths = append(paths, wdEnv)
		}
	}

	// 兼容早期版本写到用户配置目录的情况，避免升级后配置丢失
	for _, legacy := range legacyConfigPaths() {
		paths = append(paths, legacy)
	}

	return paths
}

// legacyConfigPaths 返回历史版本的配置位置，仅用于只读兼容。
func legacyConfigPaths() []string {
	var out []string
	home, _ := os.UserHomeDir()
	if cfg := strings.TrimSpace(os.Getenv("XDG_CONFIG_HOME")); cfg != "" {
		out = append(out, filepath.Join(cfg, "dronextract", ".env"))
	} else if home != "" {
		out = append(out, filepath.Join(home, ".config", "dronextract", ".env"))
	}
	if appdata := strings.TrimSpace(os.Getenv("APPDATA")); appdata != "" {
		out = append(out, filepath.Join(appdata, "dronextract", ".env"))
	}
	if home != "" {
		out = append(out, filepath.Join(home, ".dronextract", ".env"))
	}
	return out
}

// OutputDir 返回默认输出目录：程序目录下的 output/。
func OutputDir() string {
	return filepath.Join(AppDir(), "output")
}
