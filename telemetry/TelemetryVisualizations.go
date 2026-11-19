package telemetry

import (
	"encoding/csv"
	"fmt"
	"github.com/ANG13T/DroneXtract/helpers"
	"github.com/ANG13T/DroneXtract/i18n"
	"github.com/TwiN/go-color"
	"github.com/guptarohit/asciigraph"
	"os"
	"strconv"
)

var downsample_number_telemetry = helpers.GetEnvVariable("TELEMETRY_VIS_DOWNSAMPLE")

type DJI_TelemetryVisualizations struct {
	fileName string
}

func NewDJI_TelemetryVisualizations(fileName string) *DJI_TelemetryVisualizations {
	check := helpers.CheckFileFormat(fileName, ".csv")
	if check == false {
		helpers.PrintError(i18n.S("文件格式无效，必须为 CSV 文件", "INVALID FILE FORMAT. MUST BE CSV FILE"))
		return nil
	}

	parser := DJI_TelemetryVisualizations{
		fileName: fileName,
	}
	return &parser
}

var indicators = []string{"height_above_takeoff(feet)", "height_above_ground_at_drone_location(feet)", "ground_elevation_at_drone_location(feet)", "ground_elevation_at_drone_location(feet)", "altitude_above_seaLevel(feet)", "height_sonar(feet)", "speed(mph)", "distance(feet)", "mileage(feet)", "satellites", "gpslevel", "voltage(v)", "max_altitude(feet)", "max_ascent(feet)", "max_speed(mph)", "max_distance(feet)", "xSpeed(mph)", "ySpeed(mph)", "zSpeed(mph)", "compass_heading(degrees)", "pitch(degrees)", "roll(degrees)", "rc_elevator", "rc_aileron", "rc_throttle", "rc_rudder", "rc_elevator(percent)", "rc_aileron(percent)", "rc_throttle(percent)", "rc_rudder(percent)", "gimbal_heading(degrees)", "gimbal_pitch(degrees)", "gimbal_roll(degrees)", "battery_percent", "voltageCell1", "voltageCell2", "voltageCell3", "voltageCell4", "voltageCell5", "voltageCell6", "current(A)", "battery_temperature(f)", "altitude(feet)", "ascent(feet)", "flycStateRaw"}

type labelPair struct{ zh, en string }

var labels = []labelPair{
	{"起飞后高度（英尺）", "Height Above Takeoff (feet)"},
	{"无人机所在位置离地高度（英尺）", "Height Above Ground at Drone Location (feet)"},
	{"无人机所在位置地面海拔（英尺）", "Ground Elevation at Drone Location (feet)"},
	{"无人机所在位置地面海拔（英尺）", "Ground Elevation at Drone Location (feet)"},
	{"海拔高度（英尺）", "Altitude Above Sea Level (feet)"},
	{"声呐测高（英尺）", "Height Sonar (feet)"},
	{"速度（英里/小时）", "Speed (mph)"},
	{"距离（英尺）", "Distance (feet)"},
	{"里程（英尺）", "Mileage (feet)"},
	{"卫星数", "Satellites"},
	{"GPS 信号等级", "GPS Level"},
	{"电压（伏特）", "Voltage (V)"},
	{"最大高度（英尺）", "Max Altitude (feet)"},
	{"最大爬升（英尺）", "Max Ascent (feet)"},
	{"最大速度（英里/小时）", "Max Speed (mph)"},
	{"最大距离（英尺）", "Max Distance (feet)"},
	{"X 轴速度（英里/小时）", "X Speed (mph)"},
	{"Y 轴速度（英里/小时）", "Y Speed (mph)"},
	{"Z 轴速度（英里/小时）", "Z Speed (mph)"},
	{"罗盘航向（度）", "Compass Heading (degrees)"},
	{"俯仰角（度）", "Pitch (degrees)"},
	{"横滚角（度）", "Roll (degrees)"},
	{"遥控器升降通道", "RC Elevator"},
	{"遥控器横滚通道", "RC Aileron"},
	{"遥控器油门通道", "RC Throttle"},
	{"遥控器偏航通道", "RC Rudder"},
	{"遥控器升降通道（百分比）", "RC Elevator (percent)"},
	{"遥控器横滚通道（百分比）", "RC Aileron (percent)"},
	{"遥控器油门通道（百分比）", "RC Throttle (percent)"},
	{"遥控器偏航通道（百分比）", "RC Rudder (percent)"},
	{"云台航向（度）", "Gimbal Heading (degrees)"},
	{"云台俯仰角（度）", "Gimbal Pitch (degrees)"},
	{"云台横滚角（度）", "Gimbal Roll (degrees)"},
	{"电池电量（%）", "Battery Percent"},
	{"电芯 1 电压", "Voltage Cell 1"},
	{"电芯 2 电压", "Voltage Cell 2"},
	{"电芯 3 电压", "Voltage Cell 3"},
	{"电芯 4 电压", "Voltage Cell 4"},
	{"电芯 5 电压", "Voltage Cell 5"},
	{"电芯 6 电压", "Voltage Cell 6"},
	{"电流（安培）", "Current (A)"},
	{"电池温度（华氏度）", "Battery Temperature (F)"},
	{"海拔（英尺）", "Altitude (feet)"},
	{"爬升（英尺）", "Ascent (feet)"},
	{"飞控状态原始值", "Flyc State Raw"},
}

// labelAt 返回第 i 个指标的当前语言名称。
func labelAt(i int) string {
	return i18n.S(labels[i].zh, labels[i].en)
}
func (parser *DJI_TelemetryVisualizations) ExecuteTelemetryVisualizations() {
	value := GenerateOptions()
	if value == -1 {
		return
	}
	parser.GenerateGraph(value)
}

func (parser *DJI_TelemetryVisualizations) GenerateGraph(index int) {
	index = index - 1
	// 菜单传入的是 1..N，任何越界值都应被拒绝而不是让程序崩溃
	if index < 0 || index >= len(labels) {
		helpers.PrintError(fmt.Sprintf(i18n.S(
			"指标序号超出范围：有效范围 1-%d，收到 %d", "INDEX OUT OF RANGE: VALID 1-%d, GOT %d"),
			len(labels), index+1))
		return
	}
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
		helpers.PrintErrorLog(i18n.S("文件无效，无法打开", "INVALID FILE. UNABLE TO OPEN"), err)
		return
	}

	output := []float64{}

	// Print each record
	for _, record := range records {
		for in, value := range record {
			if in == index {
				val, _ := strconv.ParseFloat(value, 64)
				output = append(output, val)
			}
		}
	}

	if len(output) > downsample_number_telemetry {
		output = downsampleArray(output, downsample_number_telemetry)
	}

	graph := asciigraph.Plot(output, asciigraph.Height(10), asciigraph.Width(100), asciigraph.Caption(labelAt(index)))

	fmt.Println(color.Ize(color.Cyan, graph))
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

func GenerateOptions() int {
	helpers.GenTableHeader(i18n.S("选择要可视化的指标", "Select Value to Visualize"))
	for index := range labels {
		helpers.GenRowString(strconv.Itoa(index+1), labelAt(index))
	}
	helpers.GenRowString(strconv.Itoa(len(labels)+1), i18n.BackLabel())
	helpers.GenRowString("0", i18n.ExitLabel())
	helpers.GenTableFooter()
	return helpers.Option(0, len(labels))
}
