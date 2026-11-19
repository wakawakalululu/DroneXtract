package analysis

import (
	"encoding/csv"
	"fmt"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"os"
	"strconv"
)

var downsample_number_analysis = helpers.GetEnvVariable("ANALYSIS_DOWNSAMPLE")

var cooresponding_index = []int{8, 9, 10, 11, 14, 18, 19, 20, 22, 23, 24, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43}
var indicators = []string{"height_sonar(feet)", "speed(mph)", "distance(feet)", "mileage(feet)", "voltage(v)", "xSpeed(mph)", "ySpeed(mph)", "zSpeed(mph)", "compass_heading(degrees)", "pitch(degrees)", "roll(degrees)", "rc_elevator", "rc_aileron", "rc_throttle", "rc_rudder", "rc_elevator(percent)", "rc_aileron(percent)", "rc_throttle(percent)", "rc_rudder(percent)", "gimbal_heading(degrees)", "gimbal_pitch(degrees)", "gimbal_roll(degrees)", "battery_percent", "current(A)", "battery_temperature(f)", "altitude(feet)", "ascent(feet)"}

type labelPair struct{ zh, en string }

var labels = []labelPair{
	{"起飞后高度（feet）", "Height Sonar (feet)"},
	{"速度（mph）", "Speed (mph)"},
	{"距离（feet）", "Distance (feet)"},
	{"里程（feet）", "Mileage (feet)"},
	{"电压（V）", "Voltage (V)"},
	{"X 轴速度（mph）", "X Speed (mph)"},
	{"Y 轴速度（mph）", "Y Speed (mph)"},
	{"Z 轴速度（mph）", "Z Speed (mph)"},
	{"罗盘航向（degrees）", "Compass Heading (degrees)"},
	{"俯仰角（degrees）", "Pitch (degrees)"},
	{"横滚角（degrees）", "Roll (degrees)"},
	{"遥控器升降通道", "RC Elevator"},
	{"遥控器横滚通道", "RC Aileron"},
	{"遥控器油门通道", "RC Throttle"},
	{"遥控器偏航通道", "RC Rudder"},
	{"遥控器升降通道（%）", "RC Elevator (percent)"},
	{"遥控器横滚通道（%）", "RC Aileron (percent)"},
	{"遥控器油门通道（%）", "RC Throttle (percent)"},
	{"遥控器偏航通道（%）", "RC Rudder (percent)"},
	{"云台航向（degrees）", "Gimbal Heading (degrees)"},
	{"云台俯仰角（degrees）", "Gimbal Pitch (degrees)"},
	{"云台横滚角（degrees）", "Gimbal Roll (degrees)"},
	{"电池电量（%）", "Battery Percent"},
	{"电流（A）", "Current (A)"},
	{"电池温度（F）", "Battery Temperature (F)"},
	{"海拔（feet）", "Altitude (feet)"},
	{"爬升（feet）", "Ascent (feet)"},
}

// labelAt 返回第 i 个指标的当前语言名称。
func labelAt(i int) string {
	return i18n.S(labels[i].zh, labels[i].en)
}

var max_variance = helpers.GetEnvVariances()

type DJI_Analysis struct {
	fileName   string
	outputPath string
}

func NewDJI_Analysis(fileName string) *DJI_Analysis {
	parser := DJI_Analysis{
		fileName: fileName,
	}
	return &parser
}

func ExecuteAnalysis() {
	fileName := helpers.FileInputString()

	check := helpers.CheckFileFormat(fileName, ".csv")
	if check == false {
		helpers.PrintError(i18n.S("文件格式无效，必须为 CSV 文件", "INVALID FILE FORMAT. MUST BE CSV FILE"))
		return
	}

	suite := NewDJI_Analysis(fileName)
	suite.RunAnalysis()
}

