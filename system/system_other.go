//go:build !windows
// +build !windows

package system

import (
	"os"
	"os/exec"
)

// 在用户账号下执行
func LaunchProcessWithUser(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Env = os.Environ()
	return cmd.Start()
}
