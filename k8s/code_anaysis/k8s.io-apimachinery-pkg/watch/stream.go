package watch

type Decoder interface {
	Decode() (EventType, runtime.Object, error)
	Close()
}

type StreamWatcher struct {
	sync.Mutex
	source  Decoder
	result  chan Event
	stopped bool
}

func (sw *StreamWatcher) Stop() {
	sw.Lock()
	defer sw.Unlock()
	if !sw.stopped {
		sw.stopped = true
		sw.source.Close()
	}
}

// receive reads result from the decoder in a loop and sends down the result channel.
func (sw *StreamWatcher) receive() {
	defer close(sw.result)
	defer sw.Stop()
	defer utilruntime.HandleCrash()
	for {
		action, obj, err := sw.source.Decode()
		if err != nil {
			return
		}
		sw.result <- Event{
			Type:   action,
			Object: obj,
		}
	}
}
