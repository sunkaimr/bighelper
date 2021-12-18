/*
 * Copyright sunkai
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/ini.v1"
)

const (
	BigiotHostName = "www.bigiot.net"
	BigiotPort     = 8181

	Connectted = "WELCOME TO BIGIOT"
	Checkinok  = "checkinok"

	ConfigFileName = "bighelper.ini"
)

var (
	DeviceId = "12345"
	ApiKey   = "asdfghjkl"
)

type BigiotCommand struct {
	M    string `json:"M"`
	ID   string `json:"ID"`
	NAME string `json:"NAME"`
	K    string `json:"K"`
	V    string `json:"V"`
	C    string `json:"C"`
}

type BigiotConn struct {
	hostName  string
	port      int
	deviceID  string
	apiKey    string
	conn      net.Conn
	timeout   time.Duration
	heartBeat time.Duration
	wg        sync.WaitGroup
	exit      bool
	exitChan  chan struct{}
}

var Conn *BigiotConn

func main() {
	log.Printf("version 1.0")
	StartServer()
}

func fileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func loadConfig() (string, error) {
	configPath := filepath.Join("./", ConfigFileName)
	if !fileExist(configPath) {
		file, _ := exec.LookPath(os.Args[0])
		path, _ := filepath.Abs(file)
		index := strings.LastIndex(path, string(os.PathSeparator))
		configAbsPath := filepath.Join(path[:index], ConfigFileName)

		if !fileExist(configAbsPath) {
			return "", fmt.Errorf("config file not found in %v or %v", configAbsPath, configPath)
		}
		return configAbsPath, nil
	}

	return configPath, nil
}

func StartServer() {
	configPath, err := loadConfig()
	if err != nil {
		log.Printf("Fail to find config file: %v", err)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}

	cfg, err := ini.Load(configPath)
	if err != nil {
		log.Printf("Fail to load config file: %v", err)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}

	DeviceId = cfg.Section("bigiot").Key("device_id").String()
	if DeviceId == "" {
		log.Printf("read config from bigiot.ini, but bigiot.device_id is null")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	} else {
		log.Printf("device_id:%s", DeviceId)
	}

	ApiKey = cfg.Section("bigiot").Key("api_key").String()
	if ApiKey == "" {
		log.Printf("read config from bigiot.ini, but bigiot.api_key is null")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	} else {
		log.Printf("api_key:%s", ApiKey)
	}

	for {
		Conn = NewBigiot(BigiotHostName, BigiotPort, DeviceId, ApiKey)
		log.Print("NewBigiot success")

		err := Login(Conn)
		if err != nil {
			log.Printf("Login failed, err:%v, reconnect after 3 seconds", err)
			time.Sleep(time.Second * 3)
			continue
		} else {
			log.Print("Login success")

			go ExecCommand(Conn)
			log.Print("ExecCommand...")

			go Heartbeat(Conn)
			log.Print("Heartbeat...")

			go BroadcaseExitSignal(Conn)
		}

		select {
		case <-Conn.exitChan:
			log.Print("main received exit signal, and will retry")
			Logout(Conn)
			Conn.wg.Wait()
		}
	}
}

func NewBigiot(hostname string, port int, deviceId, apiKey string) *BigiotConn {
	return &BigiotConn{
		hostName:  hostname,
		port:      port,
		deviceID:  deviceId,
		apiKey:    apiKey,
		conn:      nil,
		timeout:   time.Second * 3,
		heartBeat: time.Second * 20,
		wg:        sync.WaitGroup{},
		exit:      false,
		exitChan:  make(chan struct{}),
	}
}

func Login(bigiotConn *BigiotConn) (err error) {
	log.Print("Login...")
	bigiotConn.conn, err = net.DialTimeout("tcp",
		fmt.Sprintf("%s:%d", bigiotConn.hostName, bigiotConn.port), bigiotConn.timeout)
	if err != nil {
		err = fmt.Errorf("connect bigiot failed, err:%v", err)
		return err
	}

	defer func() {
		if err != nil {
			bigiotConn.conn.Close()
		}
	}()

	con := make([]byte, 1204)
	bigiotConn.conn.SetReadDeadline(time.Now().Add(bigiotConn.timeout))
	len, err := bigiotConn.conn.Read(con)
	if err != nil {
		err = fmt.Errorf("get login result failed, err:%v", err)
		return err
	}

	result := &BigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		err = fmt.Errorf("Unmarshal result failed, err:%v", err)
		return err
	}

	if result.M != Connectted {
		err = fmt.Errorf("connect bigiot failed, expect %s, but got %s", Connectted, result.M)
		return err
	}

	cmd := BigiotCommand{
		M:  "checkin",
		ID: DeviceId,
		K:  ApiKey,
	}
	cmdByte, err := json.Marshal(cmd)
	if err != nil {
		err = fmt.Errorf("login instruction splicing failed, err:%v", err)
		return err
	}

	bigiotConn.conn.SetWriteDeadline(time.Now().Add(bigiotConn.timeout))
	_, err = bigiotConn.conn.Write(append(cmdByte, '\n'))
	if err != nil {
		err = fmt.Errorf("send login cmd failed, err:%v", err)
		return err
	}

	bigiotConn.conn.SetReadDeadline(time.Now().Add(bigiotConn.timeout))
	len, err = bigiotConn.conn.Read(con)
	if err != nil {
		err = fmt.Errorf("get login result failed, err:%v. maybe deviceID or APIKey invalid", err)
		return err
	}

	result = &BigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		err = fmt.Errorf("Unmarshal result failed, err:%v", err)
		return err
	}

	if result.M != Checkinok {
		err = fmt.Errorf("login failed, err:%v. expect %s, but got %s", err, Checkinok, result.M)
	}

	return err
}

func Logout(bigiotConn *BigiotConn) (err error) {
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

	cmd := BigiotCommand{
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

	result := &BigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		err = fmt.Errorf("Unmarshal result failed, err:%v", err)
		return err
	}

	if result.M != Checkinok {
		err = fmt.Errorf("logout failed, M:%s\n", result.NAME)
	}

	return err
}

func Heartbeat(bigiotConn *BigiotConn) {
	bigiotConn.wg.Add(1)
	defer func() {
		bigiotConn.wg.Done()
		bigiotConn.exit = true
		log.Print("Heartbeat exited")
	}()

	t := time.NewTicker(bigiotConn.heartBeat)

	for {
		select {
		case <-t.C:
			//log.Print("heartbeat...")
			err := sendHeartbeat(bigiotConn)
			if err != nil {
				log.Print("heartbeat failed and will exit, err:%v", err)
				return
			}
			log.Print("heartbeat success")

		case <-bigiotConn.exitChan:
			log.Print("heartbeat received exit signal")
			return
		}
	}
}

func sendHeartbeat(bigiotConn *BigiotConn) (err error) {
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

func ExecCommand(bigiotConn *BigiotConn) {
	bigiotConn.wg.Add(1)
	defer func() {
		bigiotConn.wg.Done()
		bigiotConn.exit = true
		log.Print("ExecCommand exited")
	}()

	t := time.NewTicker(time.Second)
	for {
		select {
		case <-bigiotConn.exitChan:
			log.Print("execCommand received exit signal")
			return
		case <-t.C:
			//log.Print("execCommand recv")
		}

		com, err := recvCommand(bigiotConn)
		if err != nil {
			log.Printf("recv command failed, err:%v", err)
			return
		}

		if com != "" {
			doAction(com)
		}
	}
}
func recvCommand(bigiotConn *BigiotConn) (string, error) {
	con := make([]byte, 1204)
	bigiotConn.conn.SetReadDeadline(time.Now().Add(bigiotConn.timeout))
	len, err := bigiotConn.conn.Read(con)
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return "", nil
		}

		err = fmt.Errorf("get logout result failed, err:%v", err)
		if len == 0 {
			return "", err
		}
	}

	result := &BigiotCommand{}
	err = json.Unmarshal(con[:len], result)
	if err != nil {
		log.Printf("Unmarshal result failed, err:%v, content:[%s]", err, string(con[:len]))
		return "", nil
	}

	return result.C, err
}

func BroadcaseExitSignal(bigiotConn *BigiotConn) {
	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt)
	t := time.NewTicker(time.Second)

	defer func() {
		close(bigiotConn.exitChan)
		log.Printf("BroadcaseExitSignal exit")
	}()

	for {
		select {
		case s := <-signalChan:
			log.Printf("recv signal %v, and will exit", s)
			os.Exit(0)
		case <-t.C:
		}

		if bigiotConn.exit {
			log.Printf("sameone exit, will send exit signal to all")
			return
		}
	}
}

func doAction(action string) {
	switch action {
		// 关机
	case "shutdown":
		log.Print("recv command:", action)
		out, err := exec.Command("shutdown", "-H", "-t", "15").Output()
		if err != nil {
			log.Printf("exec.Command failed, err: %v, outpout:%v", err, out)
		}

		// 重启
	case "reboot":
		log.Print("recv command:", action)
		out, err := exec.Command("shutdown", "-r", "-t", "15").Output()
		if err != nil {
			log.Printf("exec.Command failed, err: %v, outpout:%v", err, out)
		}

		// 取消 关机、重启
	case "cancel":
		log.Print("recv command:", action)
		out, err := exec.Command("shutdown", "-c").Output()
		if err != nil {
			log.Printf("exec.Command failed, err: %v, outpout:%v", err, out)
		}

		// 其他未识别的命令
	default:
		log.Print("unknown command:", action)
	}
}
