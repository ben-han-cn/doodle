package watch

//k8s.io/apiserver/pkg/endpoints/installer.go
func (a *APIInstaller) Install() ([]metav1.APIResource, *restful.WebService, []error) {
  a.registerResourceHandlers(path, a.group.Storage[path], ws)
}

func (a *APIInstaller) registerResourceHandlers(path string, storage rest.Storage, ws *restful.WebService) (*metav1.APIResource, error) {
  switch action.Verb {
    case "LIST": // List all resources of a kind.
    doc := "list objects of kind " + kind
    if hasSubresource {
      doc = "list " + subresource + " of objects of kind " + kind
    }

    watcher, isWatcher := storage.(rest.Watcher)
    handler := metrics.InstrumentRouteFunc(action.Verb, resource, subresource, requestScope, restfulListResource(lister, watcher, reqScope, false, a.minRequestTimeout))
  }
}


//restfulListResource -> handlers.ListResource
//k8s.io/apiserver/pkg/endpoints/handlers/get.go
func ListResource(r rest.Lister, rw rest.Watcher, scope RequestScope, forceWatch bool, minRequestTimeout time.Duration) http.HandlerFunc {
   return func(w http.ResponseWriter, req *http.Request) {
     opts := metainternalversion.ListOptions{}
     metainternalversion.ParameterCodec.DecodeParameters(req.URL.Query(), scope.MetaGroupVersion, &opts)
     if opts.Watch || forceWatch {
       watcher, err := rw.Watch(ctx, &opts)
       metrics.RecordLongRunning(req, requestInfo, func() {
         serveWatch(watcher, scope, req, w, timeout)
       })
     }
   }
}


//k8s.io/apiserver/pkg/endpoints/handlers/watch.go
func serveWatch(watcher watch.Interface, scope RequestScope, req *http.Request, w http.ResponseWriter, timeout time.Duration) {
  server := &WatchServer{}
  server.ServeHTTP(w, req)
}

func (s *WatchServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
  if wsstream.IsWebSocketRequest(req) {
    w.Header().Set("Content-Type", s.MediaType)
    websocket.Handler(s.HandleWS).ServeHTTP(w, req)
    return
  }

  flusher, ok := w.(http.Flusher)
  framer := s.Framer.NewFrameWriter(w)
  e := streaming.NewEncoder(framer, s.Encoder)

  ch := s.Watching.ResultChan()
  for {
    select {
    case <-cn.CloseNotify():
      return
    case <-timeoutCh:
      return
    case event, ok := <-ch:
      e.Encode(outEvent)
      if len(ch) == 0 {
        flusher.Flush()
      }
      buf.Reset()
    }
  }
}


//k8s.io/apiserver/pkg/registry/rest/rest.go
type Watcher interface {
  Watch(ctx genericapirequest.Context, options *metainternalversion.ListOptions) (watch.Interface, error)
}

//k8s.io/apiserver/pkg/storage/etcd3/watcher.go
func (w *watcher) Watch(ctx context.Context, key string, rev int64, recursive bool, pred storage.SelectionPredicate) (watch.Interface, error) {
    if recursive && !strings.HasSuffix(key, "/") {
        key += "/"
    }
    wc := w.createWatchChan(ctx, key, rev, recursive, pred)
    go wc.run()
    return wc, nil
}

func (wc *watchChan) run() {
    watchClosedCh := make(chan struct{})
    go wc.startWatching(watchClosedCh)

    var resultChanWG sync.WaitGroup
    resultChanWG.Add(1)
    go wc.processEvent(&resultChanWG)

    select {
    case err := <-wc.errChan:
        if err == context.Canceled {
            break
        }
        errResult := transformErrorToEvent(err)
        if errResult != nil {
            // error result is guaranteed to be received by user before closing ResultChan.
            select {
            case wc.resultChan <- *errResult:
            case <-wc.ctx.Done(): // user has given up all results
            }
        }
    case <-watchClosedCh:
    case <-wc.ctx.Done(): // user cancel
    }

    // We use wc.ctx to reap all goroutines. Under whatever condition, we should stop them all.
    // It's fine to double cancel.
    wc.cancel()

    // we need to wait until resultChan wouldn't be used anymore
    resultChanWG.Wait()
    close(wc.resultChan)
}
