package output

type RedisPort interface {
	Get(key string) (string, error)
	Set(key string, data string) error
	Remove(key string) error
}
