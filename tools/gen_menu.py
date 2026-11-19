# -*- coding: utf-8 -*-
"""生成中英双语终端菜单文本。

菜单边框宽度固定为 122 列内容区；中文字符按 2 列计算，保证中英两套菜单
切换后边框完全对齐。每条菜单项一行中文、一行英文，程序读取整段文本渲染，
因此无需在 Go 侧维护大段横排 banner。
"""
import os
import unicodedata

INNER = 122


def w(s):
    return sum(2 if unicodedata.east_asian_width(c) in ("W", "F") else 1 for c in s)


def center(s):
    left = (INNER - w(s)) // 2
    return " " * left + s + " " * (INNER - w(s) - left)


def line(s):
    return "║" + " " + s + " " * (INNER - w(s) - 1) + "║"


TOP = "╔" + "═" * INNER + "╗"
MID = "╠" + "═" * INNER + "╣"
BOT = "╚" + "═" * INNER + "╝"


def menu(title_zh, title_en, items, back=None):
    rows = [TOP, "║" + center(title_zh) + "║", "║" + center(title_en) + "║", MID,
            "║" + " " * INNER + "║"]
    for key, zh, en in items:
        rows.append(line("[ " + key + " ] " + zh))
        rows.append(line("        " + en))
        rows.append("║" + " " * INNER + "║")
    rows.append(MID)
    rows.append(line("[ 0 ] 退出 DroneXtract"))
    rows.append(line("      Exit DroneXtract"))
    rows.append(BOT)
    return "\n".join(rows)


FILES = {}

FILES["txt/options.txt"] = menu(
    "DJI 无人机数字取证工具包", "DJI Drone Digital Forensics Toolkit",
    [
        ("1", "无人机文件解析", "Drone File Parsing"),
        ("2", "遥测数据映射", "Telemetry Mapping"),
        ("3", "隐写分析套件", "Steganography Suite"),
        ("4", "飞行与完整性分析", "Flight and Integrity Analysis"),
        ("5", "环境依赖检查 / 安装", "Dependencies Check / Install"),
        ("L", "切换语言", "Switch language"),
    ],
)

FILES["txt/fileparsing.txt"] = menu(
    "无人机文件解析", "Drone Files Parser",
    [
        ("1", "解析 CSV 文件", "Parse CSV File"),
        ("2", "解析 KML 文件", "Parse KML File"),
        ("3", "解析 GPX 文件", "Parse GPX File"),
        ("4", "返回主菜单", "Back to Main Menu"),
    ],
)

FILES["txt/telemetrymapping.txt"] = menu(
    "DJI 遥测数据映射", "DJI Telemetry Mapping",
    [
        ("1", "飞行路径分析", "Flight Path Analysis"),
        ("2", "遥测数据可视化", "Telemetry Visualizations"),
        ("3", "返回主菜单", "Back to Main Menu"),
    ],
)

FILES["txt/steganography.txt"] = menu(
    "DJI 隐写分析套件", "DJI Steganography Suite",
    [
        ("1", "EXIF 分析", "EXIF Analysis"),
        ("2", "DNG 分析", "DNG Analysis"),
        ("3", "SRT 分析", "SRT Analysis"),
        ("4", "XMP 分析", "XMP Analysis"),
        ("5", "返回主菜单", "Back to Main Menu"),
    ],
)

FILES["txt/steganography/exif.txt"] = menu(
    "EXIF 隐写分析", "EXIF Steganography",
    [
        ("1", "打印内容", "Print Contents"),
        ("2", "导出为 TXT", "Convert to TXT"),
        ("3", "导出为 CSV", "Convert to CSV"),
        ("4", "导出为 JSON", "Convert to JSON"),
        ("5", "返回主菜单", "Back to Main Menu"),
    ],
)

FILES["txt/steganography/dng.txt"] = menu(
    "DNG 隐写分析", "DNG Steganography",
    [
        ("1", "打印内容", "Print Contents"),
        ("2", "转换为 PNG", "Convert to PNG"),
        ("3", "导出为 TXT", "Convert to TXT"),
        ("4", "导出为 CSV", "Convert to CSV"),
        ("5", "导出为 JSON", "Convert to JSON"),
        ("6", "返回主菜单", "Back to Main Menu"),
    ],
)

FILES["txt/steganography/srt.txt"] = menu(
    "SRT 隐写分析", "SRT Steganography",
    [
        ("1", "打印内容", "Print Contents"),
        ("2", "导出为 JSON", "Convert to JSON"),
        ("3", "导出为 GEOJSON", "Convert to GEOJSON"),
        ("4", "导出为 CSV", "Convert to CSV"),
        ("5", "返回主菜单", "Back to Main Menu"),
    ],
)

FILES["txt/steganography/xmp.txt"] = menu(
    "XMP 隐写分析", "XMP Steganography",
    [
        ("1", "打印内容", "Print Contents"),
        ("2", "导出为 TXT", "Convert to TXT"),
        ("3", "导出为 CSV", "Convert to CSV"),
        ("4", "导出为 JSON", "Convert to JSON"),
        ("5", "返回主菜单", "Back to Main Menu"),
    ],
)

INFO = "\n".join([
    "",
    " " * 40 + "作者 Author: Angelina Tsuboi (G4LXY)",
    "",
    " " * 22 + "GitHub: github.com/ANG13T/DroneXtract | 网站 Site: angelinatsuboi.com",
    "",
])
FILES["txt/info.txt"] = INFO

BASE = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))

if __name__ == "__main__":
    for rel, text in FILES.items():
        with open(os.path.join(BASE, rel), "w", encoding="utf-8", newline="") as f:
            f.write(text + "\n")
        bad = [l for l in text.split("\n") if w(l) != INNER + 2]
        print(rel, "ok" if not bad else "WIDTH MISMATCH %r" % bad[:1])
