control_listener
    stack
    .push(control::client::layer())
    .push(control::resolve::layer(dns_resolver.clone()))
    .push(reconnect::layer().with_fixed_backoff(config.control_backoff_delay))
    .push(proxy::timeout::layer(config.control_connect_timeout))
    .push(http_metrics::layer::<_, classify::Response>( ctl_http_metrics,))
    .push(control::grpc_request_payload::layer())
    .push(svc::watch::layer(tls_client_config.clone()))
    .push(phantom_data::layer())
    .push(control::add_origin::layer())
    .push(buffer::layer())
    .push(limit::layer(config.destination_concurrency_limit));


outbound
    let connect = connect::Stack::new()
                    .push(proxy::timeout::layer(config.outbound_connect_timeout))
                    .push(transport_metrics.connect("outbound"));
    let client_stack = connect.clone().push()...
    let endpoint_stack = client_stack.push()...
    let dst_stack = endpoint_stack.push()...
    let dst_router = dst_stack.push()...
    let addr_stack = dst_router.push()...
    let addr_router = addr_stack.push()...
    let server_stack = addr_router.push()...


inbount
    let connect = connect::Stack::new()
                    .push(proxy::timeout::layer(config.inbound_connect_timeout))
                    .push(transport_metrics.connect("inbound"))
                    .push(rewrite_loopback_addr::layer());
    let client_stack = connect .clone().push()...
    let endpoint_router = client_stack.push()...
    let dst_stack = endpoint_router.push()...
    let dst_router = dst_stack.push()...
    let source_stack = dst_router.push()...

//start tap and control listener
thread::Builder::new()
                .name("admin".into())
                .spawn(move || {

                    let mut rt = current_thread::Runtime::new().expect("initialize admin thread runtime");

                    let metrics = control::serve_http(
                        "metrics",
                        metrics_listener,
                        metrics::Serve::new(report),
                    );

                    rt.spawn(tap_daemon.map_err(|_| ()));
                    rt.spawn(serve_tap(control_listener, TapServer::new(tap_grpc)));
                    rt.spawn(metrics);
                    rt.spawn(::logging::admin().bg("dns-resolver").future(dns_bg));
                    rt.spawn( ::logging::admin() .bg("resolver") .future(resolver_bg_rx.map_err(|_| {}).flatten()),);
                    rt.spawn(::logging::admin().bg("tls-config").future(tls_cfg_bg));
                    let shutdown = admin_shutdown_signal.then(|_| Ok::<(), ()>(()));
                    rt.block_on(shutdown).expect("admin");
                    trace!("admin shutdown finished");
                })


//start inbound and outbound
main_fut = inbound.join(outbind).map(|_| {})
runtime.spawn(Box::new(main_fut));



inbound_listener
outbound_listener
metrics_listener
tap
