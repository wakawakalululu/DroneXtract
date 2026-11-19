package steganography

import (
	"os"
	"path/filepath"
	"testing"
)

// TestEmptySRTDoesNotPanic 空文件应给出提示而不是崩溃。
func TestEmptySRTDoesNotPanic(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.srt")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	suite := NewDJI_SRT_Parser(path)
	if suite == nil {
		t.Fatal("parser should be created for valid extension")
	}

	// 不应 panic
	suite.PrintFileMetadata()
	suite.GeneratePackets()
	suite.PrintAllPackets()

	if len(suite.packets) != 0 {
		t.Errorf("expected no packets, got %d", len(suite.packets))
	}
}

// TestEmptySRTExportDoesNotPanic 空文件导出应安全返回，不越界访问 packets[0]。
func TestEmptySRTExportDoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.srt")
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	suite := NewDJI_SRT_Parser(path)
	if suite == nil {
		t.Fatal("parser should be created")
	}

	out := filepath.Join(dir, "out.json")
	suite.ExportToJSON(out)
	suite.ExportToGeoJSON(filepath.Join(dir, "out.geojson"))
	suite.ExportToCSV(filepath.Join(dir, "out.csv"))
}

// TestGarbageSRTDoesNotPanic 畸形内容应被拒绝而不是崩溃。
func TestGarbageSRTDoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.srt")
	if err := os.WriteFile(path, []byte("not an srt at all\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	suite := NewDJI_SRT_Parser(path)
	if suite == nil {
		t.Fatal("parser should be created")
	}
	suite.GeneratePackets()
	suite.PrintAllPackets()
}

// TestInvalidExtensionRejected 扩展名不匹配时构造函数返回 nil。
func TestInvalidExtensionRejected(t *testing.T) {
	if p := NewDJI_SRT_Parser("file.txt"); p != nil {
		t.Error("expected nil parser for wrong extension")
	}
}
