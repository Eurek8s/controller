package sync

import (
	"github.com/hudl/fargo"
)

type Application struct {
	ResourceName string
	Environment  string
	Name         string
	Instances    []*fargo.Instance
}
