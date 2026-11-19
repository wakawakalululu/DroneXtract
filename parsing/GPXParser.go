package parsing

import (
	"fmt"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"github.com/TwiN/go-color"
	"github.com/tkrajina/gpxgo/gpx"
	"path/filepath"
)

type DJI_GPX_Parser struct {
	fileName string
}

func NewDJI_GPX_Parser(fileName string) *DJI_GPX_Parser {
	check := helpers.CheckFileFormat(fileName, ".gpx")
	if check == false {
		helpers.PrintError(i18n.S("文件格式无效，必须为 GPX 文件", "INVALID FILE FORMAT. MUST BE GPX FILE"))
		return nil
	}

	parser := DJI_GPX_Parser{
		fileName: fileName,
	}
	return &parser
}

func (parser *DJI_GPX_Parser) ParseContents() {
	gpxFile, err := gpx.ParseFile(parser.fileName)

	if err != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，读取内容出错", "INVALID FILE. ERROR READING CONTENTS"), err)
		return
	}

	gpxPath, err2 := filepath.Abs(parser.fileName)

	if err2 != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，读取内容出错", "INVALID FILE. ERROR READING CONTENTS"), err)
		return
	}

	fmt.Println(color.Ize(color.Blue, string(i18n.S("文件：", "File: ")+gpxPath+"\n")))

	fmt.Println(color.Ize(color.Blue, gpxFile.GetGpxInfo()))
}
