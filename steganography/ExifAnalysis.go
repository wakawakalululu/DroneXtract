package steganography

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/ANG13T/DroneXtract/deps"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"github.com/barasher/go-exiftool"
	"os"
	"path/filepath"
	"strings"
)

type DJI_EXIF_Parser struct {
	fileName string
}

func NewDJI_EXIF_Parser(fileName string) *DJI_EXIF_Parser {
	parser := DJI_EXIF_Parser{
		fileName: fileName,
	}
	return &parser
}

func (parser *DJI_EXIF_Parser) ExecuteEXIFAnalysis(input int) {
	switch in := input; in {
	case 1:
		parser.Read()
	case 2:
		outputPath := helpers.FileInputString()
		parser.ExportToTXT(outputPath)
	case 3:
		outputPath := helpers.FileInputString()
		parser.ExportToCSV(outputPath)
	case 4:
		outputPath := helpers.FileInputString()
		parser.ExportToJSON(outputPath)
	}
}

// requireExifTool 在使用前确认 ExifTool 可用，避免直接抛出底层错误。
func requireExifTool() bool {
	if deps.CheckExifTool().Ready() {
		return true
	}
	helpers.PrintError(i18n.S(
		"未检测到 ExifTool：请在主菜单选择 [ 5 ] 环境依赖检查并自动安装，或访问 https://exiftool.org/install.html",
		"EXIFTOOL NOT FOUND. CHOOSE [ 5 ] IN THE MAIN MENU TO INSTALL, OR VISIT https://exiftool.org/install.html"))
	return false
}

func (parser *DJI_EXIF_Parser) Read() {
	if !requireExifTool() {
		return
	}
	et, err := exiftool.NewExiftool()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			helpers.PrintError(i18n.S("未检测到 ExifTool：请在主菜单选择 [ 5 ] 环境依赖检查并自动安装，或访问 https://exiftool.org/install.html", "EXIFTOOL NOT FOUND. CHOOSE [ 5 ] IN THE MAIN MENU TO INSTALL, OR VISIT https://exiftool.org/install.html"))
		} else {
			helpers.PrintErrorLog(i18n.S("无法初始化 EXIF Tool", "COULD NOT INITIALIZE EXIF TOOL"), err)
		}
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(parser.fileName)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			helpers.PrintErrorLog(i18n.S("无法读取文件", "COULD NOT READ FILE"), fileInfo.Err)
			continue
		}

		helpers.GenTableHeader(i18n.S("EXIF 分析", "EXIF Analysis"))
		for k, v := range fileInfo.Fields {
			formattedValue := helpers.FormatSingleLine(fmt.Sprintf("%v", v))
			helpers.GenRowString(k, formattedValue)
		}
		helpers.GenTableFooter()
	}
}

func (parser *DJI_EXIF_Parser) ExportToTXT(outputPath string) {
	if len(outputPath) == 0 {
		outputPath = filepath.Join(deps.OutputDir(), "exif-analysis.txt")
	}

	check := helpers.CheckFileFormat(outputPath, ".txt")
	if check == false {
		helpers.PrintError(i18n.S("输出文件格式无效，必须为 TXT 文件", "INVALID OUTPUT FILE FORMAT. MUST BE TXT FILE"))
		return
	}

	if _, err := deps.EnsureParent(outputPath); err != nil {
		helpers.PrintErrorLog(i18n.S("无法创建输出目录", "CANNOT CREATE OUTPUT DIRECTORY"), err)
		return
	}
	file, err := os.Create(outputPath)
	if err != nil {
		helpers.PrintErrorLog(i18n.S("创建 TXT 文件失败", "FAILED TO CREATE TXT FILE"), err)
		return
	}
	defer file.Close()

	if !requireExifTool() {
		return
	}
	et, err := exiftool.NewExiftool()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			helpers.PrintError(i18n.S("未检测到 ExifTool：请在主菜单选择 [ 5 ] 环境依赖检查并自动安装，或访问 https://exiftool.org/install.html", "EXIFTOOL NOT FOUND. CHOOSE [ 5 ] IN THE MAIN MENU TO INSTALL, OR VISIT https://exiftool.org/install.html"))
		} else {
			helpers.PrintErrorLog(i18n.S("无法初始化 EXIF Tool", "COULD NOT INITIALIZE EXIF TOOL"), err)
		}
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(parser.fileName)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			helpers.PrintErrorLog(i18n.S("无法读取文件", "COULD NOT READ FILE"), fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			formattedValue := helpers.FormatSingleLine(fmt.Sprintf("%v", v))

			// Write the output to the file
			output := k + ":" + formattedValue + "\n"
			_, err = file.WriteString(output)

			if err != nil {
				helpers.PrintErrorLog(i18n.S("写入 TXT 文件失败", "FAILED TO WRITE TO TXT FILE"), err)
				return
			}
		}
	}

}

