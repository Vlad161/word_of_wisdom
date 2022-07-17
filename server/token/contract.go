package token

type Storage interface {
	Get(k string) (interface{}, error)
	Put(k string, v interface{}) error
	Delete(k string) error
}
