package main

import (
	"context"
	"github.com/deyuro/digger/internal/config"
	"github.com/deyuro/digger/internal/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func app() error {
	appCtx, cancel := context.WithCancel(context.Background())

	errChan := make(chan error)
	go func() {
		handleSIGINT()
		cancel()
	}()


	go func() {
		errChan <- run(logrus.New())
	}()
	select {
	case err := <-errChan:
		return err
	case <-appCtx.Done():
		return nil
	}
}




func run(logger *logrus.Logger) error {
	digList, err := config.Load("/home/droot/go/src/github.com/deyuro/digger/.data/diglist.json")
	svc := service.NewService(digList, logger, `127.0.0.53`, 5, 0)
	if err != nil {
		return err
	}

	svc.Run()
	return nil
}

func handleSIGUSR2(logger *logrus.Logger) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2)
	for range ch {
		level := logger.GetLevel()
		switch level {
		case logrus.DebugLevel:
			logger.Warn("switching log level to INFO")
			logger.SetLevel(logrus.InfoLevel)
		default:
			logger.Warn("switching log level to DEBUG")
			logger.SetLevel(logrus.DebugLevel)
		}
	}
}

func handleSIGINT() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	for range sigCh {
		signal.Stop(sigCh)
		return
	}
}