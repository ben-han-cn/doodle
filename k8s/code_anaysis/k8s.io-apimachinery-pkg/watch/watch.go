package watch

//k8s.io/apimachinery/pkg/watch

type Event struct {
	Type   EventType //added, modfied, deleted, error
	Object runtime.Object
}

type Interface interface {
	Stop()
	ResultChan() <-chan Event
}

type FilterFunc func(in Event) (out Evnet, keep bool)

func Filter(w Interface, f FilterFunc) Interface {
	fw := &filteredWatch{
		inner:  w,
		result: make(chan Event),
		f:      f,
	}
	go fw.loop()
	return fw
}

func (fw *filteredWatch) loop() {
	defer close(fw.result)
	for {
		e, ok := <-fw.inner.ResultChan()
		if !ok {
			break
		}
		ne, keep := fw.f(e)
		if keep {
			fw.result <- ne
		}
	}
}

type Recorder struct {
	Interface

	lock   sync.Mutex
	events []Event
}

func NewRecorder(w Interface) *Recorder {
	r := &Recorder{}
	r.Interface = Filter(w, r.record)
	return f
}

func (r *Recorder) record(in Event) (Event, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.events = append(r.events, in)
	return in, true
}
