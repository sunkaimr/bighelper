package bigiot

import (
	"bighelper/driver"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	bigiotHostName = "www.bigiot.net"
	bigiotPort     = 8181
	connectted     = "WELCOME TO BIGIOT"
	checkinok      = "checkinok"
)

type bigiotCommand struct {
	M    string `json:"M"`
	ID   string `json:"ID"`
	NAME string `json:"NAME"`
	K    string `json:"K"`
	V    string `json:"V"`
	C    string `json:"C"`
}

type bigiotConn struct {
	hostName  string
	port      int
	deviceID  string
	apiKey    string
	conn      net.Conn
	timeout   time.Duration
	heartBeat time.Duration
	wg        sync.WaitGroup
	retry     bool
	retryCh   chan struct{}
}

func StartServer(ctx context.Context, devID, apiKey string) {

	driver.RegistAction()
	for {
		conn := newBigiot(bigiotHostName, bigiotPort, devID, apiKey)
		log.Print("newBigiot success")

		err := login(conn)
		if err != nil {
			log.Printf("login failed, err:%v, reconnect after 3 seconds", err)

			t := time.NewTicker(time.Second * 3)
			select {
			case <-ctx.Done():
				log.Print("bigiot server got exit signal")
				return
			case <-t.C:
				continue
			}
		}
		log.Print("login success")

		go execCommand(ctx, conn)

		go heartbeat(ctx, conn)

		go broadcaseRetrySignal(ctx, conn)

		select {
		case <-conn.retryCh:
			log.Print("bigiot server got retry signal")
			conn.wg.Wait()
		case <-ctx.Done():
			log.Print("bigiot server got exit signal, wait all coroutines exit")
			conn.wg.Wait()
			log.Print("all coroutines have exited")
			logout(conn)
			return
		}
	}
}

func newBigiot(hostname string, port int, deviceId, apiKey string) *bigiotConn {
	return &bigiotConn{
		hostName:  hostname,
		port:      port,
		deviceID:  deviceId,
		apiKey:    apiKey,
		conn:      nil,
		timeout:   time.Second * 1,
		heartBeat: time.Second * 20,
		wg:        sync.WaitGroup{},
		retry:     false,
		retryCh:   make(chan struct{}),
	}
}

func login(conn *bigiotConn) (err error) {
	log.Print("login...")
	conn.conn, err = net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", conn.hostName, conn.port), conn.timeout)
	if err != nil {
		err = fmt.Errorf("connect bigiot failed, err:%v", err)
		return err
	}

	defer func() {
		if err != nil {
			conn.conn.Close()
		}
	}()

	con := make([]byte, 1204)
	conn.conn.SetReadDeadline(time.Now().Add(conn.timeout))
	len, err := conn.conn.Read(con)
	if err != nil {
		err = fmt.Errorf("get login result failed, err:%v", err)
		return err
	}

	result := &bigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		err = fmt.Errorf("Unmarshal result failed, err:%v", err)
		return err
	}

	if result.M != connectted {
		err = fmt.Errorf("connect bigiot failed, expect %s, but got %s", connectted, result.M)
		return err
	}

	cmd := bigiotCommand{
		M:  "checkin",
		ID: conn.deviceID,
		K:  conn.apiKey,
	}
	cmdByte, err := json.Marshal(cmd)
	if err != nil {
		err = fmt.Errorf("login instruction splicing failed, err:%v", err)
		return err
	}

	conn.conn.SetWriteDeadline(time.Now().Add(conn.timeout))
	_, err = conn.conn.Write(append(cmdByte, '\n'))
	if err != nil {
		err = fmt.Errorf("send login cmd failed, err:%v", err)
		return err
	}

	conn.conn.SetReadDeadline(time.Now().Add(conn.timeout))
	len, err = conn.conn.Read(con)
	if err != nil {
		err = fmt.Errorf("get login result failed, err:%v. maybe deviceID or APIKey invalid", err)
		return err
	}

	result = &bigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		err = fmt.Errorf("Unmarshal result failed, err:%v", err)
		return err
	}

	if result.M != checkinok {
		err = fmt.Errorf("login failed, err:%v. expect %s, but got %s", err, checkinok, result.M)
	}

	return err
}

