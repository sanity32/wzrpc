package wzrpc

type Alive struct {
	Timestamp int64
	Data      any
}

type AliveDTO struct {
	SID       SID   `json:"category"`
	Ok        bool  `json:"ok"`
	Timestamp int64 `json:"timestamp"`
	Data      any   `json:"data,omitempty"`
}

func (a AliveDTO) Norm() Alive {
	return Alive{
		Timestamp: a.Timestamp,
		Data:      a.Data,
	}
}
