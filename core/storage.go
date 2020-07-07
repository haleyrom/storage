package core

// Storage Storage
type Storage interface {
	Initialize(config ...string) error
}
