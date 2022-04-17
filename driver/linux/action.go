package linux

import (
	"fmt"
	"os/exec"
)

var LinuxAction *Action

type Action struct{}

func (la *Action) ShutDown() error {
	// 1分钟后关机
	out, err := exec.Command("shutdown", "-h", "1").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Action) Reboot() error {
	// 1分钟后重启
	out, err := exec.Command("shutdown", "-r").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Action) Cancel() error {
	// 取消关机
	out, err := exec.Command("shutdown", "-c").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}
