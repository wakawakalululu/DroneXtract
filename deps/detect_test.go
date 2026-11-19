package deps

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckAllReturnsAllItems(t *testing.T) {
	r := CheckAll()
	if len(r.Items) == 0 {
		t.Fatal("expected at least one dependency item")
	}
	names := map[string]bool{}
	for _, it := range r.Items {
		if it.Name == "" {
			t.Error("item without name")
		}
		if it.TitleZH == "" || it.TitleEN == "" {
			t.Errorf("%s: missing bilingual title", it.Name)
		}
		names[it.Name] = true
	}
	for _, want := range []string{"config", "exiftool", "go"} {
		if !names[want] {
			t.Errorf("missing item %q", want)
		}
	}
}

func TestLoadEnvFallsBackToBuiltin(t *testing.T) {
	dir := t.TempDir()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)

	for _, key := range []string{"TELEMETRY_VIS_DOWNSAMPLE", "FLIGHT_MAP_DOWNSAMPLE", "ANALYSIS_DOWNSAMPLE", "ANALYSIS_MAX_VARIANCE"} {
		os.Unsetenv(key)
	}

	values, source := LoadEnv()
	if source == ".env" {
		t.Errorf("expected builtin fallback, got %q", source)
	}
	if values["ANALYSIS_MAX_VARIANCE"] == "" {
		t.Error("builtin default for ANALYSIS_MAX_VARIANCE is empty")
	}
	if got := IntValue("ANALYSIS_DOWNSAMPLE", 7); got <= 0 {
		t.Errorf("expected positive downsample, got %d", got)
	}
	if n := len(FloatList("ANALYSIS_MAX_VARIANCE")); n == 0 {
		t.Error("expected parsed variance list")
	}
}

func TestWriteConfigCreatesFile(t *testing.T) {
	dir := t.TempDir()
	// 模拟项目目录，让配置写进临时目录而不是用户配置目录
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module test"), 0o644); err != nil {
		t.Fatal(err)
	}
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(old)
	os.Unsetenv("DRONEXTRACT_ENV")
	// 隔离用户级配置，避免影响/受本机已有配置干扰
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(dir, "xdg"))
	t.Setenv("APPDATA", filepath.Join(dir, "appdata"))

	path, written, err := WriteConfig(false)
	if err != nil {
		t.Fatalf("write config: %v", err)
	}
	if !written {
		t.Fatal("expected file to be written")
	}
	data, err := os.ReadFile(filepath.Join(dir, ".env"))
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	if !strings.Contains(string(data), "ANALYSIS_MAX_VARIANCE") {
		t.Error(".env missing expected keys")
	}

	// 已存在时不应覆盖
	if _, again, _ := WriteConfig(false); again {
		t.Error("expected existing .env to be kept")
	}
	_ = path
}

func TestReportJSONShape(t *testing.T) {
	r := CheckAll()
	out := r.JSON()
	for _, want := range []string{`"os"`, `"arch"`, `"items"`, `"status"`} {
		if !strings.Contains(out, want) {
			t.Errorf("json missing %s", want)
		}
	}
}
