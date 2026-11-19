package ui

import (
	"fmt"
	"strings"

	"github.com/ANG13T/DroneXtract/deps"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"github.com/TwiN/go-color"
)

// statusLabel 返回带颜色的状态标签。
func statusLabel(s deps.Status) string {
	switch s {
	case deps.OK:
		return color.Ize(color.Green, i18n.S("就绪", "OK"))
	case deps.Missing:
		return color.Ize(color.Red, i18n.S("缺失", "MISSING"))
	case deps.Outdated:
		return color.Ize(color.Yellow, i18n.S("需升级", "OUTDATED"))
	}
	return color.Ize(color.Yellow, i18n.S("未知", "UNKNOWN"))
}

// PrintReport 以表格形式打印检测结果。
func PrintReport(r deps.Report) {
	helpers.GenTableHeader(i18n.S("环境依赖检查", "Environment Dependencies"))
	for _, it := range r.Items {
		title := i18n.S(it.TitleZH, it.TitleEN)
		helpers.GenRowString(title, statusLabel(it.Status)+"  "+helpers.TruncateText(it.Detail, 40))
	}
	helpers.GenRowString(i18n.S("配置目录", "Config dir"), deps.AppDir())
	helpers.GenRowString(i18n.S("配置文件", "Config file"), deps.ConfigFile())
	helpers.GenRowString(i18n.S("输出目录", "Output dir"), deps.OutputDir())
	helpers.GenTableFooter()
}

// ProblemLines 返回未就绪项的修复建议，供启动时提示。
func ProblemLines(r deps.Report) []string {
	var lines []string
	for _, it := range r.Problems() {
		title := i18n.S(it.TitleZH, it.TitleEN)
		hint := i18n.S(it.HintZH, it.HintEN)
		lines = append(lines, fmt.Sprintf("%s  %s  ->  %s", statusLabel(it.Status), title, hint))
	}
	return lines
}

// PromptInstall 询问并执行安装；返回是否执行了安装。
func PromptInstall() bool {
	r := deps.CheckAll()
	problems := r.Problems()
	if len(problems) == 0 {
		fmt.Println(color.Ize(color.Green, i18n.S("  所有依赖均已就绪。", "  All dependencies are ready.")))
		return false
	}

	fmt.Println(color.Ize(color.Yellow, i18n.S("  检测到以下依赖未就绪：", "  The following dependencies are not ready:")))
	for _, line := range ProblemLines(r) {
		fmt.Println("   " + line)
	}

	fmt.Print(i18n.S("\n 是否自动安装依赖？(y/N) > ", "\n Install dependencies now? (y/N) > "))
	var answer string
	if err := helpers.ReadLine(&answer); err != nil {
		return false
	}
	answer = strings.ToLower(strings.TrimSpace(answer))
	if answer != "y" && answer != "yes" && answer != "是" {
		fmt.Println(color.Ize(color.Yellow, i18n.S("  已跳过安装。", "  Installation skipped.")))
		return false
	}

	for _, it := range problems {
		switch it.Name {
		case "exiftool":
			runInstall("ExifTool", deps.InstallExifTool)
		case "config":
			fixConfig()
		case "go":
			fmt.Println(color.Ize(color.Yellow, i18n.S(
				"  Go 仅源码构建需要，本程序已可运行，跳过。",
				"  Go is only needed to build from source; skipping.")))
		}
	}

	fmt.Println()
	PrintReport(deps.CheckAll())
	return true
}

// runInstall 执行一个安装器并打印进度。
func runInstall(name string, fn func() deps.Result) {
	fmt.Println(color.Ize(color.Cyan, i18n.S("\n  正在安装 "+name+" ...", "\n Installing "+name+" ...")))
	res := fn()

	for _, s := range res.Steps {
		if s.OK {
			fmt.Println(color.Ize(color.Green, "   [√] "+s.Command))
		} else {
			fmt.Println(color.Ize(color.Red, "   [X] "+s.Command))
			if s.Output != "" {
				fmt.Println(color.Ize(color.Red, "       "+strings.Split(s.Output, "\n")[0]))
			}
		}
	}

	if res.Succeeded() {
		fmt.Println(color.Ize(color.Green, i18n.S("   安装完成。", "   Installation finished.")))
		return
	}
	if res.Err != nil {
		fmt.Println(color.Ize(color.Red, "   "+res.Err.Error()))
	}
	for _, dir := range deps.ExifToolDirHint() {
		_ = dir
	}
	fmt.Println(color.Ize(color.Yellow, i18n.S(
		"   安装未完成，请参考上面的输出手动处理。",
		"   Installation did not complete; see output above.")))
}

// fixConfig 生成 .env 配置文件。
func fixConfig() {
	path, written, err := deps.WriteConfig(false)
	if err != nil {
		fmt.Println(color.Ize(color.Red, i18n.S("   写入 .env 失败：", "   Failed to write .env: ")+err.Error()))
		return
	}
	if !written {
		fmt.Println(color.Ize(color.Yellow, i18n.S("   .env 已存在，保持不变：", "   .env already exists, kept as is: ")+path))
		return
	}
	fmt.Println(color.Ize(color.Green, i18n.S("   已生成配置：", "   Configuration written: ")+path))
}

// RunMenu 是「环境依赖」子菜单。
func RunMenu() {
	for {
		PrintReport(deps.CheckAll())
		fmt.Println(color.Ize(color.Blue, i18n.S(
			"\n    [ 1 ] 重新检测    [ 2 ] 安装/修复    [ 3 ] 写入 .env    [ 4 ] 输出 JSON    [ 0 ] 退出",
			"\n    [ 1 ] Re-check    [ 2 ] Install/Fix    [ 3 ] Write .env    [ 4 ] JSON report    [ 0 ] Exit")))

		fmt.Print(i18n.PromptOption())
		var choice string
		if err := helpers.ReadLine(&choice); err != nil {
			return
		}
		switch strings.TrimSpace(choice) {
		case "1":
			continue
		case "2":
			PromptInstall()
		case "3":
			fixConfig()
		case "4":
			fmt.Println(deps.CheckAll().JSON())
		case "0", "q", "Q":
			return
		}
	}
}
