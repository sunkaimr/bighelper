package driver

import (
	"bighelper/driver/linux"
	"bighelper/driver/win"
)

var DriverMap = make(map[string]Driver, 2)

type Driver interface {
	ShutDown() error
	Reboot() error
	Cancel() error
}

func RegistAction() {
	DriverMap["linux"] = linux.LinuxAction
	DriverMap["windows"] = win.WinAction
}

func GetAction(driver string) Driver {
	if d, ok := DriverMap[driver]; ok {
		return d
	}
	return nil
}
