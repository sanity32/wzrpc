package wzrpc

import (
	"sync"
	"time"
)

const DEFAULT_RECORD_TTD = time.Second * 60

type Map[T any] struct {
	content          map[SID]MapRecord[T]
	contentInitiator sync.Once
	contentLock      sync.Mutex
	gcLock           sync.Mutex
	gcLastRun        time.Time
	gcInterval       time.Duration // default is 6s
	recordTTL        time.Duration // default is 60s

	RecordDestructor  func(MapRecord[T]) error
	RecordConstructor func() MapRecord[T]
}

func (cc *Map[T]) init() {
	if cc.content == nil {
		cc.content = map[SID]MapRecord[T]{}
	}
}

func (cc *Map[T]) deflock() func() {
	cc.contentLock.Lock()
	return func() {
		cc.contentLock.Unlock()
	}
}

func (cc *Map[T]) RecordTTL() time.Duration {
	if cc.recordTTL == 0 {
		cc.recordTTL = DEFAULT_RECORD_TTD
	}
	return cc.recordTTL
}

func (cc *Map[T]) Content() map[SID]MapRecord[T] {
	cc.contentInitiator.Do(cc.init)
	return cc.content
}

func (cc *Map[T]) getExpiredRecs() (rr []SID) {
	defer cc.deflock()()
	for k, record := range cc.Content() {
		if record.IsExpired(cc.RecordTTL()) {
			rr = append(rr, k)
		}
	}
	return
}

func (cc *Map[T]) createEmptyRecord() MapRecord[T] {
	var result MapRecord[T]
	if fn := cc.RecordConstructor; fn != nil {
		result = fn()
	}
	result.Reg = time.Now()
	return result
}

func (cc *Map[T]) Allocate(k SID) T {
	cc.GC()
	defer cc.deflock()()

	if _, found := cc.Content()[k]; !found {
		cc.Content()[k] = cc.createEmptyRecord()
	}
	return cc.Content()[k].Content
}

func (cc *Map[T]) Read(k SID) (value T, found bool) {
	defer cc.deflock()()
	if record, found := cc.Content()[k]; found {
		return record.Content, true
	}
	return
}

func (cc *Map[T]) Write(k SID, content T) {
	cc.Allocate(k)
	defer cc.deflock()()

	record := cc.Content()[k]
	record.Content = content
	cc.Content()[k] = record
}

func (cc *Map[T]) Delete(k SID) error {
	defer cc.deflock()()
	return cc.remove(k)
}

func (cc *Map[T]) remove(k SID) error {
	if fn := cc.RecordDestructor; fn != nil {
		if err := fn(cc.Content()[k]); err != nil {
			return err
		}
	}
	delete(cc.Content(), k)
	return nil
}

func (cc *Map[T]) GC() {
	if cc.gcLastRun.IsZero() {
		cc.gcLastRun = time.Now()
		return
	}
	if cc.gcInterval == 0 {
		cc.gcInterval = time.Second * 6
	}

	if nextDue := cc.gcLastRun.Add(cc.gcInterval); time.Now().Before(nextDue) {
		return
	}

	if !cc.gcLock.TryLock() {
		return
	}
	defer cc.gcLock.Unlock()

	cc.gcLastRun = time.Now()
	if expiredRecords := cc.getExpiredRecs(); expiredRecords != nil {
		defer cc.deflock()()
		for _, k := range expiredRecords {
			cc.remove(k)
		}
	}
}
