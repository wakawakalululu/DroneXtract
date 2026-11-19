package deps

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// Step 是安装过程中的一步，用于向 UI 报告进度。
type Step struct {
	Command string
	OK      bool
	Output  string
	// Skipped 表示这一步因权限或环境原因未执行。
	Skipped bool
}

// Result 是安装结果。
type Result struct {
	Name    string
	Started bool
	Skipped bool
	Steps   []Step
	Err     error
	// FixHint 是失败后建议用户手动执行的命令，便于复制粘贴。
	FixHint string
}

// Succeeded 表示安装流程执行成功。
func (r Result) Succeeded() bool {
	if r.Err != nil || r.Skipped {
		return false
	}
	if len(r.Steps) == 0 {
		return false
	}
	for _, s := range r.Steps {
		if !s.OK {
			return false
		}
	}
	return true
}

// Summary 返回可读的一行摘要。
func (r Result) Summary() string {
	if r.Skipped {
		return "skipped: " + r.Name
	}
	if r.Err != nil {
		return "failed: " + r.Err.Error()
	}
	if len(r.Steps) == 0 {
		return "no action"
	}
	last := r.Steps[len(r.Steps)-1]
	if last.OK {
		return "ok: " + last.Command
	}
	return "failed: " + last.Command
}

// isRoot 判断当前进程是否具备 root 权限。
func isRoot() bool {
	if runtime.GOOS == "windows" {
		return false
	}
	return os.Geteuid() == 0
}

// hasSudo 判断系统是否提供 sudo。
func hasSudo() bool {
	return lookPath("sudo") != ""
}

// isInteractive 判断当前是否有可用的终端交互（有 TTY）。
// 在 CI 或管道输入下 false，此时不应尝试需要密码的提权。
func isInteractive() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeCharDevice != 0
}

// packagedInstall 描述一种包管理器的安装方式。
type packagedInstall struct {
	manager string   // 包管理器名，用于提示
	bin     string   // 实际可执行文件
	args    []string // 安装参数
	refresh []string // 安装前的索引刷新（可为空）
	pkg     string   // 包名，用于手动提示
}

// linuxCandidates 按发行版常见顺序列出可用的包管理器。
func linuxCandidates() []packagedInstall {
	if _, err := os.Stat("/etc/debian_version"); err == nil {
		return []packagedInstall{
			{"apt", "apt-get", []string{"install", "-y", "libimage-exiftool-perl"}, []string{"update"}, "libimage-exiftool-perl"},
			{"apt", "apt", []string{"install", "-y", "libimage-exiftool-perl"}, nil, "libimage-exiftool-perl"},
			{"snap", "snap", []string{"install", "exiftool"}, nil, "exiftool"},
		}
	}
	return []packagedInstall{
		{"dnf", "dnf", []string{"install", "-y", "perl-Image-ExifTool"}, nil, "perl-Image-ExifTool"},
		{"yum", "yum", []string{"install", "-y", "perl-Image-ExifTool"}, nil, "perl-Image-ExifTool"},
		{"pacman", "pacman", []string{"-S", "--noconfirm", "perl-image-exiftool"}, nil, "perl-image-exiftool"},
		{"zypper", "zypper", []string{"install", "-y", "exiftool"}, nil, "exiftool"},
		{"apk", "apk", []string{"add", "exiftool"}, nil, "exiftool"},
		{"snap", "snap", []string{"install", "exiftool"}, nil, "exiftool"},
	}
}

// InstallExifTool 尝试自动安装 ExifTool。
//
// 策略（按平台）：
//   - Windows：winget → choco → 官方地址
//   - macOS：brew → 官方地址
//   - Linux：按发行版选择包管理器；非 root 时用 sudo（仅在有 TTY 时尝试，
//     避免在 CI / 管道环境里卡在密码提示）；无 root 且无 sudo 时直接给出
//     可复制的 sudo 命令，而不是执行注定失败的命令。
//
// 所有命令均以非交互方式执行，失败时保留输出供用户排查。
func InstallExifTool() Result {
	res := Result{Name: "exiftool"}

	switch runtime.GOOS {
	case "windows":
		return installWindows(res)
	case "darwin":
		return installDarwin(res)
	default:
		return installLinux(res)
	}
}

func installWindows(res Result) Result {
	if p := lookPath("winget"); p != "" {
		res.Started = true
		res.Steps = append(res.Steps, runStep(p, "install", "--id", "OliverBetz.ExifTool",
			"--accept-source-agreements", "--accept-package-agreements",
			"--disable-interactivity", "--silent"))
		if res.Steps[0].OK && CheckExifTool().Ready() {
			return res
		}
	}
	if p := lookPath("choco"); p != "" {
		res.Started = true
		res.Steps = append(res.Steps, runStep(p, "install", "exiftool", "-y", "--no-progress"))
		return res
	}
	if !res.Started {
		res.Err = errors.New("未检测到 winget 或 choco")
		res.FixHint = "https://exiftool.org/install.html"
	}
	return res
}

