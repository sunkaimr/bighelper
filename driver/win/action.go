package win

import (
	"github.com/CodyGuo/win"
	//"github.com/kardianos/service"
	"log"
)

/*
  TODO 详细的参数参考: https://docs.microsoft.com/zh-cn/windows/win32/api/winuser/nf-winuser-exitwindowsex?redirectedfrom=MSDN
*/

var WinAction *Action

type Action struct{}

/* 	关机
限制: 所有文件都已写入磁盘，所有软件都已关闭。如果有其他软件阻止，则无法关闭
*/
func (la *Action) ShutDown() error {
	//getPrivileges()
	//if win.ExitWindowsEx(win.EWX_SHUTDOWN|win.EWX_FORCE, 0) {
	//	return fmt.Errorf("shutdown exec failed")
	//}
	log.Print("ShutDown")
	return nil
}

func (la *Action) Reboot() error {
	//getPrivileges()
	//if win.ExitWindowsEx(win.EWX_REBOOT|win.EWX_FORCE, 0){
	//	return fmt.Errorf("reboot exec failed")
	//}
	log.Print("Reboot")
	return nil
}

func (la *Action) Cancel() error {
	log.Print("not support cancel")
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
