package action

import (
	driver2 "bighelper/action/driver"
	"fmt"
	"runtime"
)

type CmdType int

const (
	// 内置类型的命令
	CmdTypeBuiltin = 0
	// 别名类型
	CmdTypeAlias
	// 自定义类型的命令
	CmdTypeCustom
)

type Command struct {
	Type   CmdType
	Cmd    string
	Handle func(string) error
}

type Driver interface {
	ShutDown(string) error
	Reboot(string) error
	Sleep(string) error
	Cancel(string) error
	Custom(string) error
}

var DriverMap = make(map[string]Driver, 2)
var Cmds = make(map[string]*Command, 10)

func registAction() {
	DriverMap[runtime.GOOS] = &driver2.Driver{}
}

func getAction(driver string) Driver {
	if d, ok := DriverMap[driver]; ok {
		return d
	}
	return nil
}

func RegistBuiltinCommands() error {
	registAction()

	drv := getAction(runtime.GOOS)
	if drv == nil {
		return fmt.Errorf("get action failed, drv is nil")
	}

	Cmds["shutdown"] = &Command{
		Type:   CmdTypeBuiltin,
		Cmd:    "",
		Handle: drv.ShutDown,
	}
	Cmds["reboot"] = &Command{
		Type:   CmdTypeBuiltin,
		Cmd:    "",
		Handle: drv.Reboot,
	}
	Cmds["sleep"] = &Command{
		Type:   CmdTypeBuiltin,
		Cmd:    "",
		Handle: drv.Sleep,
	}
	Cmds["cancel"] = &Command{
		Type:   CmdTypeBuiltin,
		Cmd:    "",
		Handle: drv.Cancel,
	}
	return nil
}

func RegistAliasCommands(key, value string) error {
	if key == "" || value == "" {
		return nil
	}

	if c, ok := Cmds[key]; ok && c.Type == CmdTypeBuiltin {
		Cmds[value] = &Command{
			Type:   CmdTypeAlias,
			Cmd:    c.Cmd,
			Handle: c.Handle,
		}
		return nil
	}
	return fmt.Errorf("not fuond buildin command %s", key)
}

func RegistCustomCommands(key, value string) error {
	if key == "" || value == "" {
		return nil
	}

	drv := getAction(runtime.GOOS)
	if drv == nil {
		return fmt.Errorf("get action failed, drv is nil")
	}

	Cmds[key] = &Command{
		Type:   CmdTypeCustom,
		Cmd:    value,
		Handle: getAction(runtime.GOOS).Custom,
	}

	return nil
}
