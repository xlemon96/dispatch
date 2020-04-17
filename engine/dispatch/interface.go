package dispatch

type Dispatch interface {
	Start() error
	GetTodoDags() chan *dagBag
}
