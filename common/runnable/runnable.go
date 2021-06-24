package runnable

type Runnable interface {
	Init() error
	Run() error
	Shutdown()
}
