package wzrpc

type Ask struct {
	Command  string         `json:"command"`
	Options  map[string]any `json:"options"`
	Response SID            `json:"responseTo"`
}

func NewAsk(command string, responseTo ...SID) *Ask {
	var sid SID
	switch a := responseTo; len(a) {
	case 0:
		sid = RndSID()
	default:
		sid = a[0]
	}
	return &Ask{
		Command:  command,
		Response: sid,
	}
}

func (req *Ask) SetOption(k string, v any) *Ask {
	if req.Options == nil {
		req.Options = map[string]any{}
	}
	req.Options[k] = v
	return req
}

func (req *Ask) SetOptions(data any) *Ask {
	m := deployData[map[string]any](data)
	for k, v := range m {
		req.SetOption(k, v)
	}
	return req
}
