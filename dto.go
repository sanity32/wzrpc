package wzrpc

import (
	"fmt"
	"time"
)

type Registry struct {
	ask       chan Ask
	alive     Map[Alive]
	answer    Map[chan Answer]
	processed Map[bool]
}

func (w *Registry) ChAsk() chan Ask {
	return w.ask
}

func (w *Registry) AliveMap() *Map[Alive] {
	return &w.alive
}

func (w *Registry) AnswerMap() *Map[chan Answer] {
	return &w.answer
}

func (w *Registry) ProcessedMap() *Map[bool] {
	return &w.processed
}

func (g *Registry) Allocate(sid SID) error {
	g.alive.Allocate(sid)
	g.answer.Allocate(sid)
	return g.processed.Delete(sid)
}

func (g *Registry) Delete(sid SID) error {
	if err := g.alive.Delete(sid); err != nil {
		return err
	}
	if err := g.answer.Delete(sid); err != nil {
		return err
	}
	g.processed.Write(sid, true)
	return nil
}

// Follower's part
func (g *Registry) PullAsk(timeout time.Duration) (ask Ask, err error) {
	watchdog := time.NewTimer(timeout)
	defer watchdog.Stop()
	select {
	case <-watchdog.C:
		return ask, ErrAskTimeout
	case r := <-g.ask:
		return r, nil
	}
}

// Leader's part
func (g *Registry) PushAsk(ask Ask, timeout time.Duration) (err error) {
	watchdog := time.NewTimer(timeout)
	defer watchdog.Stop()
	g.Allocate(ask.Response)

	defer func() {
		if recover() != nil {
			err = ErrAskChanClosed
		}
	}()

	select {
	case <-watchdog.C:
		return ErrRequestTimeout
	case g.ask <- ask:
		return
	}

}

func (g *Registry) PushAlive(dto AliveDTO) {
	g.alive.Write(dto.SID, dto.Norm())
}

func (g *Registry) AwaitAlive(sid SID, ttl, timeout, interval time.Duration) <-chan any {
	var ch = make(chan any)
	stop := false
	lastAlive := func(sid SID) time.Time {
		if data, found := g.alive.Read(sid); found && data.Timestamp != 0 {
			return time.Unix(data.Timestamp, 0)
		}
		g.alive.Write(sid, Alive{Timestamp: time.Now().Unix()})
		return time.Now()
	}

	var watchdog = time.AfterFunc(timeout, func() { stop = true })
	go func() {
		for {
			expTime := lastAlive(sid).Add(ttl)
			if stop || time.Now().After(expTime) {
				break
			}
			time.Sleep(interval)
		}
		watchdog.Stop()
		ch <- nil
	}()
	return ch
}

// Follower part
func (g *Registry) PushAnswer(v Answer) (err error) {
	defer func() {
		if recover() != nil {
			fmt.Println("Unable to PUSH answer. Channel is appeared to be closed.")
			err = ErrAnswerChanClosed
		}
	}()
	g.answer.Allocate(v.SID) <- v
	return
}

// Leader part.
func (g *Registry) PullAnswer(sid SID, aliveTTL, timeout, interval time.Duration) (r Answer, err error) {
	ch, ok := g.answer.Read(sid)
	if !ok {
		fmt.Println("Unable to read answer for SID", sid)
		return
	}
	defer g.Delete(sid)
	select {
	case <-g.AwaitAlive(sid, aliveTTL, timeout, interval):
		return r, ErrSidNotAlive
	case r := <-ch:
		return r, nil
	}
}
