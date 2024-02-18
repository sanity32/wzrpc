package wzrpc

// import (
// 	"fmt"
// 	"time"
// )

// type Asker interface{ ChAsk() chan Ask }
// type Aliver interface{ AliveMap() *Map[Alive] }
// type Answerer interface{ AnswerMap() *Map[chan Answer] }
// type Processeder interface{ ProcessedMap() *Map[bool] }

// type AskerAliverAnswererProcesseder interface {
// 	Asker
// 	Aliver
// 	Answerer
// 	Processeder
// }

// func Allocate(w AskerAliverAnswererProcesseder, sid SID) error {
// 	w.AliveMap().Allocate(sid)
// 	w.AnswerMap().Allocate(sid)
// 	return w.ProcessedMap().Delete(sid)
// }

// func Delete(w AskerAliverAnswererProcesseder, sid SID) error {
// 	if err := w.AliveMap().Delete(sid); err != nil {
// 		return err
// 	}
// 	if err := w.AnswerMap().Delete(sid); err != nil {
// 		return err
// 	}
// 	w.ProcessedMap().Write(sid, true)
// 	return nil
// }

// // Follower's part
// func PullAsk(w Asker, timeout time.Duration) (ask Ask, err error) {
// 	watchdog := time.NewTimer(timeout)
// 	defer watchdog.Stop()
// 	select {
// 	case <-watchdog.C:
// 		return ask, ErrAskTimeout
// 	case r := <-w.ChAsk():
// 		return r, nil
// 	}
// }

// // Leader's part
// func PushAsk(w AskerAliverAnswererProcesseder, ask Ask, timeout time.Duration) (err error) {
// 	watchdog := time.NewTimer(timeout)
// 	defer watchdog.Stop()
// 	Allocate(w, ask.Response)

// 	defer func() {
// 		if recover() != nil {
// 			err = ErrAskChanClosed
// 		}
// 	}()

// 	select {
// 	case <-watchdog.C:
// 		return ErrRequestTimeout
// 	case w.ChAsk() <- ask:
// 		return
// 	}

// }

// // Error is always nil by yet
// func PushAlive(w Aliver, data AliveDTO) error {
// 	w.AliveMap().Write(data.SID, data.Norm())
// 	return nil
// }

// // Not much to be public method, used by PullAnswer
// func AwaitAlive(w Aliver, sid SID, ttl, timeout, interval time.Duration) <-chan any {
// 	var ch = make(chan any)
// 	stop := false
// 	lastAlive := func(w Aliver, sid SID) time.Time {
// 		if data, found := w.AliveMap().Read(sid); found && data.Timestamp != 0 {
// 			return time.Unix(data.Timestamp, 0)
// 		}
// 		w.AliveMap().Write(sid, Alive{Timestamp: time.Now().Unix()})
// 		return time.Now()
// 	}

// 	var watchdog = time.AfterFunc(timeout, func() { stop = true })
// 	go func() {
// 		for {
// 			expTime := lastAlive(w, sid).Add(ttl)
// 			if stop || time.Now().After(expTime) {
// 				break
// 			}
// 			time.Sleep(interval)
// 		}
// 		watchdog.Stop()
// 		ch <- nil
// 	}()
// 	return ch
// }

// // Follower part
// func PushAnswer(w Answerer, v Answer) (err error) {
// 	defer func() {
// 		if recover() != nil {
// 			fmt.Println("Unable to PUSH answer. Channel is appeared to be closed.")
// 			err = ErrAnswerChanClosed
// 		}
// 	}()
// 	w.AnswerMap().Allocate(v.SID) <- v
// 	return
// }

// // Leader part.
// func PullAnswer(w AskerAliverAnswererProcesseder, sid SID, aliveTTL, timeout, interval time.Duration) (r Answer, err error) {
// 	ch, ok := w.AnswerMap().Read(sid)
// 	if !ok {
// 		fmt.Println("Unable to read answer for SID", sid)
// 		return
// 	}
// 	defer Delete(w, sid)
// 	select {
// 	case <-AwaitAlive(w, sid, aliveTTL, timeout, interval):
// 		return r, ErrSidNotAlive
// 	case r := <-ch:
// 		return r, nil
// 	}
// }
