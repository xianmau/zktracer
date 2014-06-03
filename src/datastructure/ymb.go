package datastructure

type YMB struct {
	Brokers     []BrokerInfo
	Loggers     []LoggerInfo
	Apps        []AppInfo
	Topics      []TopicInfo
	ZoneId      string
	RemoteZones []string

	BrokersTopics []BrokerTopic
}
