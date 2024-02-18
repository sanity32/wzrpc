package wzrpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func NewServer(host string, port int) *Server {
	return &Server{
		Settings:  NewServerSettings(host, port),
		Mux:       http.NewServeMux(),
		performer: NewModel(),
	}
}

type Server struct {
	performer   *Model
	Settings    ServerSettings
	Mux         *http.ServeMux
	httpserver  *http.Server
	netlistener net.Listener
}

func (s *Server) Shutdown() {
	if l := s.netlistener; l != nil {
		l.Close()
	}
}

func (s Server) servingPort() string {
	return fmt.Sprintf(":%v", s.Settings.Port)
}

func parseServerRequest[T any](req *http.Request) (r T, err error) {
	buff, err := io.ReadAll(req.Body)
	if err != nil {
		return
	}
	return r, json.Unmarshal(buff, &r)
}

func out[T any](w http.ResponseWriter, data T) error {
	j, err := json.MarshalIndent(data, "", "	")
	if err != nil {
		return err
	}
	_, err = w.Write(j)
	return err
}

func (serv *Server) Serve() (err error) {
	outSuccess := func(w http.ResponseWriter) {
		j, _ := json.MarshalIndent(map[string]bool{"success": true}, "", "	")
		w.Write(j)
	}

	outErr := func(w http.ResponseWriter, err error) {
		if err == nil {
			outSuccess(w)
		} else {
			j, _ := json.MarshalIndent(map[string]string{"error": err.Error()}, "", "	")
			w.Write(j)
		}
	}

	t := serv.Settings
	serv.netlistener, err = net.Listen("tcp", serv.servingPort())
	if err != nil {
		return err
	}

	// h(t.LeaderAskEndpoint, e.LeaderAsk)
	serv.Mux.HandleFunc(t.LeaderAskEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if req, err := parseServerRequest[Ask](r); err != nil {
			outErr(w, err)
		} else {
			outErr(w, serv.performer.LeaderAsk(req))
		}
	})

	// h(t.LeaderAnswerEndpoint, e.LeaderAnswer)
	serv.Mux.HandleFunc(t.LeaderAnswerEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if req, err := parseServerRequest[LeaderAnswerRequest](r); err != nil {
			outErr(w, err)
		} else {
			switch answer, err := serv.performer.LeaderAnswer(req.SID); err {
			case nil:
				out(w, answer)
			default:
				outErr(w, err)
			}
		}
	})

	// h(t.FollowerAskEndpoint, e.FollowerAsk)
	serv.Mux.HandleFunc(t.FollowerAskEndpoint, func(w http.ResponseWriter, r *http.Request) {
		switch ask, err := serv.performer.FollowerAsk(); err {
		case nil:
			out(w, ask)
		case ErrAskTimeout:
			out(w, map[string]bool{"timeout": true})
		default:
			outErr(w, err)
		}
	})

	// h(t.FollowerAliveEndpoint, e.FollowerAlive)
	serv.Mux.HandleFunc(t.FollowerAliveEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if req, err := parseServerRequest[AliveDTO](r); err != nil {
			outErr(w, err)
		} else {
			serv.performer.FollowerAlive(req)
			outSuccess(w)
		}
	})

	// h(t.FollowerAnswerEndpoint, e.FollowerAnswer)
	serv.Mux.HandleFunc(t.FollowerAnswerEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if answer, err := parseServerRequest[Answer](r); err != nil {
			outErr(w, err)
		} else {
			outErr(w, serv.performer.FollowerAnswer(answer))
		}
	})

	serv.Mux.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("simple domestic starter page"))
	})

	serv.Mux.HandleFunc("/renice/", func(w http.ResponseWriter, r *http.Request) {
		pid := strings.TrimPrefix(r.URL.Path, "/renice/")
		fmt.Println("Waltz is renicing", pid)
		if n, err := strconv.Atoi(pid); err != nil {
			outErr(w, err)
		} else {
			PID(n).Bump()
			outSuccess(w)
		}
	})
	serv.httpserver = &http.Server{
		Addr:    serv.servingPort(),
		Handler: serv.Mux,
	}
	return serv.httpserver.Serve(serv.netlistener)
}
