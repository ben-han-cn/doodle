package evnet

//k8s.io/api/core/v1/types.go
type EventSource struct {
	Component string
	Host      string
}

type Event struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	InvolvedObject      ObjectReference
	Reason              string
	Message             string
	Source              EventSource
	FirstTimestamp      metav1.Time
	LastTimestamp       metav1.Time
	Count               int32
	Type                string
	EventTime           metav1.MicroTime
	Series              *EventSeries
	Action              string
	Related             *ObjectReference
	ReportingController string
	ReportingInstance   string
}

//event.go
// EventSink knows how to store events (client.Client implements it.)
type EventSink interface {
	Create(event *v1.Event) (*v1.Event, error)
	Update(event *v1.Event) (*v1.Event, error)
	Patch(oldEvent *v1.Event, data []byte) (*v1.Event, error)
}

// EventRecorder knows how to record events on behalf of an EventSource.
type EventRecorder interface {
	Event(object runtime.Object, eventtype, reason, message string)
	Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{})
	PastEventf(object runtime.Object, timestamp metav1.Time, eventtype, reason, messageFmt string, args ...interface{})
	AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{})
}

// EventBroadcaster knows how to receive events and send them to any EventSink, watcher, or log.
type EventBroadcaster interface {
	StartEventWatcher(eventHandler func(*v1.Event)) watch.Interface
	StartRecordingToSink(sink EventSink) watch.Interface
	StartLogging(logf func(format string, args ...interface{})) watch.Interface
	// NewRecorder returns an EventRecorder that can be used to send events to this EventBroadcaster
	// with the event source set to the given event source.
	NewRecorder(scheme *runtime.Scheme, source v1.EventSource) EventRecorder
}

type EventFilterFunc func(event *v1.Event) bool

type EventCorrelator struct {
	filterFunc EventFilterFunc
	aggregator *EventAggregator
	logger     *eventLogger
}

type EventCorrelateResult struct {
	Event *v1.Event
	Patch []byte
	Skip  bool
}

type EventAggregator struct {
	sync.RWMutex
	cache                *lru.Cache
	keyFunc              EventAggregatorKeyFunc
	messageFunc          EventAggregatorMessageFunc
	maxEvents            uint
	maxIntervalInSeconds uint
	clock                clock.Clock
}

// EventAggregatorKeyFunc is responsible for grouping events for aggregation
// It returns a tuple of the following:
// aggregateKey - key the identifies the aggregate group to bucket this event
// localKey - key that makes this event in the local group
type EventAggregatorKeyFunc func(event *v1.Event) (aggregateKey string, localKey string)

// EventAggregatorMessageFunc is responsible for producing an aggregation message
type EventAggregatorMessageFunc func(event *v1.Event) string

/*
generate event module use EventRecorder to create event
use EventBroadcaster to broadcast event to api server, log file and other event sinker
before send event to api server, EventCorrelator will aggregate event, based on the result,
make the decision whether to skip event or to update event(has related events) or post
new events to etcd
*/
