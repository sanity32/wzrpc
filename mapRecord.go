package wzrpc

import "time"

type MapRecord[T any] struct {
	Reg     time.Time
	Content T
}

func (rec *MapRecord[T]) Read() T {
	return rec.Content
}

func (rec *MapRecord[T]) IsExpired(ttl time.Duration) bool {
	expTime := rec.Reg.Add(ttl)
	return time.Now().After(expTime)
}
