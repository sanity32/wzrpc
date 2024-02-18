package wzrpc

import "errors"

func NewRegistry() *Registry {
	return &Registry{
		ask:    make(chan Ask),
		answer: *NewAnswerChanMap(),
	}
}

func NewAnswerChanMap() *Map[chan Answer] {
	r := Map[chan Answer]{}
	r.content = map[SID]MapRecord[chan Answer]{}
	r.RecordConstructor = func() MapRecord[chan Answer] {
		return MapRecord[chan Answer]{
			Content: make(chan Answer),
		}
	}
	r.RecordDestructor = func(mr MapRecord[chan Answer]) (err error) {
		defer func() {
			if recover() != nil {
				err = errors.New("unable to close the channel")
			}
		}()
		close(mr.Read())
		return
	}
	return &r
}
