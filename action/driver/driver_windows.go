package driver

import (
	"fmt"
	"github.com/CodyGuo/win"
	"os/exec"
	"strings"
)

type Driver struct{}

func (la *Driver) ShutDown(cmd string) error {
	out, err := exec.Command("shutdown", "-s", "-t", "60").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Driver) Reboot(cmd string) error {
	out, err := exec.Command("shutdown", "-r", "-t", "60").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Driver) Sleep(cmd string) error {
	out, err := exec.Command("shutdown", "-h").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Driver) Cancel(cmd string) error {
	out, err := exec.Command("shutdown", "-a").Output()
	if err != nil {
		return fmt.Errorf("exec command failed, err: %v, outpout:%v", err, out)
	}
	return nil
}

func (la *Driver) Custom(cmd string) error {
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

func getPrivileges() {
	var hToken win.HANDLE
	var tkp win.TOKEN_PRIVILEGES

	win.OpenProcessToken(win.GetCurrentProcess(), win.TOKEN_ADJUST_PRIVILEGES|win.TOKEN_QUERY, &hToken)
	win.LookupPrivilegeValueA(nil, win.StringToBytePtr(win.SE_SHUTDOWN_NAME), &tkp.Privileges[0].Luid)
	tkp.PrivilegeCount = 1
	tkp.Privileges[0].Attributes = win.SE_PRIVILEGE_ENABLED
	win.AdjustTokenPrivileges(hToken, false, &tkp, 0, nil, nil)
}
