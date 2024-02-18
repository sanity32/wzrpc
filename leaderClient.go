package wzrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func NewLeaderClient(host string, port int) *LeaderClient {
	return &LeaderClient{Settings: NewServerSettings(host, port)}
}

type LeaderClient struct {
	Settings ServerSettings
}

func (l *LeaderClient) Performer() Leader {
	return l
}

func (y *LeaderClient) url(suffix string) string {
	s, _ := strings.CutPrefix(suffix, "/")
	return fmt.Sprintf("http://%v:%v/%v", y.Settings.Host, y.Settings.Port, s)
}

func (y *LeaderClient) req(url string, data any) (r []byte, err error) {
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(data); err != nil {
		return r, err
	}

	req, err := http.NewRequest(http.MethodPost, url, &buff)
	if err != nil {
		return r, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: time.Second * 60}
	res, err := client.Do(req)
	if err != nil {
		return r, err
	}
	return io.ReadAll(res.Body)
}

func (y *LeaderClient) LeaderAsk(ask Ask) error {
	url := y.url(y.Settings.LeaderAskEndpoint)
	_, err := y.req(url, ask)
	return err
}

func (y *LeaderClient) LeaderAnswer(k SID) (r Answer, err error) {
	url := y.url(y.Settings.LeaderAnswerEndpoint)
	req := LeaderAnswerRequest{SID: k}
	if body, err := y.req(url, req); err == nil {
		return r, json.Unmarshal(body, &r)
	}
	return
}
