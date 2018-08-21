package fannel

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sm, err := subnet.NewLocalManager(cfg)
	nm, err := network.NewNetworkManager(ctx, sm)
	go func() {
		nm.Run(ctx)
	}()
	<-sigs
	cancel()
}

//network/manager.go
type Manager struct {
	ctx context.Context
	sm  subnet.Manager
	bm  backend.Manager
}

func NewNetworkManager(ctx context.Context, sm subnet.Manager) (*Manager, error) {
	bm := backend.NewManager(ctx, sm, extIface)
	return &Manager{
		ctx: ctx,
		sm:  sm,
		bm:  bm,
	}
}

func (m *Manager) Run(ctx context.Context) {
	m.networks[""] = NewNetwork(ctx, m.sm, m.bm, "", m.ipMasq)
	for _, n := range m.networks {
		wg.Add(1)
		go func(n *Network) {
			m.runNetwork(n)
			wg.Done()
		}(n)
	}
}

func (n *Network) run() {
	n.Config, err = n.sm.GetNetworkConfig(n.ctx, n.Name)
	be, err := n.bm.GetBackend(n.Config.BackendType)
	n.bn, err = be.RegisterNetwork(n.ctx, n.Name, n.Config)

	writeSubnetFile(opts.subnetFile, n.Config.Network, m.ipMasq, bn)

	go func() {
		n.bn.Run(ctx) //---->backend/hostgw/network.go
	}()

	go func() {
		subnet.WatchLease(ctx, n.sm, n.Name, n.bn.Lease().Subnet, evts)
	}()

	dur := n.bn.Lease().Expiration.Sub(time.Now()) - renewMargin
	for {
		select {
		case <-time.After(dur):
			err := n.sm.RenewLease(n.ctx, n.Name, n.bn.Lease())
			dur = n.bn.Lease().Expiration.Sub(time.Now()) - renewMargin
		case e := <-evts:
			switch e.Type {
			case subnet.EventAdded:
				n.bn.Lease().Expiration = e.Lease.Expiration
				dur = n.bn.Lease().Expiration.Sub(time.Now()) - renewMargin
			case subnet.EventRemoved:
				return errInterrupted
			}
		}
	}
}

//backend/manager.go
func NewManager(ctx context.Context, sm subnet.Manager, extIface *ExternalInterface) Manager {
	return &manager{
		ctx:      ctx,
		sm:       sm,
		extIface: extIface,
		active:   make(map[string]Backend),
	}
}

func (bm *manager) GetBackend(backendType string) (Backend, error) {
	betype := strings.ToLower(backendType)
	befunc, ok := backendCtors[betype]
	be, err := befunc(bm.sm, bm.extIface)
	go func() {
		be.Run(bm.ctx)
	}()
	return be, nil
}

//backend/hostgw/hostgw.go
func (be *HostgwBackend) RegisterNetwork(ctx context.Context, netname string, config *subnet.Config) (backend.Network, error) {
	n := &network{
		name:     netname,
		extIface: be.extIface,
		sm:       be.sm,
	}

	attrs := subnet.LeaseAttrs{
		PublicIP:    ip.FromIP(be.extIface.ExtAddr),
		BackendType: "host-gw",
	}

	l, _ := be.sm.AcquireLease(ctx, netname, &attrs)
	n.lease = l
	return n, nil
}

//backend/common.go
type Backend interface {
	Run(ctx context.Context)
	RegisterNetwork(ctx context.Context, network string, config *subnet.Config) (Network, error)
}

type Network interface {
	Lease() *subnet.Lease
	MTU() int
	Run(ctx context.Context)
}

const raceRetries = 10

//use etcd to coordinate
func (m *LocalManager) AcquireLease() (*Lease, error) {
	for i := 0; i < raceRetries; i++ {
		l, err := m.tryAcquireLease(ctx, network, config, attrs.PublicIP, attrs)
		switch err {
		case nil:
			return l, nil
		case errTryAgain:
			continue
		default:
			return nil, err
		}
	}
}

func (m *LocalManager) tryAcquireLease() (*Lease, error) {
	leases, _, err := m.registry.getSubnets(ctx, network)
	sn, err := m.allocateSubnet(config, leases)
	exp, err := m.registry.createSubnet(ctx, network, sn, attrs, subnetTTL)
	switch {
	case err == nil:
		return &Lease{
			Subnet:     sn,
			Attrs:      *attrs,
			Expiration: exp,
		}, nil
	case isErrEtcdNodeExist(err):
		return nil, errTryAgain
	default:
		return nil, err
	}
}

//backend/hostgw/network.go
func (n *network) Run(ctx context.Context) {
	go func() {
		subnet.WatchLeases(ctx, n.sm, n.name, n.lease, evts)
	}()

	n.rl = make([]netlink.Route, 0, 10)
	wg.Add(1)
	go func() {
		n.routeCheck(ctx)
	}()

	for {
		select {
		case evtBatch := <-evts:
			n.handleSubnetEvents(evtBatch)

		case <-ctx.Done():
			return
		}
	}
}

func (n *network) handleSubnetEvents(batch []subnet.Event) {
	for _, evt := range batch {
		switch evt.Type {
		case subnet.EventAdded:
			route := netlink.Route{
				Dst:       evt.Lease.Subnet.ToIPNet(),
				Gw:        evt.Lease.Attrs.PublicIP.ToIP(),
				LinkIndex: n.linkIndex,
			}

			routeList, err := netlink.RouteListFiltered(netlink.FAMILY_V4, &netlink.Route{
				Dst: route.Dst,
			}, netlink.RT_FILTER_DST)
			if len(routeList) > 0 && !routeList[0].Gw.Equal(route.Gw) {
				if err := netlink.RouteDel(&route); err != nil {
					log.Errorf("Error deleting route to %v: %v", evt.Lease.Subnet, err)
					continue
				}
			}
			n.addToRouteList(route)
		case subnet.EventRemoved:
			route := netlink.Route{
				Dst:       evt.Lease.Subnet.ToIPNet(),
				Gw:        evt.Lease.Attrs.PublicIP.ToIP(),
				LinkIndex: n.linkIndex,
			}
			netlink.RouteDel(&route)
			n.removeFromRouteList(route)

		default:
			log.Error("Internal error: unknown event type: ", int(evt.Type))
		}
	}
}
