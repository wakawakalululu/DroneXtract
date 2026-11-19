package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ANG13T/DroneXtract/analysis"
	"github.com/ANG13T/DroneXtract/assets"
	"github.com/ANG13T/DroneXtract/deps"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"github.com/ANG13T/DroneXtract/parsing"
	"github.com/ANG13T/DroneXtract/steganography"
	"github.com/ANG13T/DroneXtract/telemetry"
	"github.com/ANG13T/DroneXtract/ui"
	"github.com/TwiN/go-color"
	"github.com/iskaa02/qalam/gradient"
)

var category_banners = []string{"fileparsing.txt", "telemetrymapping.txt", "steganography.txt"}
var option_values = []int{5, 4, 6}

// LangBadge 返回标题栏上的语言图标：中文显示「文」，英文显示「A」。
func LangBadge() string {
	return color.Ize(color.Yellow, "[ "+i18n.Get().Icon()+" ]") + " " + i18n.Get().Name()
}

func Banner() {
	banner := assets.Text("banner.txt")
	info := assets.Text("info.txt")
	options := assets.Text("options.txt")
	g, _ := gradient.NewGradient("#1179ef", "cyan")
	g.Print("\n" + string(banner))
	fmt.Println(color.Ize(color.Cyan, info))
	fmt.Println("   " + LangBadge() + color.Ize(color.Cyan, i18n.S("   按 L 切换语言", "   Press L to switch language")))
	fmt.Println(color.Ize(color.Cyan, string(options)+"\n"))
}

func Option() {
	fmt.Print(i18n.PromptOption())
	var selection string
	if scanErr := helpers.ReadLine(&selection); scanErr != nil {
		fmt.Println(color.Ize(color.Red, i18n.InputClosed()))
		os.Exit(1)
	}
	if helpers.IsLanguageCommand(selection) {
		helpers.ToggleLanguage()
		Banner()
		Option()
		return
	}
	num, err := strconv.Atoi(selection)
	if err != nil {
		fmt.Println(color.Ize(color.Red, i18n.InvalidInput()))
		Option()
		return
	}
	if num == 0 {
		fmt.Println(color.Ize(color.Blue, i18n.Farewell()))
		os.Exit(1)
	}
	if num > 0 && num < 6 {
		DisplayOption(num)
		return
	}
	if num == 6 {
		ui.RunMenu()
		Banner()
		Option()
		return
	}
	fmt.Println(color.Ize(color.Red, i18n.InvalidInput()))
	Option()
}

func DisplayOption(x int) {
	if x == 0 {
		fmt.Println(color.Ize(color.Blue, i18n.Farewell()))
		os.Exit(1)
	} else if x == 5 {
		ui.RunMenu()
		Banner()
		Option()
	} else if x == 6 {
		Banner()
		Option()
	} else if x > 0 && x < 6 {
		DisplayOptionInformation(x)
	} else {
		helpers.PrintError(i18n.InvalidInput())
		Option()
	}
}

func DisplayOptionInformation(option int) {
	returnVal := 0
	if option == 4 {
		analysis.ExecuteAnalysis()
		Banner()
		Option()
	} else {
		contents := assets.Text(category_banners[option-1])
		fmt.Println(color.Ize(color.Cyan, contents))
		returnVal = helpers.Option(0, option_values[option-1])
	}

	if returnVal == -1 {
		Banner()
		Option()
	} else {
		if option == 1 {
			parsing.ExecuteParser(returnVal)
			Banner()
			Option()
		} else if option == 2 {
			telemetry.ExecuteTelemetry(returnVal)
			Banner()
			Option()
		} else {
			steganography.ExecuteSteganography(returnVal)
			Banner()
			Option()
		}
	}

}

// applyStartupLang 处理启动语言：--lang 参数优先，其次 DX_LANG 环境变量。
func applyStartupLang() {
	for i, arg := range os.Args[1:] {
		var value string
		switch {
		case arg == "--lang" || arg == "-l":
			if i+2 < len(os.Args) {
				value = os.Args[i+2]
			}
		case strings.HasPrefix(arg, "--lang="):
			value = strings.TrimPrefix(arg, "--lang=")
		default:
			continue
		}
		if lang, ok := i18n.Parse(value); ok {
			i18n.Set(lang)
			return
		}
		fmt.Println(color.Ize(color.Red, "--lang/--lang= 仅支持 zh 或 en (only zh / en)"))
		return
	}
	if lang, ok := i18n.FromEnv(); ok {
		i18n.Set(lang)
	}
}

// runCheck 处理 --check / --json，打印依赖报告后退出。
func runCheck(asJSON bool) {
	if asJSON {
		fmt.Println(deps.CheckAll().JSON())
		return
	}
	r := deps.CheckAll()
	ui.PrintReport(r)
	if r.HasProblem() {
		for _, line := range ui.ProblemLines(r) {
			fmt.Println("  " + line)
		}
		fmt.Println()
		fmt.Println(color.Ize(color.Yellow, i18n.S(
			"  提示：运行 ./DroneXtract --install 自动安装，或在主菜单选 [ 5 ]。",
			"  Tip: run ./DroneXtract --install, or choose [ 5 ] in the main menu.")))
		os.Exit(2)
	}
}

func usageText() string {
	return strings.TrimSpace(`
  DroneXtract - DJI 无人机数字取证工具 / DJI drone digital forensics toolkit

  用法 / Usage:
    DroneXtract                启动交互界面 / start interactive UI
    DroneXtract --lang=en      以英文界面启动 / start in English
    DroneXtract --check        检测环境依赖 / check dependencies
    DroneXtract --json         以 JSON 输出依赖报告 / JSON dependency report
    DroneXtract --install      自动安装缺失依赖 / install missing dependencies
    DroneXtract --help         显示本帮助 / show this help

  环境变量 / Environment:
    DX_LANG=zh|en              界面语言 / UI language
`)
}

func main() {
	applyStartupLang()

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--check", "-c":
			runCheck(false)
			return
		case "--json":
			runCheck(true)
			return
		case "--install":
			deps.ApplyEnv()
			ui.PromptInstall()
			return
		case "--help", "-h":
			fmt.Println(usageText())
			return
		}
	}

	deps.ApplyEnv()
	Banner()
	Option()
}
