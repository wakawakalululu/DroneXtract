// Package i18n 提供运行时的中英文文本切换。
//
// 菜单边框宽度固定，中英两套菜单共用同一个渲染器，因此切换语言后边框仍然对齐。
package i18n

import (
	"os"
	"strings"
)

// Lang 表示当前界面语言。
type Lang int

const (
	// ZH 简体中文（默认）
	ZH Lang = iota
	// EN English
	EN
)

// current 是全局当前语言，默认中文。
var current = ZH

// Get 返回当前语言。
func Get() Lang { return current }

// Set 直接设置当前语言。
func Set(l Lang) { current = l }

// Parse 解析语言标识，支持 zh/cn/chinese/中文 与 en/english，无法识别时返回 false。
func Parse(s string) (Lang, bool) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "zh", "cn", "zh-cn", "chinese", "中文", "简体中文":
		return ZH, true
	case "en", "en-us", "english", "英文":
		return EN, true
	}
	return ZH, false
}

// FromEnv 读取 DX_LANG 环境变量，用于在启动前指定语言。
func FromEnv() (Lang, bool) {
	return Parse(os.Getenv("DX_LANG"))
}

// Toggle 在中英文之间切换，并返回切换后的语言。
func Toggle() Lang {
	if current == ZH {
		current = EN
	} else {
		current = ZH
	}
	return current
}

// Name 返回语言名称，用于状态提示。
func (l Lang) Name() string {
	if l == ZH {
		return "中文"
	}
	return "English"
}

// Icon 返回语言切换图标。
func (l Lang) Icon() string {
	if l == ZH {
		return "文"
	}
	return "A"
}

// IsEN 便于判断是否需要使用英文字段。
func (l Lang) IsEN() bool { return l == EN }

// S 从双语字符串中按当前语言取值。
// zh 为中文，en 为英文；缺少英文时回退到中文，避免出现空白。
func S(zh, en string) string {
	if current == EN && en != "" {
		return en
	}
	return zh
}

// ---- 通用交互文本 ----

func PromptOption() string { return S("\n 请输入选项 > ", "\n ENTER INPUT > ") }
func PromptFile() string   { return S("\n 请输入文件路径 > ", "\n ENTER FILE PATH > ") }
func PromptOutput() string { return S("\n 请输入输出路径 > ", "\n ENTER OUTPUT PATH > ") }

func InvalidInput() string { return S("  [!] 输入无效", "  [!] INVALID INPUT") }
func InputClosed() string {
	return S("  [!] 输入已结束，程序退出。", "  [!] INPUT CLOSED, EXITING")
}
func Farewell() string  { return S(" 再见，飞行愉快！", " Farewell and fly high!") }
func Cancelled() string { return S("  已取消", "  CANCELLED") }

// ExitLabel 返回菜单中“退出”一项的文案。
func ExitLabel() string { return S("退出 DroneXtract", "Exit DroneXtract") }

// BackLabel 返回菜单中“返回”一项的文案。
func BackLabel() string { return S("返回主菜单", "Back to Main Menu") }

// ---- 表格工具 ----

// DisplayWidth 返回字符串的终端显示宽度（CJK 全角字符算 2 列）。
func DisplayWidth(s string) int {
	width := 0
	for _, r := range s {
		if IsWide(r) {
			width += 2
		} else {
			width++
		}
	}
	return width
}

// IsWide 判断字符在终端中是否占两列。
func IsWide(r rune) bool {
	switch {
	case r >= 0x1100 && r <= 0x115F,
		r >= 0x2E80 && r <= 0x303E,
		r >= 0x3041 && r <= 0x33FF,
		r >= 0x3400 && r <= 0x4DBF,
		r >= 0x4E00 && r <= 0x9FFF,
		r >= 0xA000 && r <= 0xA4CF,
		r >= 0xAC00 && r <= 0xD7A3,
		r >= 0xF900 && r <= 0xFAFF,
		r >= 0xFE30 && r <= 0xFE6F,
		r >= 0xFF00 && r <= 0xFF60,
		r >= 0xFFE0 && r <= 0xFFE6,
		r >= 0x20000 && r <= 0x3FFFD:
		return true
	}
	return false
}

// Pad 按显示宽度右侧补空格。
func Pad(s string, width int) string {
	gap := width - DisplayWidth(s)
	if gap < 0 {
		gap = 0
	}
	return s + strings.Repeat(" ", gap)
}

// Truncate 按显示宽度截断，超出部分以省略号表示。
func Truncate(s string, width int) string {
	if DisplayWidth(s) <= width {
		return s
	}
	if width <= 1 {
		return ""
	}
	var b strings.Builder
	used := 0
	for _, r := range s {
		step := 1
		if IsWide(r) {
			step = 2
		}
		if used+step > width-1 {
			break
		}
		b.WriteRune(r)
		used += step
	}
	return b.String() + "…"
}

// Fit 先截断再补齐，保证显示宽度恰好为 width。
func Fit(s string, width int) string {
	return Pad(Truncate(s, width), width)
}

// Center 按显示宽度居中。
func Center(s string, width int) string {
	s = Truncate(s, width)
	rest := width - DisplayWidth(s)
	if rest < 0 {
		rest = 0
	}
	left := rest / 2
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", rest-left)
}

// SingleLine 压成单行，避免多行内容破坏表格边框。
func SingleLine(s string) string {
	s = strings.NewReplacer("\r\n", " ", "\n", " ", "\t", " ").Replace(s)
	return strings.Join(strings.Fields(s), " ")
}

// ---- 菜单渲染 ----

const menuInner = 122

// MenuItem 是菜单中的一项。
type MenuItem struct {
	Key   string
	Label string
}

// RenderMenu 渲染一个固定宽度的菜单框，中英文都能保持边框对齐。
func RenderMenu(title string, items []MenuItem, exitLabel string) string {
	top := "╔" + strings.Repeat("═", menuInner) + "╗"
	mid := "╠" + strings.Repeat("═", menuInner) + "╣"
	bottom := "╚" + strings.Repeat("═", menuInner) + "╝"

	rows := []string{
		top,
		"║" + Center(title, menuInner) + "║",
		mid,
		"║" + strings.Repeat(" ", menuInner) + "║",
	}
	for _, it := range items {
		line := "[ " + it.Key + " ] " + it.Label
		rows = append(rows, "║"+" "+Pad(line, menuInner-1)+"║")
		rows = append(rows, "║"+strings.Repeat(" ", menuInner)+"║")
	}
	rows = append(rows, mid)
	rows = append(rows, "║"+" "+Pad("[ 0 ] "+exitLabel, menuInner-1)+"║")
	rows = append(rows, bottom)
	return strings.Join(rows, "\n")
}

// HeaderPair 用于导出文件的双语表头。
type HeaderPair struct {
	ZH string
	EN string
}

// Headers 按当前语言生成导出文件的表头列，保证 CSV / JSON 等数据与界面语言一致。
func Headers(pairs []HeaderPair) []string {
	out := make([]string, 0, len(pairs))
	for _, p := range pairs {
		out = append(out, S(p.ZH, p.EN))
	}
	return out
}
