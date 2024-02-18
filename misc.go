package wzrpc

import "encoding/json"

const MAX_LEADER_ACTION_RECURSION_LEVEL = 9

type LeaderAnswerRequest struct {
	SID SID `json:"category"`
}

type LeaderAskResponse struct {
	Success  bool   `json:"success"`
	ErrorMsg string `json:"error"`
}

func deployDataTo(src any, dest any) error {
	j, err := json.Marshal(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, &dest)
}

func deployData[T any](src any) (dest T) {
	deployDataTo(src, &dest)
	return
}
