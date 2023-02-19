package common

type Cockroach struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Mongo struct {
	Url    string
	DbName string
}

type Kafka struct {
	Url   string
	Topic string
}

type Pulsar struct {
	Url   string
	Topic string
}