func (parser *DJI_Analysis) RunAnalysis() {
	file, err := os.Open(parser.fileName)
	if err != nil {
		helpers.PrintErrorLog(i18n.S(i18n.S("文件无效，无法打开", "INVALID FILE. UNABLE TO OPEN"), "INVALID FILE. UNABLE TO OPEN"), err)
		return
	}
	defer file.Close()

	// Create a new CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	records, err := reader.ReadAll()

	if err != nil {
		helpers.PrintErrorLog(i18n.S("文件无效，无法读取内容", "INVALID FILE. CANNOT READ CONTENTS"), err)
		return
	}

	if len(records) == 0 {
		helpers.PrintError(i18n.S("CSV 文件为空，无法进行飞行分析", "CSV FILE IS EMPTY, CANNOT RUN ANALYSIS"))
		return
	}

	// 列数不足时逐项统计会越界，先给出明确提示
	maxCol := 0
	for _, idx := range cooresponding_index {
		if idx > maxCol {
			maxCol = idx
		}
	}
	if len(records[0]) <= maxCol {
		helpers.PrintError(fmt.Sprintf(i18n.S(
			"CSV 列数不足：需要至少 %d 列，当前只有 %d 列。请使用 Airdata 导出的完整 CSV。",
			"NOT ENOUGH COLUMNS: NEED AT LEAST %d, GOT %d. USE A FULL AIRDATA CSV EXPORT."),
			maxCol+1, len(records[0])))
		return
	}

	if len(max_variance) < len(indicators) {
		helpers.PrintError(fmt.Sprintf(i18n.S(
			"方差阈值数量不足：需要 %d 项，当前 %d 项。请检查 .env 中的 ANALYSIS_MAX_VARIANCE。",
			"NOT ENOUGH VARIANCE THRESHOLDS: NEED %d, GOT %d. CHECK ANALYSIS_MAX_VARIANCE IN .env."),
			len(indicators), len(max_variance)))
		return
	}

	helpers.PrintLog(i18n.S("==== 开始飞行分析 ====", "==== RUNNING FLIGHT ANALYSIS ===="))

	invalid_count := 0
	valid_count := 0

	for index, _ := range indicators {
		variance := parser.GetVariance(index, records)
		valid := variance <= float64(max_variance[index])
		if valid {
			helpers.PrintValidLog("[√] " + labelAt(index) + " > " + fmt.Sprintf("%.2f", variance))
			valid_count++
		} else {
			helpers.PrintInvalidLog("[X] " + labelAt(index) + " > " + fmt.Sprintf("%.2f", variance) + " " + fmt.Sprintf("%.2f", max_variance[index]))
			invalid_count++
		}
	}

	helpers.PrintLog(i18n.S("==== 飞行分析结束 ====", "==== FINISHED FLIGHT ANALYSIS ===="))
	helpers.PrintValidLog("[√] " + i18n.S(strconv.Itoa(valid_count)+" 项方差正常", strconv.Itoa(valid_count)+" VALID VARIANCE"))
	helpers.PrintInvalidLog("[X] " + i18n.S(strconv.Itoa(invalid_count)+" 项方差异常", strconv.Itoa(invalid_count)+" INVALID VARIANCE"))
	helpers.PrintLog("==================================")
}

func GetMax(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0.0
	}

	maxValue := numbers[0]
	for _, num := range numbers {
		if num > maxValue {
			maxValue = num
		}
	}

	return maxValue
}

func (parser *DJI_Analysis) GetVariance(value int, records [][]string) float64 {
	result := parser.GetCSVValues(value, records)

	max := GetMax(result)

	min := GetMin(result)

	variance := max - min

	return variance
}

func GetMin(numbers []float64) float64 {
	if len(numbers) == 0 {
		return 0.0
	}

	minValue := numbers[0]
	for _, num := range numbers {
		if num < minValue {
			minValue = num
		}
	}

	return minValue
}

func (parser *DJI_Analysis) GetCSVValues(index int, records [][]string) []float64 {
	output := []float64{}

	// Print each record
	for _, record := range records {
		for in, value := range record {
			if in == cooresponding_index[index] {
				val, _ := strconv.ParseFloat(value, 64)
				output = append(output, val)
			}
		}
	}

	if len(output) > downsample_number_analysis {
		output = downsampleArray(output, downsample_number_analysis)
	}

	return output
}

func downsampleArray(data []float64, targetLength int) []float64 {
	length := len(data)
	if targetLength >= length {
		return data
	}

	ratio := float64(length) / float64(targetLength)
	result := make([]float64, targetLength)
	resultIndex := 0

	for i := 0; i < targetLength; i++ {
		rangeStart := int(float64(i) * ratio)
		rangeEnd := int(float64(i+1) * ratio)

		// Calculate the average within the range
		sum := 0.0
		for j := rangeStart; j < rangeEnd; j++ {
			sum += data[j]
		}
		average := sum / float64(rangeEnd-rangeStart)

		result[resultIndex] = average
		resultIndex++
	}

	result[0] = data[0]

	return result
}
