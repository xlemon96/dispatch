package dispatch

type Dispatch interface {
	Start() error
	GetTodoDag() chan *dagBag
}
