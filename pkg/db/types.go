package db

type Datastore interface {
	Close() error
	GetType() string
}
