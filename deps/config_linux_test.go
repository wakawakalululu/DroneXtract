//go:build linux || darwin

package deps

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestConfigFileIsUnderAppDir 验证配置固定落在程序所在目录。
func TestConfigFileIsUnderAppDir(t *testing.T) {
	got := ConfigFile()
	want := filepath.Join(AppDir(), ".env")
	if got != want {
		t.Errorf("ConfigFile = %q, want %q", got, want)
	}
}

// TestConfigPathsPrefersAppDir 验证程序目录优先于工作目录和历史用户目录。
func TestConfigPathsPrefersAppDir(t *testing.T) {
	t.Setenv("DRONEXTRACT_ENV", "")
	paths := ConfigPaths()
	if len(paths) == 0 {
		t.Fatal("no candidates")
	}
	if paths[0] != ConfigFile() {
		t.Errorf("first candidate should be app dir config, got %q", paths[0])
	}
	// 历史用户目录只作为兜底，不应出现在前面
	for i, p := range paths {
		if strings.Contains(p, filepath.Join("dronextract", ".env")) && i < 2 {
			t.Errorf("legacy user config %q should not outrank app dir", p)
		}
	}
}

// TestExplicitEnvWins 验证 DRONEXTRACT_ENV 可以覆盖默认位置。
func TestExplicitEnvWins(t *testing.T) {
	explicit := filepath.Join(t.TempDir(), "custom.env")
	t.Setenv("DRONEXTRACT_ENV", explicit)
	if got := ConfigPaths()[0]; got != explicit {
		t.Errorf("explicit config should be first, got %q", got)
	}
}

// TestOutputDirUnderAppDir 验证默认输出目录跟随程序目录。
func TestOutputDirUnderAppDir(t *testing.T) {
	if got, want := OutputDir(), filepath.Join(AppDir(), "output"); got != want {
		t.Errorf("OutputDir = %q, want %q", got, want)
	}
}

// TestEnsureParentCreatesDir 验证输出目录会自动创建。
func TestEnsureParentCreatesDir(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "output", "nested", "a.csv")
	if _, err := EnsureParent(target); err != nil {
		t.Fatalf("ensure parent: %v", err)
	}
	info, err := os.Stat(filepath.Dir(target))
	if err != nil || !info.IsDir() {
		t.Fatalf("expected directory to exist: %v", err)
	}
}

// TestInstallNeedsRoot 验证无 root 且无 sudo 时不会执行注定失败的命令。
func TestInstallNeedsRoot(t *testing.T) {
	if isRoot() {
		t.Skip("running as root")
	}
	if hasSudo() {
		t.Skip("sudo available")
	}
	res := InstallExifTool()
	if !res.Skipped && res.Err == nil && res.Started {
		t.Error("expected install to be skipped without root/sudo")
	}
}
