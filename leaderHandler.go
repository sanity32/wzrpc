package wzrpc

import (
	"fmt"
	"net"
	"time"
)

func NewServerlessLeaderHandler(host string, port int) *LeaderHandler {
	r := LeaderHandler{
		Host:         host,
		Port:         port,
		ProbeTimeout: time.Second,
		Relay:        NewLeaderClient(host, port),
		Server:       NewServer(host, port),
	}
	r.ServerMode = false
	r.SetupOk = true
	return &r
}

func NewLeaderHandler(host string, port int) *LeaderHandler {
	r := LeaderHandler{
		Host:         host,
		Port:         port,
		ProbeTimeout: time.Second,
		Relay:        NewLeaderClient(host, port),
		Server:       NewServer(host, port),
	}
	r.Setup(true)
	return &r
}

type LeaderHandler struct {
	ServerMode   bool
	Host         string
	Port         int
	Relay        *LeaderClient
	Server       *Server
	ProbeTimeout time.Duration
	SetupOk      bool
}

func (s *LeaderHandler) addr() string {
	return fmt.Sprintf("%v:%v", s.Host, s.Port)
}

func (s *LeaderHandler) isServerAlreadyExist() bool {
	_, err := net.DialTimeout("tcp", s.addr(), s.ProbeTimeout)
	return err == nil
}

func (s *LeaderHandler) Setup(performAnyway bool) {
	if !s.SetupOk || performAnyway {
		if !s.isServerAlreadyExist() {
			go s.Server.Serve()
			s.ServerMode = true
		}
		s.SetupOk = true
	}
}

func (s *LeaderHandler) Leader() Leader {
	if s.ServerMode {
		return s.Server.performer
	}
	return s.Relay
}
