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
	"bighelper/bigiot"
	"bighelper/config"
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Printf("version 1.0")
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Printf("fail loadconfig, err:%v", err)
		time.Sleep(time.Second * 5)
		os.Exit(1)
	}

	devID := cfg.Section("bigiot").Key("device_id").String()
	if devID == "" {
		log.Printf("read config from bigiot.ini, but bigiot.device_id is null")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	} else {
		log.Printf("device_id:%s", devID)
	}

	apiKey := cfg.Section("bigiot").Key("api_key").String()
	if apiKey == "" {
		log.Printf("read config from bigiot.ini, but bigiot.api_key is null")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	} else {
		log.Printf("api_key:%s", apiKey)
	}

	exitCh := make(chan os.Signal)
	signal.Notify(exitCh, os.Interrupt)
	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		<-exitCh
		cancel()
		log.Printf("main recv exit signal")
		time.Sleep(time.Second * 5)
		log.Printf("waited 10s to force quit")
		os.Exit(0)
	}()

	bigiot.StartServer(ctx, devID, apiKey)
}
