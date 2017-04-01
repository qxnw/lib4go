package registry

type ValueWatcher interface {
	GetValue() []byte
	GetError() error
}
type ChildrenWatcher interface {
	GetValue() []string
	GetError() error
}
