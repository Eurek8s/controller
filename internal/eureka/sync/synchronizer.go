package sync

import (
	"errors"
	"fmt"
	"github.com/eurek8s/controller/internal/eureka/client"
	"github.com/go-logr/logr"
	"github.com/hudl/fargo"
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"time"
)

var (
	totalHeartbeats = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eurek8s_total_heartbeats",
			Help: "Number of heartbeats processed",
		},
		[]string{"environment", "appName", "appInstance"},
	)
	heartbeatFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eurek8s_heartbeat_failures",
			Help: "Number of failed heartbeats",
		},
		[]string{"environment", "appName", "appInstance"},
	)
	totalRegistrations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eurek8s_total_registration",
			Help: "Number of registrations processed",
		},
		[]string{"environment", "appName", "appInstance"},
	)
	registrationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eurek8s_registration_failures",
			Help: "Number of failed registrations",
		},
		[]string{"environment", "appName", "appInstance"},
	)
	totalDeregistrations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eurek8s_total_deregistration",
			Help: "Number of deregistrations processed",
		},
		[]string{"environment", "appName", "appInstance"},
	)
	deregistrationFailures = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "eurek8s_deregistration_failures",
			Help: "Number of failed deregistrations",
		},
		[]string{"environment", "appName", "appInstance"},
	)
)

func init() {
	metrics.Registry.MustRegister(
		totalHeartbeats, heartbeatFailures,
		totalRegistrations, registrationFailures,
		totalDeregistrations, deregistrationFailures,
	)
}

type Synchronizer struct {
	client         *client.EurekaClient
	applications   map[string]*Application
	registerChan   chan *Application
	deregisterChan chan string
	log            logr.Logger
}

func New(client *client.EurekaClient, log logr.Logger) *Synchronizer {
	return &Synchronizer{
		client:         client,
		applications:   make(map[string]*Application),
		registerChan:   make(chan *Application),
		deregisterChan: make(chan string),
		log:            log,
	}
}

func (s *Synchronizer) goProcess() {
	tickChan := time.NewTicker(time.Second * 10).C

	for {
		select {
		case _ = <-tickChan:
			s.heartbeat()
		case application := <-s.registerChan:
			if err := s.RegisterApplicationSync(application); err != nil {
				s.log.Error(err, "Error trying to Register App (Channel)")
			}
		case key := <-s.deregisterChan:
			s.deregister(key)
		}
	}
}

func (s *Synchronizer) Register(app *Application) {
	s.registerChan <- app
}

func (s *Synchronizer) Deregister(resourceName string) {
	s.deregisterChan <- resourceName
}

func (s *Synchronizer) heartbeat() {


	for _, app := range s.applications {
		for _, i := range app.Instances {
			uniqueId := i.UniqueID(*i)

			totalHeartbeats.
				WithLabelValues(app.Environment, app.Name, uniqueId).
				Inc()

			log := s.log.WithValues("environment", app.Environment, "app", app.Name, "uniqueId", uniqueId)
			log.Info("sending heartbeat request for instance")

			if err := s.client.HeartBeatInstance(app.Environment, i); err != nil {
				log.Error(err, "unable to heartbeat instance")

				heartbeatFailures.
					WithLabelValues(app.Environment, app.Name, uniqueId).
					Inc()
			}
		}
	}
}

func (s *Synchronizer) RegisterApplicationSync(n *Application) error {
	resourceName := n.ResourceName

	if app, contains := s.applications[resourceName]; !contains {
		for _, i := range n.Instances {
			uniqueId := i.UniqueID(*i)

			totalRegistrations.
				WithLabelValues(n.Environment, n.Name, uniqueId).
				Inc()

			log := s.log.WithValues("environment", n.Environment, "app", n.Name, "uniqueId", uniqueId)
			log.Info("trying to register instance")

			if err := s.client.RegisterInstance(n.Environment, i); err != nil {
				log.Error(err, "unable to register instance")

				registrationFailures.
					WithLabelValues(n.Environment, n.Name, uniqueId).
					Inc()

				return errors.New(fmt.Sprintf("error trying to register new application. Resource: %s", resourceName))
			} else {
				s.applications[resourceName] = n
				return nil
			}
		}
		return errors.New(fmt.Sprintf("error. Invalid application. No instances set to be registered. Resource: %s", resourceName))
	} else {
		instances := getInstancesToDeregister(app.Instances, n.Instances)

		for _, i := range instances {
			s.deregisterInstance(app, i)
		}

		delete(s.applications, resourceName)
		return s.RegisterApplicationSync(n)
	}
}

func (s *Synchronizer) deregister(key string) {
	if app, ok := s.applications[key]; !ok {
		s.log.Error(errors.New("unable to deregister app"), "app not found", "key", key)
	} else {
		for _, i := range app.Instances {
			s.deregisterInstance(app, i)
		}

		delete(s.applications, key)
	}
}

func (s *Synchronizer) deregisterInstance(app *Application, i *fargo.Instance) {
	uniqueId := i.UniqueID(*i)

	totalDeregistrations.
		WithLabelValues(app.Environment, app.Name, uniqueId).
		Inc()

	log := s.log.WithValues("environment", app.Environment, "app", app.Name, "uniqueId", uniqueId)
	log.Info("trying to deregister instance")

	if err := s.client.DeregisterInstance(app.Environment, i); err != nil {
		log.Error(err, "unable to deregister instance")

		deregistrationFailures.
			WithLabelValues(app.Environment, app.Name, uniqueId).
			Inc()
	}
}

func (s *Synchronizer) Start() {
	go s.goProcess()
}

func (s *Synchronizer) Stop() {
	close(s.registerChan)
	close(s.deregisterChan)
}

func getInstancesToDeregister(old, new []*fargo.Instance) []*fargo.Instance {
	var result []*fargo.Instance

	for _, o := range old {
		var found bool
		for _, n := range new {
			if o.InstanceId == n.InstanceId {
				found = true
				break
			}
		}

		if !found {
			result = append(result, o)
		}
	}

	return result
}