func (parser *DJI_EXIF_Parser) ExportToCSV(outputPath string) {
	if len(outputPath) == 0 {
		outputPath = filepath.Join(deps.OutputDir(), "exif-analysis.csv")
	}

	check := helpers.CheckFileFormat(outputPath, ".csv")
	if check == false {
		helpers.PrintError(i18n.S("输出文件格式无效，必须为 CSV 文件", "INVALID OUTPUT FILE FORMAT. MUST BE CSV FILE"))
		return
	}

	if _, err := deps.EnsureParent(outputPath); err != nil {
		helpers.PrintErrorLog(i18n.S("无法创建输出目录", "CANNOT CREATE OUTPUT DIRECTORY"), err)
		return
	}
	file, err := os.Create(outputPath)
	if err != nil {
		helpers.PrintErrorLog(i18n.S("创建 CSV 文件失败", "FAILED TO CREATE CSV FILE"), err)
		return
	}
	defer file.Close()

	if !requireExifTool() {
		return
	}
	et, err := exiftool.NewExiftool()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			helpers.PrintError(i18n.S("未检测到 ExifTool：请在主菜单选择 [ 5 ] 环境依赖检查并自动安装，或访问 https://exiftool.org/install.html", "EXIFTOOL NOT FOUND. CHOOSE [ 5 ] IN THE MAIN MENU TO INSTALL, OR VISIT https://exiftool.org/install.html"))
		} else {
			helpers.PrintErrorLog(i18n.S("无法初始化 EXIF Tool", "COULD NOT INITIALIZE EXIF TOOL"), err)
		}
		return
	}
	defer et.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	fileInfos := et.ExtractMetadata(parser.fileName)

	for _, fileInfo := range fileInfos {
		if fileInfo.Err != nil {
			helpers.PrintErrorLog(i18n.S("无法读取文件", "COULD NOT READ FILE"), fileInfo.Err)
			continue
		}

		for k, v := range fileInfo.Fields {
			formattedValue := helpers.FormatSingleLine(fmt.Sprintf("%v", v))

			// Write the output to the file
			output := []string{k, formattedValue}
			err := writer.Write(output)

			if err != nil {
				helpers.PrintErrorLog(i18n.S("写入 CSV 文件失败", "FAILED TO WRITE TO CSV FILE"), err)
				return
			}
		}
	}
}

func (parser *DJI_EXIF_Parser) ExportToJSON(outputPath string) {
	if len(outputPath) == 0 {
		outputPath = filepath.Join(deps.OutputDir(), "exif-analysis.json")
	}

	check := helpers.CheckFileFormat(outputPath, ".json")
	if check == false {
		helpers.PrintError(i18n.S("输出文件格式无效，必须为 JSON 文件", "INVALID OUTPUT FILE FORMAT. MUST BE JSON FILE"))
		return
	}

	if _, err := deps.EnsureParent(outputPath); err != nil {
		helpers.PrintErrorLog(i18n.S("无法创建输出目录", "CANNOT CREATE OUTPUT DIRECTORY"), err)
		return
	}
	file, err := os.Create(outputPath)
	if err != nil {
		helpers.PrintErrorLog(i18n.S("创建 JSON 文件失败", "FAILED TO CREATE JSON FILE"), err)
		return
	}
	defer file.Close()

	if !requireExifTool() {
		return
	}
	et, err := exiftool.NewExiftool()
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found") {
			helpers.PrintError(i18n.S("未检测到 ExifTool：请在主菜单选择 [ 5 ] 环境依赖检查并自动安装，或访问 https://exiftool.org/install.html", "EXIFTOOL NOT FOUND. CHOOSE [ 5 ] IN THE MAIN MENU TO INSTALL, OR VISIT https://exiftool.org/install.html"))
		} else {
			helpers.PrintErrorLog(i18n.S("无法初始化 EXIF Tool", "COULD NOT INITIALIZE EXIF TOOL"), err)
		}
		return
	}
	defer et.Close()

	fileInfos := et.ExtractMetadata(parser.fileName)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(fileInfos[0])

	if err != nil {
		helpers.PrintErrorLog(i18n.S("JSON 编码失败", "FAILED TO ENCODE JSON"), err)
		return
	}
}
