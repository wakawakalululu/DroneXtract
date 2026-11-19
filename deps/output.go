package deps

import (
	"os"
	"path/filepath"
)

// EnsureParent 确保目标文件的父目录存在，避免 Linux 上因 ../output 不存在而写入失败。
// 返回规范化后的路径。
func EnsureParent(path string) (string, error) {
	if path == "" {
		return path, nil
	}
	dir := filepath.Dir(path)
	if dir == "" || dir == "." {
		return path, nil
	}
	if _, err := os.Stat(dir); err == nil {
		return path, nil
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return path, err
	}
	return path, nil
}
