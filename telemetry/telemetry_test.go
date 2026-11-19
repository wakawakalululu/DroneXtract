package telemetry

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/TwiN/go-color"
)

// TestMain 让包内测试也能读到项目根目录的配置，
// 否则 from-package 运行时采样率会退化为 0 导致图表无数据。
func TestMain(m *testing.M) {
	if _, err := os.Stat("../.env"); err == nil {
		if err := os.Setenv("DRONEXTRACT_ENV", "../.env"); err != nil {
			panic(err)
		}
	}
	os.Exit(m.Run())
}

func TestTelemetryVisualizations(t *testing.T) {
	t.Logf("Running Telemetry Visualization Test Unit")
	fmt.Println(color.Ize(color.Cyan, "Running Telemetry Visualization Test Unit"))

	filename := `../test-data/Airdata-Files/AirdataCSV.csv`

	suite := NewDJI_TelemetryVisualizations(filename)
	if suite == nil {
		t.Fatal("failed to create telemetry suite")
	}

	// 菜单传入 1..N，1 表示第一个指标（起飞后高度）
	suite.GenerateGraph(1)

	t.Logf("[Telemetry Visualizations #1] Graph Display Case - PASS")
}

func TestFlightPathMap(t *testing.T) {
	t.Logf("Running Flight Path Test Unit")
	fmt.Println(color.Ize(color.Cyan, "Running Flight Path Map Test Unit"))

	filename := `../test-data/Airdata-Files/AirdataCSV.csv`

	suite := NewDJI_Flight_Path_Map(filename, filepath.Join(t.TempDir(), "flight-map.png"))

	suite.ExecuteFlightPathAnalysis()

	t.Logf("[Telemetry Visualizations #1] Flight Map Case - PASS")
}
