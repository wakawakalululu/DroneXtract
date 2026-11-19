package helpers

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/ANG13T/DroneXtract/deps"
	"github.com/ANG13T/DroneXtract/i18n"
	"github.com/TwiN/go-color"
)

// stdinReader 全局复用同一个缓冲读取器，避免多次创建导致缓冲区中已读入的数据丢失。
var stdinReader = bufio.NewReader(os.Stdin)

// ReadLine 读取一行输入，供 main 等包复用同一输入缓冲。
func ReadLine(out *string) error {
	return readLine(out)
}

// readLine 读取一行输入并清理行尾的 CR/LF。Windows 下回车可能被当作行结束符，
// 也可能以 CRLF 形式残留，因此统一按 Split 处理。
func readLine(out *string) error {
	line, err := stdinReader.ReadString('\n')
	if err != nil && line == "" {
		return err
	}
	for strings.HasSuffix(line, "\r") {
		line = strings.TrimSuffix(line, "\r")
	}
	*out = strings.TrimRight(line, "\n")
	return nil
}

func IsLanguageCommand(s string) bool {
	switch strings.TrimSpace(s) {
	case "L", "l", "文", "A", "a":
		return true
	}
	return false
}

// ToggleLanguage 切换当前语言并打印提示。
func ToggleLanguage() {
	i18n.Toggle()
	fmt.Println(color.Ize(color.Yellow, i18n.S("  [文] 已切换为中文界面", "  [A] Switched to English")))
}

func FileInputString() string {
	fmt.Print(i18n.PromptFile())
	var selection string
	if err := ReadLine(&selection); err != nil {
		exitOnInputClosed()
	}
	if IsLanguageCommand(selection) {
		ToggleLanguage()
		return FileInputString()
	}
	return selection
}

func OutputPathString() string {
	fmt.Print(i18n.PromptOutput())
	var selection string
	if err := ReadLine(&selection); err != nil {
		exitOnInputClosed()
	}
	if IsLanguageCommand(selection) {
		ToggleLanguage()
		return OutputPathString()
	}
	return selection
}

// exitOnInputClosed 处理输入终止（如管道输入结束），避免菜单无限递归。
func exitOnInputClosed() {
	fmt.Println(color.Ize(color.Red, i18n.InputClosed()))
	os.Exit(1)
}

// Option 读取一个菜单选项；输入 L 可随时切换语言。
func Option(min int, max int) int {
	return OptionWithBack(min, max, true)
}

// OptionWithBack 读取选项。hasBack 为 true 时，max-1 表示“返回主菜单”；
// 供没有返回项的子菜单使用，避免把正常选项误判为返回。
func OptionWithBack(min int, max int, hasBack bool) int {
	fmt.Print(i18n.PromptOption())
	var selection string
	if err := ReadLine(&selection); err != nil {
		exitOnInputClosed()
	}
	if IsLanguageCommand(selection) {
		ToggleLanguage()
		return OptionWithBack(min, max, hasBack)
	}
	num, err := strconv.Atoi(selection)
	if err != nil {
		fmt.Println(color.Ize(color.Red, i18n.InvalidInput()))
		return OptionWithBack(min, max, hasBack)
	}
	if num == min {
		fmt.Println(color.Ize(color.Blue, i18n.Farewell()))
		os.Exit(1)
		return 0
	}
	if hasBack && num == max-1 {
		return -1
	}
	if num > min && num < max+1 {
		return num
	}
	fmt.Println(color.Ize(color.Red, i18n.InvalidInput()))
	return OptionWithBack(min, max, hasBack)
}

func GenTableHeader(name string) {
	fmt.Println(color.Ize(color.Blue, "\n    ╔══════════════════════════════════════════════════════════════════════════════╗"))
	fmt.Println(color.Ize(color.Blue, "    ║"+CenterText(name, 78)+"║"))
	fmt.Println(color.Ize(color.Blue, "    ╠══════════════════════════════════════════════════════════════════════════════╣"))
}

func CheckFileFormat(path string, exten string) bool {
	extension := strings.ToLower(filepath.Ext(path))
	return (extension == exten)
}

func GenRowString(intro string, input string) {
	var val = "    ║ " + FitText(intro+": "+input, 76) + " ║"
	fmt.Println(color.Ize(color.Blue, val))
}

func GenTableHeaderModified(name string) {
	fmt.Println(color.Ize(color.Blue, "    ╠══════════════════════════════════════════════════════════════════════════════╣"))
	fmt.Println(color.Ize(color.Blue, "    ║"+CenterText(name, 78)+"║"))
	fmt.Println(color.Ize(color.Blue, "    ╠══════════════════════════════════════════════════════════════════════════════╣"))
}

func GenTableFooter() {
	fmt.Println(color.Ize(color.Blue, "    ╚══════════════════════════════════════════════════════════════════════════════╝"))
}

func PrintLog(message string) {
	fmt.Println(color.Ize(color.Cyan, message))
}

func PrintError(message string) {
	fmt.Println(color.Ize(color.Red, "[ERROR] "+message))
}

func PrintErrorLog(message string, err error) {
	fmt.Println(color.Ize(color.Red, message))
	log.Println(color.Ize(color.Red, "[ERROR]"), err)
}

func PrintValidLog(message string) {
	fmt.Println(color.Ize(color.Green, message))
}

func PrintInvalidLog(message string) {
	fmt.Println(color.Ize(color.Red, message))
}

// loadEnv 加载配置：优先 .env，缺失时使用内嵌默认值。
func loadEnv() {
	deps.ApplyEnv()
}

func GetEnvVariable(key string) int {
	loadEnv()

	val, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		PrintError(i18n.S("环境变量无效：", "INVALID ENV VARIABLE: ") + key)
	}

	return val
}

func GetEnvVariances() []float64 {
	loadEnv()

	return stringToFloat64Array(os.Getenv("ANALYSIS_MAX_VARIANCE"))
}

func stringToFloat64Array(input string) []float64 {
	strValues := strings.Fields(input)
	values := make([]float64, len(strValues))

	for i, strValue := range strValues {
		value, err := strconv.ParseFloat(strValue, 64)
		if err != nil {
			fmt.Printf(i18n.S("解析第 %d 项数值时出错：%v\n", "Error parsing value at index %d: %v\n"), i, err)
			return nil
		}
		values[i] = value
	}

	return values
}

// CenterText 按终端显示宽度居中填充，使中英文混排的表格边框保持对齐。
func CenterText(text string, width int) string {
	return i18n.Center(text, width)
}

// FitText 按显示宽度补齐到指定列宽，超长内容自动截断。
func FitText(text string, width int) string {
	return i18n.Fit(text, width)
}

// TruncateText 按显示宽度截断文本，超出部分以省略号表示。
func TruncateText(text string, width int) string {
	return i18n.Truncate(text, width)
}

// FormatSingleLine 将文本压成单行，避免多行内容破坏表格边框。
func FormatSingleLine(text string) string {
	return i18n.SingleLine(text)
}
