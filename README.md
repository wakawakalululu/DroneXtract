<p align="center">
  <img align="center" alt="DroneXtract logo" width="600" src="https://github.com/ANG13T/DroneXtract/blob/main/assets/graphic.png">
</p>

> **中文汉化 + 工程化增强版 / Chinese-localized fork**
>
> 本仓库基于 [ANG13T/DroneXtract](https://github.com/ANG13T/DroneXtract)（MIT License，
> 原作者 Angelina Tsuboi）二次开发，在原版基础上完成了中英双语界面、跨平台依赖自检与
> 自动安装、跨平台构建与安装脚本，以及若干稳定性修复。原始版权与许可证完整保留。
>
> This is a community fork of [ANG13T/DroneXtract](https://github.com/ANG13T/DroneXtract).
> Original work Copyright (c) 2023 Angelina Tsuboi, MIT licensed. See [CHANGES.md](CHANGES.md)
> for everything added or fixed in this fork.



## About
DroneXtract is a comprehensive digital forensics suite for DJI drones made with Golang. It can be used to analyze drone sensor values and telemetry data, visualize drone flight maps, audit for criminal activity, and extract pertinent data within multiple file formats. 

## Preview

<img alt="DroneXtract logo" width="600" src="https://github.com/ANG13T/DroneXtract/blob/main/assets/Preview.png">

## Features
DroneXtract features four main suites for drone forensics and auditing. They include the following:

### DJI File Parsing
You can visualize and extract information from DJI file formats such as CSV, KML, and GPX using the parsing tool.
The parsed information can be saved into an alternative file format when inputted an output file path.
The image below includes an example of a parsed file output and the type of data extracted from the file.

<img alt="DroneXtract logo" height="300" src="https://github.com/ANG13T/DroneXtract/blob/main/assets/demo-1.png">


### Steganography
Steganography refers to the process of revealing information stored within files.
The DroneXtract steganography suite allows you to extract telemetry and valuable data from image and video formats.
Additionally, the extracted data can be exported to four different file formats.

<img alt="DroneXtract logo" height="300" src="https://github.com/ANG13T/DroneXtract/blob/main/assets/demo-2.png">

### Telemetry Visualization
The telemetry visualization suite contains a flight path mapping generator and a telemetry graph visualizer.
The flight path mapping generator creates an image of a map indicating the locations the drone traveled to enroute and the path it took.
The telemetry graph visualizer plots a graph for each of the relevant telemetry or sensor values to be used for auditing purposes. 

<img alt="DroneXtract logo" width="600" src="https://github.com/ANG13T/DroneXtract/blob/main/assets/demo-3.png">

### Flight and Integrity Analysis
The flight and integrity analysis tool iterates through all the telemetry values the drone logged during its flight.
Once the values are collected, it calculates the maximum variance assumed by the value and checks for suspicious data gaps.
This tool can be used to check for anomalous data or any file corruption that may have taken place.

<img alt="DroneXtract logo" height="300" src="https://github.com/ANG13T/DroneXtract/blob/main/assets/demo-4.png">

## Install (Linux / macOS)

```bash
./build.sh all        # 产出 linux/amd64, linux/arm64, darwin/{amd64,arm64}, windows/amd64
./install.sh          # 安装到 ~/.local/bin + 装依赖 + 生成配置
./install.sh --prefix /opt/dronextract
./install.sh --no-deps
```

`install.sh` 会：检测并安装 ExifTool（自动处理 `sudo`，非 root 时给出可复制的命令）、
在程序目录生成 `.env`、把二进制装到 `~/.local/bin/dronextract`。

Prebuilt binaries are static ELF/Mach-O files with no runtime dependencies.

### Linux 适配要点

- **权限**：包管理器安装会自动判断 root / `sudo`；在无 TTY（CI、管道）环境下不会挂起等密码，而是直接打印 `sudo ...` 命令。
- **配置位置**：**始终保存在程序所在目录**（二进制旁边的 `.env`），不写用户主目录，整目录拷贝即可迁移。
  查找优先级：`$DRONEXTRACT_ENV` → `<程序目录>/.env` → `<工作目录>/.env`（仅作覆盖）→ 旧版本用户目录（只读兼容）→ 内嵌默认值。
- **输出目录**：默认写到 `<程序目录>/output/`，并自动创建；不再依赖工作目录或预置的 `output/`。
- **终端**：尊重 `NO_COLOR`；中文按 2 列计算宽度，`--check` 等表格在 Linux 终端不会错位。
- **构建**：`build.sh` 产出静态链接二进制并补上执行位（交叉编译产物默认可执行位会丢失）。

## Usage
To build from source, you will need Go installed.

```bash
$ export GO111MODULE=on
$ go get ./...
$ go run main.go
```

## Language / 语言

DroneXtract ships with a built-in bilingual interface (English / 简体中文).
The default language is **Chinese**; every menu also prints the English label underneath.

- Press `L` at any prompt (main menu, sub-menu, file path, output path) to switch language.
- Start directly in a language: `./DroneXtract --lang=en` (or `DX_LANG=en`); default is Chinese.
  `--lang` also accepts `zh`.
- The language icon in the title bar shows the active language: `[ 文 ] 中文` or `[ A ] English`.
- The choice also applies to exported data: CSV headers and table titles are written
  in the current language, e.g. `时间戳,返航点纬度,...` vs `TIMESTAMP,HOME.LATITUDE,...`.

带中英双语界面：标题栏的 `[ 文 ] / [ A ]` 即语言切换图标，在任意输入提示下按 `L`
即可切换，界面文案与导出数据的表头会同步切换。启动时可用 `--lang=en` 或
`--lang=zh` 指定语言（也可设 `DX_LANG`），默认中文。

> 菜单文本可通过 `txt/` 目录直接修改；改完执行 `./sync_txt.sh` 同步到内嵌资源。
> 程序在任意目录运行都能显示完整菜单与默认配置（`assets/default.env`）。

## Configuration / 配置

配置固定保存在**程序所在目录**的 `.env`。没有 `.env` 时程序自动使用内嵌默认值，
不影响任何功能；需要自定义时复制示例文件：

```bash
cp .env.example .env    # 然后编辑 .env
```

| 变量 | 含义 |
|---|---|
| `TELEMETRY_VIS_DOWNSAMPLE` | 遥测可视化降采样点数 |
| `FLIGHT_MAP_DOWNSAMPLE` | 飞行路径地图降采样点数 |
| `ANALYSIS_DOWNSAMPLE` | 飞行分析降采样点数 |
| `ANALYSIS_MAX_VARIANCE` | 27 项指标的方差阈值（空格分隔） |

覆盖优先顺序：`$DRONEXTRACT_ENV` → 程序目录 `.env` → 工作目录 `.env` → 内嵌默认值。

## Dependencies / 依赖

DroneXtract needs very little at runtime. Everything is checked, and fixed, from one place.

| Dependency | Needed for | Auto-install |
|---|---|---|
| `.env` config | downsampling + variance thresholds | yes (writes defaults) |
| [ExifTool](https://exiftool.org/) | EXIF / DNG / XMP analysis | yes (`winget` / `choco` / `brew` / `apt` / `dnf` / `pacman`) |
| Go 1.20+ | building from source only | n/a (binary needs nothing) |

Check or install from the CLI:

```bash
./DroneXtract --check      # table report, exit code 2 when something is missing
./DroneXtract --json       # machine readable report
./DroneXtract --install    # auto-install missing dependencies
```

Or interactively: main menu → `[ 5 ]` → re-check / install / write `.env` / JSON report.

依赖检测与安装集中在主菜单 `[ 5 ] 环境依赖检查 / 安装`：可重新检测、自动安装、
生成 `.env`、输出 JSON 报告。命令行等价入口为 `--check` / `--json` / `--install`。

- 缺少 `.env` 时自动使用内嵌默认配置，不会因读不到配置而让采样率变成 0。
- 缺少 ExifTool 时，EXIF/DNG/XMP 分析会先给出可操作提示，而不是底层报错。
- 没有可用包管理器时，会给出官方下载地址与手动安装指引。

## Airdata Usage
In order to parse out DJI flight .TXT logs, please use [Airdata's Flight Data Analysis tool](https://app.airdata.com/main?a=upload).
- Airdata CSV output file can be used for the CSV parser, flight path map, and telemetry visualizations
- Airdata KML output file can be used for the KML parser
- Airdata GPX output file can be used for GPX parser 

## Configuration Detail (upstream)
There are a set of environment variables utilized in DroneXtract. In order to tailor the values to your specific drone / investigation scenario, you can go to the `.env` file and adjust the following values:

`TELEMETRY_VIS_DOWNSAMPLE` downsampling number for values to be used for telemetry visualization

`FLIGHT_MAP_DOWNSAMPLE` downsampling number for values to be used for flight path mapping

`ANALYSIS_DOWNSAMPLE` downsampling number for values to be used for integrity analysis

`ANALYSIS_MAX_VARIANCE` maximum variance allowed between max and min value for analysis values


## Learning and Resources
To learn more about DJI drone digital forensics and the features of DroneXtract, refer to [this article](https://medium.com/@angelinatsuboi/a-comprehensive-guide-to-digital-forensics-with-dji-drones-fd7ef5af2891).

## Contributing
DroneXtract is open to any contributions. Please fork the repository and make a pull request with the features or fixes you want to be implemented.

### Testing
DroneXtract contains testing files for all four suites. 
An example command set for testing the steganography suite
```bash
$ cd steganography
$ go test
```
Example files for testing are included in the `test-data` directory.
If a specific test case includes an output such as a file or image, it will be stored within the `output` directory.
There are already some files contained in the `output` for reference purposes.

## Upcoming
- DUML parser for firmware integrity checking
- DJI Flight Log TXT parsing for the parsing suite
- GEOJSON parsing output for SRT files in the steganography suite

## Support
If you enjoyed DroneXtract, please consider [becoming a sponsor](https://github.com/sponsors/ANG13T) in order to fund my future projects. 

To check out my other works, visit my [GitHub profile](https://github.com/ANG13T).
