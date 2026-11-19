package steganography

import (
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
)

type DJI_XMP_Parser struct {
	fileName string
}

func NewDJI_XMP_Parser(fileName string) *DJI_XMP_Parser {
	check := helpers.CheckFileFormat(fileName, ".xmp")
	if check == false {
		helpers.PrintError(i18n.S("文件格式无效，必须为 XMP 文件", "INVALID FILE FORMAT. MUST BE XMP FILE"))
		return nil
	}

	parser := DJI_XMP_Parser{
		fileName: fileName,
	}
	return &parser
}

func (parser *DJI_XMP_Parser) ExecuteXMPAnalysis(input int) {
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

func (parser *DJI_XMP_Parser) Read() {
	exif := NewDJI_EXIF_Parser(parser.fileName)
	exif.Read()
}

func (parser *DJI_XMP_Parser) ExportToTXT(outputPath string) {
	exif := NewDJI_EXIF_Parser(parser.fileName)
	exif.ExportToTXT(outputPath)
}

func (parser *DJI_XMP_Parser) ExportToCSV(outputPath string) {
	exif := NewDJI_EXIF_Parser(parser.fileName)
	exif.ExportToCSV(outputPath)
}

func (parser *DJI_XMP_Parser) ExportToJSON(outputPath string) {
	exif := NewDJI_EXIF_Parser(parser.fileName)
	exif.ExportToJSON(outputPath)
}
