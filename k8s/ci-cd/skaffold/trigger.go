package trigger

type Trigger interface {
	Start() (<-chan bool, func()) //return implus and clean up function
	WatchForChanges(io.Writer)    //for debug
	Debounce() bool               //whether to omit too often changes
}

//pollTrigger  based on interval
//manualTrigger based on use import

type Watcher interface {
	//when specified file changes, call onChange
	Register(deps func() ([]string, error), onChange func(Events)) error
	//run file check based on trigger event, when file changed
	//invoke onChange from Register when related file changed
	//invoke onChange if any file changed
	Run(ctx context.Context, trigger Trigger, onChange func() error) error
}
