package common

// Response is representation of the response payload
type Response struct {
	Message    string      `json:"message,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"statusCode,omitempty"`
	TotalCount int64       `json:"totalCount,omitempty"`
}

type Cockroach struct {
	Host     string `json:"host,omitempty"`
	Port     string `json:"port,omitempty"`
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
}

type Mongo struct {
	Uri    string `json:"uri,omitempty"`
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

var (
	UrlEncoded = "url.encoded"
	UrlCreated = "url.created"
)

type EventDataDetail struct {
	UserId   string `json:"userId,omitempty"`
	LongUrl  string `json:"longUrl,omitempty"`
	ShortUrl string `json:"shortUrl,omitempty"`
}

type EventModel struct {
	AggregateId   string `json:"aggregate_id,omitempty"`
	AggregateType string `json:"aggregate_type,omitempty"`
	EventData     []byte `json:"event_data,omitempty"`
	EventDate     int64  `json:"event_date,omitempty"`
	EventName     string `json:"event_name,omitempty"`
	UserId        string `json:"user_id,omitempty"`
}
