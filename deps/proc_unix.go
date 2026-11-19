//go:build !windows

package deps

import "syscall"

// sysProcAttr 让子进程进入独立进程组：
// 终端 Ctrl+C 只中断本程序，不会把包管理器置于半途状态。
func sysProcAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setpgid: true}
}
