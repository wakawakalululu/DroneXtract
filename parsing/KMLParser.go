package parsing

import (
	"encoding/xml"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"io/ioutil"
	"strconv"
	"strings"
)

type DJI_KML_Parser struct {
	fileName string
}

type KML struct {
	Placemarks []Placemark `xml:"Document>Placemark"`
}

type Placemark struct {
	Name       string     `xml:"name"`
	Point      Point      `xml:"Point"`
	LineString LineString `xml:"LineString"`
}

type ExtendedData struct {
	Data []Data `xml:"Data"`
}

type Data struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

type Point struct {
	Coordinates string `xml:"coordinates"`
}

type LineString struct {
	Coordinates string `xml:"coordinates"`
}

func NewDJI_KML_Parser(fileName string) *DJI_KML_Parser {
	check := helpers.CheckFileFormat(fileName, ".kml")
	if check == false {
		helpers.PrintError(i18n.S("文件格式无效，必须为 KML 文件", "INVALID FILE FORMAT. MUST BE KML FILE"))
		return nil
	}

	parser := DJI_KML_Parser{
		fileName: fileName,
	}
	return &parser
}

func (parser *DJI_KML_Parser) ParseContents() {
	content, err := ioutil.ReadFile(parser.fileName)
	if err != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，读取内容出错", "INVALID FILE. ERROR READING CONTENTS"), err)
	}

	// Parse the KML data
	var kml KML
	err = xml.Unmarshal(content, &kml)
	if err != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，读取内容出错", "INVALID FILE. ERROR READING CONTENTS"), err)
	}

	// Process the placemarks
	for _, placemark := range kml.Placemarks {

		if placemark.Point.Coordinates != "" {
			coorValues := strings.Split(placemark.Point.Coordinates, ",")
			GenTableHeader(i18n.S("返航点信息", "Home Point Information"), true)
			helpers.GenRowString(i18n.S("坐标", "Coordinates"), "("+coorValues[0]+","+coorValues[1]+")")
			helpers.GenRowString("高度", coorValues[2]+i18n.S(" 英尺", " ft"))
			helpers.GenTableFooter()
		}

		if placemark.LineString.Coordinates != "" {

			lines := strings.Split(placemark.LineString.Coordinates, "\n")

			for cIndex, coor := range lines {
				coorValues := strings.Split(coor, ",")
				if len(coor) > 0 && len(coorValues) > 1 {
					GenTableHeader(i18n.S("坐标点 ", "Coordinate Point ")+strconv.Itoa(cIndex+1), false)
					helpers.GenRowString(i18n.S("坐标", "Coordinates"), "("+coorValues[0]+", "+coorValues[1]+")")
					helpers.GenRowString("高度", coorValues[2]+i18n.S(" 英尺", " ft"))
					helpers.GenTableFooter()
				}
			}
		}
	}
}
