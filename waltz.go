package wzrpc

import "time"

type Settings struct {
	AskLeaderRequestTimeout   time.Duration
	AskFollowerStandbyTimeout time.Duration
	AlivePollingInterval      time.Duration
	AlivePollingTimeout       time.Duration
	AliveTTL                  time.Duration
}

func NewModel() *Model {
	return &Model{
		registry: NewRegistry(),
		Settings: *DefaultSettings(),
	}
}

type Model struct {
	registry *Registry
	Settings Settings
}

func (m *Model) LeaderAsk(ask Ask) error {
	return m.registry.PushAsk(ask, m.Settings.AskLeaderRequestTimeout)
	// return PushAsk(m.data, ask, m.Settings.AskLeaderRequestTimeout)
}

func (m *Model) LeaderAnswer(k SID) (Answer, error) {
	return m.registry.PullAnswer(k,
		m.Settings.AliveTTL,
		m.Settings.AlivePollingTimeout,
		m.Settings.AlivePollingInterval,
	)
}

func (m *Model) FollowerAsk() (Ask, error) {
	return m.registry.PullAsk(m.Settings.AskFollowerStandbyTimeout)
}

func (m *Model) FollowerAlive(followerMsg AliveDTO) {
	m.registry.PushAlive(followerMsg)
}

func (m *Model) FollowerAnswer(a Answer) error {
	return m.registry.PushAnswer(a)
}

func DefaultSettings() *Settings {
	return &Settings{
		AskLeaderRequestTimeout:   time.Second * 30,
		AskFollowerStandbyTimeout: time.Second * 30,
		AlivePollingInterval:      time.Second * 3,
		AlivePollingTimeout:       time.Second * 60,
		AliveTTL:                  time.Second * 7,
	}
}
