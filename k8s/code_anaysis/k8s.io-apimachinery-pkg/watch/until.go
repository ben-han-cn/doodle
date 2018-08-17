package watch

type ConditionFunc func(event Event) (bool, error)

//only after previous condition meet, the following conditions will
//be checked, the event make last condition meet, will be checked against
//next condtion
func Until(timeout time.Duration, watcher Interface, conditions ...ConditionFunc) (*Event, error) {
	ch := watcher.ResultChan()
	defer watcher.Stop()
	var after <-chan time.Time
	if timeout > 0 {
		after = time.After(timeout)
	} else {
		ch := make(chan struct{})
		defer close(ch)
		after = ch
	}

	var lastEvent *Event
	for _, condition := range conditions {
		// check the next condition against the previous event and short circuit waiting for the next watch
		if lastEvent != nil {
			done, err := condition(*lastEvent)
			if err != nil {
				return lastEvent, err
			}
			if done {
				continue
			}
		}
	ConditionSucceeded:
		for {
			select {
			case event, ok := <-ch:
				if !ok {
					return lastEvent, ErrWatchClosed
				}
				lastEvent = &event

				// TODO: check for watch expired error and retry watch from latest point?
				done, err := condition(event)
				if err != nil {
					return lastEvent, err
				}
				if done {
					break ConditionSucceeded
				}

			case <-after:
				return lastEvent, wait.ErrWaitTimeout
			}
		}
	}
	return lastEvent, nil
}
