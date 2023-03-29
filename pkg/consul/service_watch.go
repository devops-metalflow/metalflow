package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"metalflow/pkg/global"
)

type Service struct {
	ServiceID   string `json:"serviceId"`
	ServiceName string `json:"serviceName"`
	Address     string `json:"address"`
	Port        int    `json:"port"`
	Status      string `json:"status"`
	ServerOs    string `json:"serverOs"`
}

type Registry struct {
	Addr         string
	Client       *consulapi.Client
	StatusChan   chan *Service
	ShutdownChan chan string
	watchers     map[string]*watch.Plan
}

func NewRegistry(addr string) (*Registry, error) {
	config := consulapi.DefaultConfig()
	config.Address = addr
	c, err := consulapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Registry{
		Addr:         addr,
		Client:       c,
		StatusChan:   make(chan *Service),
		ShutdownChan: make(chan string),
		watchers:     make(map[string]*watch.Plan),
	}, nil
}

func (r *Registry) StartWatch() error {
	wp, err := newWatchPlan("services", nil, r.getWholeSvcHandler())
	if err != nil {
		fmt.Printf("new watch plan failed: %v\n", err)
		return err
	}
	go func() {
		err = runWatchPlan(wp, r.Addr)
		if err != nil {
			fmt.Println("run consul watch failed: ", err)
			return
		}
	}()
	return err
}

// getWholeSvcHandler uses to get whole service watch
func (r *Registry) getWholeSvcHandler() Handler {
	return func(_ uint64, data any) {
		switch d := data.(type) {
		// "services" watch type returns map[string][]string type. follow:https://www.consul.io/docs/dynamic-app-config/watches#services
		case map[string][]string:
			// watchers的更改始终在同一个goroutine内，因此是安全的
			for k := range d {
				if _, ok := r.watchers[k]; ok || k == "consul" {
					continue
				}
				// 如果没有该service，则开启新的goroutine监控service
				r.insertServiceWatch(k)
			}

			// read watchers and delete deregister services
			// 通过对总的services与暂存的watchers对比，停掉已经挂掉的service
			for k, plan := range r.watchers {
				if _, ok := d[k]; !ok {
					global.Log.Infof("%s服务已掉线", k)
					plan.Stop()
					delete(r.watchers, k)
					r.ShutdownChan <- k
				}
			}
			global.Log.Infof("consul: %v", d)
			global.Log.Infof("watchers: %v", r.watchers)
		default:
			fmt.Printf("can't decide the watch type: %v\n", &d)
		}
	}
}

// getSingleSvcHandler uses to get single service handler.
func (r *Registry) getSingleSvcHandler() Handler {
	return func(_ uint64, data any) {
		if d, ok := data.([]*consulapi.ServiceEntry); ok {
			for _, entry := range d {
				serviceAddr := entry.Service.Address
				svc := &Service{
					ServiceID:   entry.Service.ID,
					ServiceName: entry.Service.Service,
					Address:     serviceAddr,
					Port:        entry.Service.Port,
					Status:      entry.Checks.AggregatedStatus(),
				}
				// get os info from KV.
				p, _, err := r.Client.KV().Get(serviceAddr, nil)
				if err != nil {
					global.Log.Errorf("get PV Key [%s] value failed: %v", serviceAddr, err)
				}
				if p != nil {
					svc.ServerOs = string(p.Value)
				}
				global.Log.Infof("服务%s的状态变化：%s", entry.Service.Service, entry.Checks.AggregatedStatus())
				r.StatusChan <- svc
			}
		}
	}
}

// insertServiceWatch start add new single service watch.
func (r *Registry) insertServiceWatch(serviceName string) {
	serviceOpts := map[string]any{
		"service": serviceName,
	}
	servicePlan, err := newWatchPlan("service", serviceOpts, r.getSingleSvcHandler())
	if err != nil {
		fmt.Printf("new service watch failed: %v", err)
	}

	go func() {
		err = runWatchPlan(servicePlan, r.Addr)
		if err != nil {
			fmt.Printf("run single servcie:%s watch failed, error: %v\n", serviceName, err)
		}
	}()
	r.watchers[serviceName] = servicePlan
}

type Handler func(uint64, any)

// newWatchPlan  uses to generate watch plan.
func newWatchPlan(watchType string, opts map[string]any, handler Handler) (*watch.Plan, error) {
	var options = map[string]any{
		"type": watchType,
	}
	// combine params
	for k, v := range opts {
		options[k] = v
	}
	pl, err := watch.Parse(options)
	if err != nil {
		return nil, err
	}
	pl.Handler = watch.HandlerFunc(handler)
	return pl, nil
}

func runWatchPlan(plan *watch.Plan, address string) error {
	err := plan.Run(address)
	if err != nil {
		fmt.Printf("run consul addr: %s error: %v", address, err)
		return err
	}
	return nil
}
