package wzrpc

type ServerSettings struct {
	Host                   string
	Port                   int
	LeaderAskEndpoint      string
	LeaderAnswerEndpoint   string
	FollowerAskEndpoint    string
	FollowerAliveEndpoint  string
	FollowerAnswerEndpoint string
}

func NewServerSettings(host string, port int) ServerSettings {
	return ServerSettings{
		Host:                   host,
		Port:                   port,
		LeaderAskEndpoint:      "/leader/ask",
		LeaderAnswerEndpoint:   "/leader/answer",
		FollowerAskEndpoint:    "/follower/ask",
		FollowerAliveEndpoint:  "/follower/alive",
		FollowerAnswerEndpoint: "/follower/answer",
	}
}
