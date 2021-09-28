package service

import (
	"fmt"
	"github.com/deyuro/digger/internal/config"
	"github.com/sirupsen/logrus"
	"os/exec"
	"time"
)

type service struct {
	digList    config.DigList
	logger     *logrus.Logger
	resolverIP string
	times      int
	wait       time.Duration
	output     bool
}

func NewService(digList config.DigList, logger *logrus.Logger, resolverIP string, times int, wait time.Duration) *service {
	return &service{digList: digList, logger: logger, resolverIP: fmt.Sprintf("@%s", resolverIP), times: times, wait: wait, output: true}
}

func (s service) Run() {
	if s.times == 0 {
		s.infinity()
	} else {
		s.repeatedly()
	}
}

func (s *service) infinity() {
	for {
		for _, v := range s.digList.List {
			s.dig(v.Name)
		}
	}
}

func (s *service) repeatedly() {
	for i := 0; i <= s.times; i++ {
		for _, v := range s.digList.List {
			s.dig(v.Name)
		}
	}
}

func (s *service) dig(host string) {
	time.Sleep(s.wait)
	s.cmd(host)
}

func (s *service) cmd(host string) {
	cmd := exec.Command("dig", host, s.resolverIP)
	if s.output {
		stdOut, _ := cmd.Output()

		s.logger.Info(string(stdOut))
	}
}
