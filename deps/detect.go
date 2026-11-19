// Package deps 负责运行环境依赖的检测、安装与配置。
//
// 目前 DroneXtract 依赖三类外部条件：
//  1. Go 工具链（仅源码运行需要，二进制运行时无需）
//  2. ExifTool（EXIF / DNG / XMP 分析必需）
//  3. .env 配置（采样率与方差阈值，缺失时用内嵌默认值）
//
// 所有检查都返回结构化结果，UI 层据此渲染；安装动作只在用户明确确认后执行。
package deps

import (
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
	"time"
)

// Status 表示单项依赖的状态。
type Status int

const (
	// OK 已就绪
	OK Status = iota
	// Missing 未安装
	Missing
	// Outdated 版本过低
	Outdated
	// Unknown 无法判断
	Unknown
)

// Item 是一项依赖的检测结果。
type Item struct {
	Name    string // 内部标识
	TitleZH string
	TitleEN string
	Status  Status
	Detail  string // 版本号、路径或错误摘要
	HintZH  string // 修复建议
	HintEN  string
	Fixable bool // 是否支持自动安装
}

// Ready 表示该项可用。
func (i Item) Ready() bool { return i.Status == OK }

// Report 汇总一次完整检测。
type Report struct {
	OS       string
	Arch     string
	Items    []Item
	Duration time.Duration
}

// HasProblem 表示存在未就绪项。
func (r Report) HasProblem() bool {
	for _, it := range r.Items {
		if !it.Ready() {
			return true
		}
	}
	return false
}

// Problems 返回所有未就绪项。
func (r Report) Problems() []Item {
	out := make([]Item, 0, len(r.Items))
	for _, it := range r.Items {
		if !it.Ready() {
			out = append(out, it)
		}
	}
	return out
}

// commandTimeout 限制外部命令探测时间，避免在异常 PATH 上卡住。
const commandTimeout = 5 * time.Second

// lookPath 在 PATH 中查找可执行文件，Windows 下自动补 .exe。
func lookPath(names ...string) string {
	for _, name := range names {
		if p, err := exec.LookPath(name); err == nil {
			return p
		}
		if runtime.GOOS == "windows" && !strings.HasSuffix(strings.ToLower(name), ".exe") {
			if p, err := exec.LookPath(name + ".exe"); err == nil {
				return p
			}
		}
	}
	return ""
}

// output 执行命令并返回首行输出，失败返回错误。
func output(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	done := make(chan struct{})
	var out []byte
	var err error
	go func() {
		out, err = cmd.CombinedOutput()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(commandTimeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return "", os.ErrDeadlineExceeded
	}
	text := strings.TrimSpace(string(out))
	if idx := strings.IndexAny(text, "\r\n"); idx >= 0 {
		text = text[:idx]
	}
	return text, err
}

// CheckExifTool 检测 ExifTool。
func CheckExifTool() Item {
	item := Item{
		Name:    "exiftool",
		TitleZH: "ExifTool（EXIF / DNG / XMP 分析）",
		TitleEN: "ExifTool (EXIF / DNG / XMP analysis)",
		HintZH:  "运行「安装依赖」自动安装，或手动执行：winget install ExifTool",
		HintEN:  "Run \"Install dependencies\", or: winget install ExifTool",
		Fixable: true,
	}
	path := lookPath("exiftool", "exiftool.exe")
	if path == "" {
		item.Status = Missing
		item.Detail = "not found in PATH"
		return item
	}
	if v, err := output(path, "-ver"); err == nil && v != "" {
		item.Status = OK
		item.Detail = "v" + v + " @ " + path
		return item
	}
	item.Status = Unknown
	item.Detail = path
	return item
}

// CheckGo 检测 Go 工具链（源码构建需要）。
func CheckGo() Item {
	item := Item{
		Name:    "go",
		TitleZH: "Go 工具链（源码构建）",
		TitleEN: "Go toolchain (build from source)",
		HintZH:  "二进制运行无需 Go；如需源码构建请安装 Go 1.20+",
		HintEN:  "Not needed for the binary; install Go 1.20+ to build from source",
	}
	path := lookPath("go")
	if path == "" {
		item.Status = Missing
		item.Detail = "not found in PATH"
		return item
	}
	if v, err := output(path, "version"); err == nil {
		item.Status = OK
		item.Detail = v
		return item
	}
	item.Status = Unknown
	item.Detail = path
	return item
}

// CheckConfig 检测 .env 配置是否完整可用。
func CheckConfig() Item {
	item := Item{
		Name:    "config",
		TitleZH: ".env 运行配置",
		TitleEN: ".env runtime configuration",
		Fixable: true,
	}
	required := []string{
		"TELEMETRY_VIS_DOWNSAMPLE",
		"FLIGHT_MAP_DOWNSAMPLE",
		"ANALYSIS_DOWNSAMPLE",
		"ANALYSIS_MAX_VARIANCE",
	}
	values, source := LoadEnv()
	var missing []string
	for _, key := range required {
		if strings.TrimSpace(values[key]) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) == 0 {
		item.Status = OK
		if source == "builtin default" {
			item.Detail = "builtin default (" + ConfigFile() + ")"
		} else {
			item.Detail = source
		}
		return item
	}
	item.Status = Missing
	item.Detail = "missing: " + strings.Join(missing, ", ")
	item.HintZH = "运行「修复配置」生成 .env，或继续使用内嵌默认值"
	item.HintEN = "Run \"Fix configuration\" to write .env, or keep built-in defaults"
	return item
}

// CheckAll 执行完整检测。
func CheckAll() Report {
	start := time.Now()
	items := []Item{CheckConfig(), CheckExifTool(), CheckGo()}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Name < items[j].Name })
	return Report{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		Items:    items,
		Duration: time.Since(start),
	}
}

// JSON 便于把检测结果落盘或交给外部脚本。
func (r Report) JSON() string {
	type jsonItem struct {
		Name   string `json:"name"`
		Status string `json:"status"`
		Detail string `json:"detail"`
	}
	payload := struct {
		OS    string     `json:"os"`
		Arch  string     `json:"arch"`
		Items []jsonItem `json:"items"`
	}{OS: r.OS, Arch: r.Arch}
	for _, it := range r.Items {
		payload.Items = append(payload.Items, jsonItem{it.Name, statusText(it.Status), it.Detail})
	}
	data, _ := json.MarshalIndent(payload, "", "  ")
	return string(data)
}

func statusText(s Status) string {
	switch s {
	case OK:
		return "ok"
	case Missing:
		return "missing"
	case Outdated:
		return "outdated"
	}
	return "unknown"
}

// envPath 返回 .env 的写入位置：始终写入程序所在目录。
func envPath() string {
	return ConfigFile()
}
