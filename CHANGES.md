# 变更说明 / Changelog

本仓库基于 [ANG13T/DroneXtract](https://github.com/ANG13T/DroneXtract) 二次开发。
原始代码版权归 Angelina Tsuboi 所有，遵循 MIT License。以下按主题汇总本分支的改动。

---

## 1. 中英双语界面

- 新增 `i18n` 包：运行期语言状态、`S(zh, en)` 取值、菜单渲染、CJK 等宽对齐。
- 8 套终端菜单改为**上下双语**（每项一行中文、一行英文），
  中文字符按 2 列计算宽度，切换语言后边框仍严格对齐。
- 提示语、表格标题、字段名、错误信息、飞行分析统计全部双语化。
- 遥测可视化 45 项指标、飞行分析 27 项指标建立双语对照表。
- 语言切换图标：标题栏显示 `[ 文 ] 中文` 或 `[ A ] English`，
  在**任意输入提示**下按 `L` 即可切换。
- 启动可指定语言：`--lang=en` / `--lang=zh` / `DX_LANG=en`，默认中文。

## 2. 导出数据双语

- 新增 `i18n.Headers([]HeaderPair)`，CSV 表头随当前语言切换：
  - 中文：`时间戳,返航点纬度,返航点经度,日期,GPS 纬度,...`
  - 英文：`TIMESTAMP,HOME.LATITUDE,HOME.LONGITUDE,DATE,GPS.LATITUDE,...`
- JSON / GeoJSON 的机器字段名保持不变，避免破坏下游解析。

## 3. 依赖检测与自动安装

新增 `deps` 包，集中处理第三方依赖：

| 依赖 | 用途 | 处理方式 |
|---|---|---|
| `.env` 配置 | 采样率 + 方差阈值 | 缺失时使用内嵌默认值，可一键写盘 |
| ExifTool | EXIF / DNG / XMP 分析 | 自动安装，支持多包管理器 |
| Go 1.20+ | 仅源码构建 | 检测并提示 |

命令行入口：

```bash
./DroneXtract --check      # 表格报告，缺依赖时退出码 2
./DroneXtract --json       # 机器可读报告
./DroneXtract --install    # 自动安装缺失依赖
./DroneXtract --help
```

交互入口：主菜单 `[ 5 ] 环境依赖检查 / 安装`（重新检测 / 安装修复 / 写配置 / JSON 报告）。

安装器按平台分流：

- **Windows**：`winget` → `choco` → 官方地址
- **macOS**：`brew` → 官方地址
- **Linux**：按发行版选择 `apt` / `dnf` / `yum` / `pacman` / `zypper` / `apk` / `snap`；
  自动判断 root / `sudo`，在无 TTY 环境（CI、管道）不会挂起等密码，
  而是直接给出可复制的 `sudo ...` 命令。

## 4. 跨平台适配

- 所有配置与输出**固定保存在程序所在目录**，不再写用户主目录，整目录拷贝即可迁移。
  查找优先级：`$DRONEXTRACT_ENV` → `<程序目录>/.env` → `<工作目录>/.env` → 内嵌默认值。
- 默认输出目录改为 `<程序目录>/output/`，导出时自动创建父目录
  （原实现依赖 `../output`，在 Linux 上常因目录不存在而失败）。
- `go run` 场景识别 `go-build` 临时目录并回退到工作目录，不污染系统临时目录。
- 菜单与配置内嵌进二进制（`go:embed`），任意目录运行都有完整界面。
- Shell 脚本统一 LF 行尾，并用 `.gitattributes` 锁定，避免 Linux 报
  `bad interpreter: /bin/sh^M`。
- 尊重 `NO_COLOR` 环境变量。
- 新增 `build.sh`：一次产出 linux/amd64、linux/arm64、darwin/amd64、darwin/arm64、
  windows/amd64 静态二进制，并自动补执行位。
- 新增 `install.sh`：安装到 `~/.local/bin`，按发行版装 ExifTool，在程序目录生成 `.env`。
- 新增 `check_deps.sh`：CI 用的依赖自检，缺依赖返回退出码 2。

## 5. 稳定性修复

修复的原有缺陷（均会导致崩溃或静默错误）：

1. `steganography.ExecuteSteganography` 中 `filePath :=` 遮蔽包级变量，
   用户输入的文件路径被丢弃，SRT/XMP/DNG 入口必然报「文件格式无效」并空指针崩溃。
2. `helpers.Option` 的 `max-1` 规则对无「返回」项的子菜单误判，
   选中第 1 项会直接退回主菜单。
3. SRT 直接导出时从未解析文件，`packets` 为空，导出得到空文件。
4. 空 CSV 触发 `records[0]` 越界崩溃。
5. 空/畸形 SRT 导出时触发 `packets[0]` 越界。
6. 飞行路径地图渲染失败时直接 `panic`（离线环境必然触发），已改为错误提示；
   并在渲染前探测瓦片服务可达性，给出可操作的替代方案。
7. 遥测指标序号越界（`GenerateGraph(0)` → 下标 -1）。
8. 飞行分析在 CSV 列数不足或阈值数量不足时静默算出错误的 0 方差，
   现改为明确报错并提示使用完整 Airdata CSV。
9. 输入流结束（EOF）时菜单无限递归重提示。
10. `txt/` 与 `.env` 只能从当前目录读取，换目录运行菜单空白、采样率退化为 0。

## 6. 测试

- `go test ./...` 全部通过（analysis / deps / parsing / steganography / telemetry）。
- 新增 `deps` 单元测试：依赖项完整性、内嵌配置回退、配置写入、JSON 结构。
- 新增 Linux/macOS 专属测试（`//go:build linux || darwin`）：
  配置路径优先级、输出目录、非 root 安装跳过。
- 新增 `steganography/robustness_test.go`：空 SRT、空文件导出、畸形内容、错误扩展名。
- 修复上游两个长期失败的测试：`telemetry` 传入非法指标序号、`analysis`
  在包内运行读不到根目录 `.env`。

## 7. 待完成 / 已知限制

- 飞行路径地图依赖在线瓦片，离线环境只做降级提示，不生成图片。
- EXIF / DNG / XMP 分析需要外部 ExifTool，未安装时给出安装指引。
- `helpers/errors.go` 中 `Message` 相关方法为上游遗留占位，尚未实现 DUML 解析。
- `SRTAnalysis.go` 中 `ExportToMGJSON` 仍为注释状态。
