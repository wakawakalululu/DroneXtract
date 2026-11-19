package telemetry

import (
	"encoding/csv"
	"github.com/ANG13T/DroneXtract/deps"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"os"
	"strconv"

	"image/color"

	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
)

var downsample_number_map = helpers.GetEnvVariable("FLIGHT_MAP_DOWNSAMPLE")

type Coordinate struct {
	Latitude  float64
	Longitude float64
}

type DJI_Flight_Path_Map struct {
	fileName   string
	outputPath string
}

func NewDJI_Flight_Path_Map(fileName string, outputPath string) *DJI_Flight_Path_Map {
	parser := DJI_Flight_Path_Map{
		fileName:   fileName,
		outputPath: outputPath,
	}
	return &parser
}

func (parser *DJI_Flight_Path_Map) ExecuteFlightPathAnalysis() {
	parser.PrintGPSCoordinates()
}

func (parser *DJI_Flight_Path_Map) PrintGPSCoordinates() {
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
		helpers.PrintError(i18n.S("CSV 文件为空，无法提取坐标", "CSV FILE IS EMPTY, NO COORDINATES"))
		return
	}

	columns := records[0]

	coors := []Coordinate{}

	// Print each record
	for _, record := range records {

		lat_val := 0.0
		lon_val := 0.0

		for in, value := range record {
			val, _ := strconv.ParseFloat(value, 64)

			if columns[in] == "latitude" && val != 0 {
				lat_val = val
			}

			if columns[in] == "longitude" && val != 0 {
				lon_val = val
			}
		}

		if lat_val != 0.0 && lon_val != 0.0 {
			coor_val := Coordinate{Latitude: lat_val, Longitude: lon_val}
			coors = append(coors, coor_val)
		}
	}

	GenerateMapOutput(coors, parser.outputPath)
}

func GenerateMapOutput(coors []Coordinate, outputPath string) {
	if len(outputPath) == 0 {
		outputPath = "flight-path-map.png"
	}

	if len(coors) > downsample_number_map {
		coors = downsampleCoordinates(coors, downsample_number_map)
	}

	check := helpers.CheckFileFormat(outputPath, ".png")
	if check == false {
		helpers.PrintError(i18n.S("输出文件格式无效，必须为 PNG 文件", "INVALID OUTPUT FILE FORMAT. MUST BE PNG FILE"))
		return
	}
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)

	for index, coor := range coors {
		ctx.AddObject(
			sm.NewMarker(
				s2.LatLngFromDegrees(coor.Latitude, coor.Longitude),
				color.RGBA{0xff, 0, 0, 0xff},
				16.0,
			),
		)

		if index < len(coors)-1 {
			point1 := s2.LatLngFromDegrees(coors[index].Latitude, coors[index].Longitude)
			point2 := s2.LatLngFromDegrees(coors[index+1].Latitude, coors[index+1].Longitude)

			pos := []s2.LatLng{point1, point2}

			path := sm.NewPath(pos, color.RGBA{0xff, 0, 0, 0xff}, 6.0)

			ctx.AddObject(path)
		}
	}

	PrintCoordinates(coors, len(coors) <= 10)

	if !deps.CheckTileServer() {
		helpers.PrintError(i18n.S(
			"无法连接在线地图服务：本功能需要联网加载瓦片。坐标已在上方输出，可改用 SRT 导出 GEOJSON 做离线分析。",
			"CANNOT REACH MAP TILE SERVICE: INTERNET REQUIRED. COORDINATES PRINTED ABOVE; USE SRT GEOJSON EXPORT FOR OFFLINE ANALYSIS."))
		return
	}

	img, err := ctx.Render()
	if err != nil {
		helpers.PrintError(i18n.S(
			"生成地图失败：无法加载在线地图瓦片，请检查网络后重试（坐标表已在上方打印）",
			"FAILED TO RENDER MAP: CANNOT LOAD MAP TILES. CHECK NETWORK AND RETRY (COORDINATES ABOVE) "))
		helpers.PrintErrorLog("", err)
		return
	}

	if _, err := deps.EnsureParent(outputPath); err != nil {
		helpers.PrintErrorLog(i18n.S("无法创建输出目录", "CANNOT CREATE OUTPUT DIRECTORY"), err)
		return
	}
	if err := gg.SavePNG(outputPath, img); err != nil {
		helpers.PrintErrorLog(i18n.S("保存 PNG 失败", "FAILED TO SAVE PNG"), err)
		return
	}

	helpers.PrintLog(i18n.S("飞行路径地图已生成：", "Created Flight Path Map at ") + outputPath)
}

func PrintCoordinates(coordinates []Coordinate, downsampled bool) {
	if downsampled {
		helpers.GenTableHeader(i18n.S("降采样后的 GPS 坐标", "Downsampled GPS Coordinates"))
	} else {
		helpers.GenTableHeader(i18n.S("GPS 坐标", "GPS Coordinates"))
	}

	for in, coor := range coordinates {
		str_lat := strconv.FormatFloat(coor.Latitude, 'f', -1, 64)
		str_lon := strconv.FormatFloat(coor.Longitude, 'f', -1, 64)
		helpers.GenRowString(i18n.S("坐标 ", "Coordinate ")+strconv.Itoa(in+1), "("+str_lat+", "+str_lon+")")
	}

	helpers.GenTableFooter()
}

func downsampleCoordinates(coordinates []Coordinate, targetLength int) []Coordinate {
	length := len(coordinates)
	if targetLength >= length {
		return coordinates
	}

	ratio := float64(length) / float64(targetLength)
	result := make([]Coordinate, targetLength)
	resultIndex := 0

	for i := 0; i < targetLength; i++ {
		rangeStart := int(float64(i) * ratio)
		rangeEnd := int(float64(i+1) * ratio)

		// Calculate the average of latitude and longitude within the range
		sumLat := 0.0
		sumLon := 0.0
		for j := rangeStart; j < rangeEnd; j++ {
			sumLat += coordinates[j].Latitude
			sumLon += coordinates[j].Longitude
		}
		averageLat := sumLat / float64(rangeEnd-rangeStart)
		averageLon := sumLon / float64(rangeEnd-rangeStart)

		result[resultIndex] = Coordinate{
			Latitude:  averageLat,
			Longitude: averageLon,
		}
		resultIndex++
	}

	return result
}
