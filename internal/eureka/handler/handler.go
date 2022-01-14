package handler

import (
	"context"
	"fmt"
	discoveryv1 "github.com/eurek8s/controller/api/v1"
	eurek8ssyncer "github.com/eurek8s/controller/internal/eureka/sync"
	"github.com/eurek8s/controller/internal/eureka/util"
	"github.com/go-logr/logr"
	"github.com/hudl/fargo"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	DefaultEnvironment = "qa"

	NoZone      = "no-zone"
	DefaultZone = "default-zone"

	FinalizerName = "finalizers.eurekaapplication.discovery.eurek8s.com"

	protocolHttp  = "http"
	protocolHttps = "https"

	httpsPort = 443
)

type Handler struct {
	EurekaSyncer *eurek8ssyncer.Syncer
	log          logr.Logger
}

func New(syncer *eurek8ssyncer.Syncer, log logr.Logger) *Handler {
	return &Handler{EurekaSyncer: syncer, log: log}
}

type hostPort struct {
	host string
	port int32
}

func (h *Handler) Handle(
	ctx context.Context,
	c client.Client,
	spec *discoveryv1.EurekaApplication,
	resourceName string,
) error {
	environment := spec.Spec.Environment
	if environment == "" {
		environment = DefaultEnvironment
	}

	if spec.ObjectMeta.DeletionTimestamp.IsZero() {
		if !util.ContainsString(spec.ObjectMeta.Finalizers, FinalizerName) {
			h.log.Info("Registering finalizer", "name", FinalizerName)
			spec.ObjectMeta.Finalizers = append(spec.ObjectMeta.Finalizers, FinalizerName)
			if err := c.Update(ctx, spec); err != nil {
				return err
			}
		}
	} else {
		if util.ContainsString(spec.ObjectMeta.Finalizers, FinalizerName) {
			h.log.Info("Deregistering finalizer", "name", FinalizerName)
			h.EurekaSyncer.Deregister(resourceName)

			spec.ObjectMeta.Finalizers = util.RemoveString(spec.ObjectMeta.Finalizers, FinalizerName)
			if err := c.Update(ctx, spec); err != nil {
				return err
			}
		}

		return nil
	}

	disabled := spec.Spec.Disabled

	if disabled {
		h.EurekaSyncer.Deregister(resourceName)
	} else if app, err := getEurekaApplication(ctx, c, spec, environment, resourceName); err != nil {
		return err
	} else {
		h.EurekaSyncer.Register(app)
	}

	return nil
}

func getHostPorts(
	ctx context.Context,
	c client.Client,
	ingress networkingv1.Ingress,
) ([]hostPort, error) {
	var hostPorts []hostPort
	for _, rule := range ingress.Spec.Rules {
		for _, path := range rule.HTTP.Paths {
			var port int32

			if path.Backend.Service.Port.Name != "" {
				var service v1.Service
				namespacedName := types.NamespacedName{Namespace: ingress.Namespace, Name: path.Backend.Service.Name}
				if err := c.Get(ctx, namespacedName, &service); err != nil {
					return nil, err
				}

				for _, sport := range service.Spec.Ports {
					if sport.Name == path.Backend.Service.Name {
						port = sport.Port
						break
					}
				}
			} else {
				port = path.Backend.Service.Port.Number
			}

			hostPorts = append(hostPorts, hostPort{host: rule.Host, port: port})
		}
	}

	return hostPorts, nil
}

func getEurekaApplication(
	ctx context.Context,
	c client.Client,
	spec *discoveryv1.EurekaApplication,
	environment string,
	resourceName string,
) (*eurek8ssyncer.Application, error) {
	zone := spec.Spec.Zone
	if zone == "" {
		zone = DefaultZone
	}

	metadata := map[string]string{"zone": zone}
	if zone == NoZone {
		metadata = map[string]string{}
	}

	var ingress networkingv1.Ingress
	nn := types.NamespacedName{Namespace: spec.Namespace, Name: spec.Spec.IngressName}
	if err := c.Get(ctx, nn, &ingress); err != nil {
		return nil, err
	}

	app := &eurek8ssyncer.Application{
		ResourceName: resourceName,
		Environment:  environment,
		Name:         spec.Spec.AppName,
	}

	hostPorts, err := getHostPorts(ctx, c, ingress)
	if err != nil {
		return nil, err
	}

	for _, hostPort := range hostPorts {
		rawHost, rawPort := hostPort.host, hostPort.port
		protocol := protocolHttp
		if rawPort == httpsPort {
			protocol = protocolHttps
		}
		host := fmt.Sprintf("%s://%s:%d", protocol, rawHost, rawPort)

		statusUrl, err := util.JoinPathStr(host, spec.Spec.Paths.Status)
		if err != nil {
			return nil, errors.Wrap(err, "invalid host or path set for application status address")
		}

		healthCheckUrl, err := util.JoinPathStr(host, spec.Spec.Paths.HealthCheck)
		if err != nil {
			return nil, errors.Wrap(err, "invalid host or path set for application healthcheck address")
		}

		homeUrl, err := util.JoinPathStr(host, spec.Spec.Paths.Home)
		if err != nil {
			return nil, errors.Wrap(err, "invalid host or path set for application home address")
		}

		i := &fargo.Instance{
			UniqueID: func(i fargo.Instance) string {
				return strings.ToLower(fmt.Sprintf("%s:%s:%d", i.App, i.HostName, i.Port))
			},
			InstanceId:       strings.ToLower(fmt.Sprintf("%s:%s:%d", app.Name, rawHost, rawPort)),
			HostName:         rawHost,
			IPAddr:           rawHost,
			App:              app.Name,
			VipAddress:       app.Name,
			SecureVipAddress: app.Name,
			HomePageUrl:      homeUrl,
			StatusPageUrl:    statusUrl,
			HealthCheckUrl:   healthCheckUrl,
			Status:           fargo.UP,
			Port:             int(rawPort),
			PortEnabled:      true,
			DataCenterInfo:   fargo.DataCenterInfo{Name: fargo.MyOwn},
			Metadata:         fargo.InstanceMetadata{},
		}

		for key, value := range metadata {
			i.SetMetadataString(key, value)
		}

		if rawPort == httpsPort {
			i.SecurePort = i.Port
			i.SecurePortEnabled = true
			i.PortEnabled = false
		}

		app.Instances = append(app.Instances, i)
	}

	return app, nil
}