func logout(bigiotConn *bigiotConn) (err error) {
	bigiotConn.wg.Add(1)
	defer func() {
		bigiotConn.wg.Done()
	}()

	if bigiotConn.conn == nil {
		return nil
	}

	defer func() {
		if err == nil {
			bigiotConn.conn.Close()
		}
	}()

	cmd := bigiotCommand{
		M:  "checkout",
		ID: bigiotConn.deviceID,
		K:  bigiotConn.apiKey,
	}
	cmdByte, err := json.Marshal(cmd)
	if err != nil {
		err = fmt.Errorf("logout instruction splicing failed, err:%v", err)
		return err
	}

	bigiotConn.conn.SetWriteDeadline(time.Now().Add(bigiotConn.timeout))
	_, err = bigiotConn.conn.Write(append(cmdByte, '\n'))
	if err != nil {
		err = fmt.Errorf("send logout cmd failed, err:%v", err)
		return err
	}

	con := make([]byte, 1204)
	bigiotConn.conn.SetReadDeadline(time.Now().Add(bigiotConn.timeout))
	len, err := bigiotConn.conn.Read(con)
	if err != nil {
		err = fmt.Errorf("get logout result failed, err:%v", err)
		return err
	}

	result := &bigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		err = fmt.Errorf("Unmarshal result failed, err:%v", err)
		return err
	}

	if result.M != checkinok {
		err = fmt.Errorf("logout failed, M:%s\n", result.NAME)
	}

	return err
}

func heartbeat(ctx context.Context, conn *bigiotConn) {
	log.Print("heartbeat...")

	conn.wg.Add(1)
	defer func() {
		conn.wg.Done()
	}()

	t := time.NewTicker(conn.heartBeat)
	for {
		select {
		case <-t.C:
			err := sendHeartbeat(conn)
			if err != nil {
				conn.retry = true
				log.Printf("heartbeat failed and will retry, err:%v", err)
				return
			}
			log.Printf("heartbeat success")

		case <-ctx.Done():
			log.Printf("heartbeat got exit signal")
			return
		case <-conn.retryCh:
			log.Printf("heartbeat got retry signal")
			return
		}
	}
}

func sendHeartbeat(bigiotConn *bigiotConn) (err error) {
	bigiotConn.conn.SetWriteDeadline(time.Now().Add(bigiotConn.timeout))
	n, err := bigiotConn.conn.Write([]byte("{\"M\":\"beat\"}\n"))
	if err != nil {
		err = fmt.Errorf("send heartbeat failed, err:%v", err)
	}
	if n == 0 {
		err = fmt.Errorf("send heartbeat failed, send 0 byte")
	}

	return err
}

func execCommand(ctx context.Context, conn *bigiotConn) {
	log.Print("execCommand...")

	conn.wg.Add(1)
	defer func() {
		conn.wg.Done()
	}()

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Print("execCommand got exit signal")
			return
		case <-t.C:
			//log.Print("execCommand recv")
		case <-conn.retryCh:
			log.Print("execCommand got retry signal")
			return
		}

		command, err := recvCommand(conn)
		if err != nil {
			conn.retry = true
			log.Printf("recv command failed, err:%v", err)
			return
		}

		if command != "" {
			doAction(command)
		}
	}
}

func broadcaseRetrySignal(ctx context.Context, conn *bigiotConn) {
	log.Print("broadcaseRetrySignal...")

	conn.wg.Add(1)
	defer func() {
		if conn.retry {
			close(conn.retryCh)
			log.Printf("broadcase notifying  all collaborators to exit")
		}
		conn.wg.Done()
	}()

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			log.Printf("broadcase got exit signal, and will exit")
			return
		case <-t.C:
		}

		if conn.retry {
			log.Printf("sameone exit, will broadcase retry signal to all")
			return
		}
	}
}

func recvCommand(conn *bigiotConn) (string, error) {
	con := make([]byte, 1204)
	conn.conn.SetReadDeadline(time.Now().Add(conn.timeout))
	len, err := conn.conn.Read(con)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return "", nil
		}

		err = fmt.Errorf("conn read failed, err:%v", err)
		if len == 0 {
			return "", err
		}
	}

	result := &bigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		log.Printf("Unmarshal result failed, err:%v, content:[%s]", err, string(con[:len]))
		return "", nil
	}

	return result.C, err
}

func doAction(action string) {

	drv := driver.GetAction(runtime.GOOS)
	if drv == nil {
		log.Printf("doAction failed, get cannot get driver")
		return
	}

	switch action {
	// 关机
	case "shutdown":
		log.Print("recv command:", action)
		err := drv.ShutDown()
		if err != nil {
			log.Printf("%v", err)
		}

		// 重启
	case "reboot":
		log.Print("recv command:", action)
		err := drv.Reboot()
		if err != nil {
			log.Printf("%v", err)
		}

		// 取消 关机、重启
	case "cancel":
		log.Print("recv command:", action)
		err := drv.Cancel()
		if err != nil {
			log.Printf("%v", err)
		}

		// 其他未识别的命令
	default:
		log.Print("unknown command:", action)
	}
}
