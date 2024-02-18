package wzrpc

import (
	"encoding/json"
	"errors"
)

type Answer struct {
	SID    SID    `json:"category"`
	Result any    `json:"result"`
	Error  string `json:"error"`
}

func (a Answer) GetError() error {
	if msg := a.Error; msg != "" {
		return errors.New(msg)
	}
	return nil
}

func (a Answer) DeployTo(dest any) error {
	if j, err := json.Marshal(a.Result); err != nil {
		return err
	} else {
		return json.Unmarshal(j, &dest)
	}
}
