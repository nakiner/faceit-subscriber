package main

import (
	"context"
	"fmt"
	"github.com/go-kit/log/level"
	"github.com/nakiner/faceit-subscriber/configs"
	"github.com/nakiner/faceit-subscriber/tools/logging"
	"github.com/nakiner/faceit/pkg/queue/user"
	"github.com/nakiner/faceit/pkg/store/nats"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load config
	cfg := configs.NewConfig()
	if err := cfg.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}
	// Print config
	if err := cfg.Print(); err != nil {
		fmt.Fprintf(os.Stderr, "read config: %s", err)
		os.Exit(1)
	}

	logger, err := logging.NewLogger(cfg.Logger.Level, cfg.Logger.TimeFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %s", err)
		os.Exit(1)
	}
	ctx = logging.WithContext(ctx, logger)

	nc, err := nats.NewClient(&cfg.Nats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init nats: %s", err)
		os.Exit(1)
	}
	defer nc.Close()

	subscriber := user.NewSubscriber(nc)
	go func() {
		err = subscriber.UpdateUser(func(u *user.User) {
			level.Info(logger).Log("msg", "got update from channel", "user_id", u.ID)
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to init NewSubscriber.UpdateUser: %s", err)
		}
	}()

	ch := make(chan bool)
	<-ch
}
