package wzrpc

import "errors"

var (
	ErrUnableToDeployData = errors.New("unable to deploy data")
	ErrSidNotAlive        = errors.New("session not alive")
	ErrAskTimeout         = errors.New("timed out while awaiting an ask")
	ErrLeaderAskFailed    = errors.New("leader failed to ask")
	ErrAskChanClosed      = errors.New("ask channel is closed")
	ErrAnswerChanClosed   = errors.New("answer channel is closed")
	ErrRequestIsNil       = errors.New("provided request is a nil")
	ErrRequestTimeout     = errors.New("request timeout")
	ErrMaxRecursion       = errors.New("max action recursion level reached")
)
