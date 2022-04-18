package linux

import (
	"bighelper/driver"
	"fmt"
	"os/exec"
	"strings"
)

var LinuxAction *Action

func init(){
	driver.RegistAction("linux", &Action{})
}

func (la *Action) ShutDown(cmd string) error {
	// 1分钟后关机
	out, err := exec.Command("shutdown", "-h", "1").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Action) Reboot(cmd string) error {
	// 1分钟后重启
	out, err := exec.Command("shutdown", "-r").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Action) Cancel(cmd string) error {
	// 取消关机
	out, err := exec.Command("shutdown", "-c").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Action) Custom(cmd string) error {
	cmds := strings.Split(cmd, " ")

	if len(cmds) == 0 {
		return fmt.Errorf("Custom command cannot be emapy")
	} else if len(cmds) == 1 {
		out, err := exec.Command(cmds[0]).Output()
		if err != nil {
			return fmt.Errorf("exec command[%s] failed, err: %v, outpout:%v", cmds[0], err, out)
		}
	} else {
		out, err := exec.Command(cmds[0], cmds[1:]...).Output()
		if err != nil {
			return fmt.Errorf("exec command[%s %s] failed, err: %v, outpout:%v", cmds[0], cmds[1:], err, out)
		}
	}

	return nil
}
