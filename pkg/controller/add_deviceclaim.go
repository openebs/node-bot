package controller

import (
	"github.com/openebs/node-disk-manager/pkg/controller/deviceclaim"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, deviceclaim.Add)
}
