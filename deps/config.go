package deps

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ANG13T/DroneXtract/assets"
	"github.com/joho/godotenv"
)

// DefaultConfig 返回内嵌的默认配置文本。
func DefaultConfig() string {
	return assets.DefaultEnv()
}

// ConfigPath 返回首个存在且可读的配置文件路径，找不到返回空串。
func ConfigPath() string {
	for _, p := range ConfigPaths() {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return ""
}

// LoadEnv 按优先级加载配置并返回键值与来源描述：
//  1. 配置文件（见 ConfigPaths，先命中者优先）
//  2. 内嵌默认配置（保证任意目录运行都有有效参数）
//
// 已存在的进程环境变量具有更高优先级，不被覆盖。
func LoadEnv() (map[string]string, string) {
	values := map[string]string{}
	source := "builtin default"

	if path := ConfigPath(); path != "" {
		if parsed, err := godotenv.Read(path); err == nil {
			for k, v := range parsed {
				values[k] = v
			}
			source = path
		}
	}

	defaults, _ := godotenv.Unmarshal(DefaultConfig())
	for k, v := range defaults {
		if _, ok := values[k]; !ok {
			values[k] = v
		}
	}

	// 进程环境变量优先级最高（忽略本程序自己写回的那些，避免误判来源）
	for k := range values {
		if appliedByUs[k] {
			continue
		}
		if env := os.Getenv(k); env != "" {
			values[k] = env
			source = "environment"
		}
	}

	return values, source
}

// ApplyEnv 把配置写入进程环境，供既有代码通过 os.Getenv 读取。
// 只在环境变量为空时写入，避免覆盖用户显式设置。
// appliedByUs 记录由本程序写回环境变量的键，
// 让后续 LoadEnv 能把它们与用户外部设置的环境变量区分开。
var appliedByUs = map[string]bool{}

func ApplyEnv() {
	values, _ := LoadEnv()
	for k, v := range values {
		if os.Getenv(k) == "" {
			if err := os.Setenv(k, v); err == nil {
				appliedByUs[k] = true
			}
		}
	}
}

// IntValue 读取整数配置，失败时返回 fallback。
func IntValue(key string, fallback int) int {
	values, _ := LoadEnv()
	raw := strings.TrimSpace(values[key])
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

// FloatList 读取以空格分隔的浮点列表，失败时返回 nil。
func FloatList(key string) []float64 {
	values, _ := LoadEnv()
	fields := strings.Fields(values[key])
	if len(fields) == 0 {
		return nil
	}
	out := make([]float64, 0, len(fields))
	for _, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil
		}
		out = append(out, v)
	}
	return out
}

// WriteConfig 把配置写入当前目录的 .env，返回写入路径。
// overwrite 为 false 时，已存在则跳过。
func WriteConfig(overwrite bool) (string, bool, error) {
	path := envPath()
	if !overwrite {
		if _, err := os.Stat(path); err == nil {
			return path, false, nil
		}
	}
	content := DefaultConfig()
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	// 配置可能写入用户目录（Linux 的 ~/.config/dronextract），需先建目录
	if dir := filepath.Dir(path); dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return path, false, err
		}
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return path, false, err
	}
	return path, true, nil
}
