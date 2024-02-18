package wzrpc

import "github.com/gofrs/uuid"

type SID string

func RndSID() SID {
	return SID(uuid.Must(uuid.NewV4()).String())
}

func (sid SID) Str() string {
	return string(sid)
}
