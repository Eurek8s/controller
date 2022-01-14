package client

import (
	"fmt"
	"github.com/hudl/fargo"
	"github.com/pkg/errors"
)

type EurekaClient struct {
	connections map[string]fargo.EurekaConnection
}

func New(addresses map[string][]string) *EurekaClient {
	connections := make(map[string]fargo.EurekaConnection)
	for k, v := range addresses {
		connections[k] = fargo.NewConn(v...)
	}

	return &EurekaClient{connections: connections}
}

func (c *EurekaClient) call(
	environment string,
	i *fargo.Instance,
	f func(c fargo.EurekaConnection, i *fargo.Instance) error,
) error {
	if conn, ok := c.connections[environment]; !ok {
		return errors.New(fmt.Sprintf("cannot find eureka connection for environment \"%s\"", environment))
	} else if err := f(conn, i); err != nil {
		statusCode, _ := fargo.HTTPResponseStatusCode(err)

		return errors.Wrap(err, fmt.Sprintf("invalid status code received: %d", statusCode))
	}

	return nil
}

func (c *EurekaClient) HeartBeatInstance(environment string, i *fargo.Instance) error {
	return c.call(
		environment,
		i,
		func(c fargo.EurekaConnection, i *fargo.Instance) error { return c.HeartBeatInstance(i) },
	)
}

func (c *EurekaClient) RegisterInstance(environment string, i *fargo.Instance) error {
	return c.call(
		environment,
		i,
		func(c fargo.EurekaConnection, i *fargo.Instance) error { return c.RegisterInstance(i) },
	)
}

func (c *EurekaClient) DeregisterInstance(environment string, i *fargo.Instance) error {
	return c.call(
		environment,
		i,
		func(c fargo.EurekaConnection, i *fargo.Instance) error { return c.DeregisterInstance(i) },
	)
}

func (c *EurekaClient) GetApp(environment, appName string) (*fargo.Application, error) {
	if conn, ok := c.connections[environment]; !ok {
		return nil, errors.New(fmt.Sprintf("cannot find eureka connection for environment \"%s\"", environment))
	} else if app, err := conn.GetApp(appName); err != nil {
		return app, err
	} else {
		return nil, errors.New(fmt.Sprintf("unable to get app: %s", appName))
	}
}
