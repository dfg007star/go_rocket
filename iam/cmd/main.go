package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/jcmturner/gokrb5/v8/config"

	"github.com/dfg007star/go_rocket/platform/pkg/closer"
)

const configPath = "../deploy/compose/iam/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)
}
