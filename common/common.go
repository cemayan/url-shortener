package common

type Cockroach struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Mongo struct {
	Url    string `json:"url,omitempty"`
	DbName string `json:"dbName,omitempty"`
}

type Kafka struct {
	Url   string `json:"url,omitempty"`
	Topic string `json:"topic,omitempty"`
}

type Pulsar struct {
	Url   string `json:"url,omitempty"`
	Topic string `json:"topic,omitempty"`
}
