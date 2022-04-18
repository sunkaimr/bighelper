package service

import (
	"context"
	"log"
	"runtime"

	"github.com/kardianos/service"

	"bighelper/bigiot"
)

type program struct {
	ctx      context.Context
	deviceID string
	apiKey   string
}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go bigiot.StartServer(p.ctx, p.deviceID, p.apiKey)
	return nil
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func StartService(ctx context.Context, devID, apiKey string) {
	switch runtime.GOOS {
	case "linux":
		bigiot.StartServer(ctx, devID, apiKey)

	case "windows":
		program := &program{
			ctx:      ctx,
			deviceID: devID,
			apiKey:   apiKey,
		}
		s, err := service.New(program, &service.Config{
			Name:        "bighelper",
			DisplayName: "bighelper",
			Description: "基于贝壳物联平台远程控制工具",
		})
		if err != nil {
			log.Fatal(err)
		}

		err = s.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
