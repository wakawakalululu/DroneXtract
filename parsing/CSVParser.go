package parsing

import (
	"encoding/csv"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"os"
	"strconv"
)

type DJI_CSV_Parser struct {
	fileName string
}

func NewDJI_CSV_Parser(fileName string) *DJI_CSV_Parser {
	check := helpers.CheckFileFormat(fileName, ".csv")
	if check == false {
		helpers.PrintError(i18n.S("文件格式无效，必须为 CSV 文件", "INVALID FILE FORMAT. MUST BE CSV FILE"))
		return nil
	}

	parser := DJI_CSV_Parser{
		fileName: fileName,
	}
	return &parser
}

func (parser *DJI_CSV_Parser) ParseContents() {
	file, err := os.Open(parser.fileName)
	if err != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，无法打开", "INVALID FILE. UNABLE TO OPEN"), err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()

	if err != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，无法打开", "INVALID FILE. UNABLE TO OPEN"), err)
		return
	}

	if len(records) == 0 {
		helpers.PrintError(i18n.S("CSV 文件为空，没有可解析的内容", "CSV FILE IS EMPTY, NOTHING TO PARSE"))
		return
	}

	columns := records[0]

	// Print each record
	for count, record := range records {
		row_out := i18n.S("第 "+strconv.Itoa(count)+" 行", "Row "+strconv.Itoa(count))
		GenTableHeader(row_out, count == 0)
		for in, value := range record {
			helpers.GenRowString(columns[in], value)
		}
		helpers.GenTableFooter()
	}
}
