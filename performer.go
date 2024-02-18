package wzrpc

type Leader interface {
	LeaderAsk(ask Ask) error
	LeaderAnswer(k SID) (Answer, error)
}

type Follower interface {
	FollowerAsk() (Ask, error)
	FollowerAlive(AliveDTO)
	FollowerAnswer(Answer) error
}

type Performer interface {
	Leader
	Follower
}

type Performerer interface {
	WaltzPerformer() Performer
}

func LeaderAction(l Leader, command string, opts any, recur int) (resp any, err error) {
	ask := NewAsk(command).SetOptions(opts)
	if err = l.LeaderAsk(*ask); err != nil {
		return resp, ErrLeaderAskFailed
	}
	answer, err := l.LeaderAnswer(ask.Response)
	if err == ErrSidNotAlive {
		if recur < MAX_LEADER_ACTION_RECURSION_LEVEL {
			return LeaderAction(l, command, opts, recur+1)
		}
		return resp, ErrMaxRecursion
	}
	return answer.Result, answer.GetError()
}
