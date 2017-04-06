package registry

type ValueWatcher interface {
	GetValue() ([]byte, int32)
	GetError() error
}
type ChildrenWatcher interface {
	GetValue() ([]string, int32)
	GetError() error
}
