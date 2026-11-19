//go:build windows

package deps

import "syscall"

// sysProcAttr Windows 下无需特殊进程组处理。
func sysProcAttr() *syscall.SysProcAttr {
	return nil
}
