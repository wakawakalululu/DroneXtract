// Package assets 把终端界面所需的文本资源编译进二进制，并提供程序目录的解析。
//
// 直接读取 txt/*.txt 会让程序只能在仓库根目录运行（换目录后菜单变成空白）。
// 这里用 go:embed 内嵌，并在缺失时回退到磁盘文件，方便二次修改文案。
package assets

import (
	"embed"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

//go:embed txt default.env
var files embed.FS

// DefaultEnv 返回内嵌的默认配置，供缺少 .env 时回退。
func DefaultEnv() string {
	if data, err := files.ReadFile("default.env"); err == nil {
		return string(data)
	}
	return ""
}

// Dir 返回程序所在目录。所有配置、输出与磁盘文案资源都以它为基准。
//
// 特殊处理：通过 go run 运行时，可执行文件位于系统临时目录，
// 此时回退到当前工作目录，避免把配置写进临时目录。
func Dir() string {
	if exe, err := os.Executable(); err == nil {
		if resolved, err := filepath.EvalSymlinks(exe); err == nil {
			exe = resolved
		}
		dir := filepath.Dir(exe)
		if !isTempDir(dir) {
			return dir
		}
	}
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "."
}

// isTempDir 判断目录是否为 go run 的临时产物目录。
func isTempDir(dir string) bool {
	temp := os.TempDir()
	if temp == "" {
		return false
	}
	absDir, err1 := filepath.Abs(dir)
	absTemp, err2 := filepath.Abs(temp)
	if err1 != nil || err2 != nil {
		return false
	}
	if !strings.HasPrefix(absDir, absTemp) {
		return false
	}
	return strings.Contains(absDir, "go-build")
}

// Text 返回 txt 目录下相对路径的文本内容。
//
// 查找顺序：
//  1. <程序目录>/txt/<rel>   随程序分发的可覆盖文案，便于用户改菜单
//  2. <工作目录>/txt/<rel>   源码目录内直接调试
//  3. 内嵌资源               保证任意位置运行都不会出现空菜单
func Text(rel string) string {
	for _, base := range []string{Dir(), "."} {
		if data, err := ioutil.ReadFile(filepath.Join(base, "txt", rel)); err == nil {
			return string(data)
		}
	}
	if data, err := files.ReadFile(path.Join("txt", rel)); err == nil {
		return string(data)
	}
	return ""
}