func installDarwin(res Result) Result {
	if p := lookPath("brew"); p != "" {
		res.Started = true
		res.Steps = append(res.Steps, runStep(p, "install", "exiftool"))
		return res
	}
	res.Err = errors.New("未检测到 Homebrew")
	res.FixHint = "https://exiftool.org/install.html"
	return res
}

func installLinux(res Result) Result {
	candidates := linuxCandidates()

	var chosen *packagedInstall
	for i := range candidates {
		if lookPath(candidates[i].bin) != "" {
			chosen = &candidates[i]
			break
		}
	}
	if chosen == nil {
		res.Err = errors.New("未找到可用的包管理器")
		res.FixHint = "https://exiftool.org/install.html"
		return res
	}

	// 已装但未在 PATH：检查常见目录后给出 PATH 提示。
	defer func() {
		if CheckExifTool().Ready() {
			res.Err = nil
			res.FixHint = ""
		}
	}()

	// 决定是否提权
	var runner string
	var prefix []string
	if isRoot() {
		runner = chosen.bin
	} else if hasSudo() {
		if !isInteractive() {
			// 无 TTY 时 sudo 会挂起等待密码，改为直接给出命令
			res.Skipped = true
			res.Err = errors.New("当前非 root 且无交互终端，无法自动提权")
			res.FixHint = buildSudoHint(chosen)
			return res
		}
		runner = lookPath("sudo")
		prefix = []string{"-n", chosen.bin} // -n: 不提示密码，失败立即返回
	} else {
		res.Skipped = true
		res.Err = errors.New("当前用户无 root 权限且未安装 sudo")
		res.FixHint = buildSudoHint(chosen)
		return res
	}

	res.Started = true

	// 先刷新包索引，失败不阻断安装（离线镜像或已是最新）
	if chosen.refresh != nil {
		args := append(append([]string{}, prefix...), chosen.refresh...)
		step := runStep(runner, args...)
		// 刷新失败只记录，不视为致命错误
		res.Steps = append(res.Steps, step)
	}

	args := append(append([]string{}, prefix...), chosen.args...)
	step := runStep(runner, args...)
	res.Steps = append(res.Steps, step)

	if step.OK {
		if CheckExifTool().Ready() {
			return res
		}
		// 安装成功但 PATH 未刷新
		res.Err = nil
		res.FixHint = newPathHint()
		return res
	}

	res.FixHint = buildSudoHint(chosen)
	return res
}

// buildSudoHint 生成用户可直接复制的提权安装命令。
func buildSudoHint(c *packagedInstall) string {
	return "sudo " + c.bin + " " + strings.Join(c.args, " ")
}

// newPathHint 提示用户刷新 shell 的 PATH 缓存。
func newPathHint() string {
	return "hash -r  # 或重开终端，让新装的 exiftool 进入 PATH"
}

// runStep 执行一条安装命令并捕获输出。
func runStep(bin string, args ...string) Step {
	step := Step{Command: bin + " " + strings.Join(args, " ")}
	cmd := exec.Command(bin, args...)
	cmd.Env = append(os.Environ(),
		"DEBIAN_FRONTEND=noninteractive",
		"APT_LISTCHANGES_FRONTEND=none",
		"NEEDRESTART_MODE=a",
	)
	// 独立进程组，避免包管理器的子进程继承终端信号
	cmd.SysProcAttr = sysProcAttr()

	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if len(text) > 2000 {
		text = text[len(text)-2000:]
	}
	step.Output = text

	if err != nil {
		// sudo -n 在无免密授权时返回 1，并带有明确提示
		if strings.Contains(text, "a password is required") ||
			strings.Contains(text, "no tty present") ||
			strings.Contains(text, "terminal is required") {
			step.Output = "需要 root 权限但无法自动提权"
			step.Skipped = true
		}
	}
	step.OK = err == nil
	return step
}

// ExecutablePath 返回当前可执行文件路径，供提示用户加入 PATH。
func ExecutablePath() string {
	if p, err := os.Executable(); err == nil {
		return p
	}
	return ""
}

// ExifToolDirHint 返回 ExifTool 可能的安装目录，供提示用户加入 PATH。
func ExifToolDirHint() []string {
	switch runtime.GOOS {
	case "windows":
		return []string{
			filepath.Join(os.Getenv("ProgramFiles"), "ExifTool"),
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Microsoft", "WinGet", "Packages"),
		}
	case "darwin":
		return []string{"/opt/homebrew/bin", "/usr/local/bin", "/usr/bin"}
	default:
		return []string{"/usr/bin", "/usr/local/bin", "/snap/bin", "/var/lib/snapd/snap/bin"}
	}
}
