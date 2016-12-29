package steps

import (
	"fmt"
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// DestroyServer is the step that destroys the target server in CloudControl.
type DestroyServer struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step DestroyServer) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	client := state.Get("client").(*compute.Client)
	server := state.Get("server").(*compute.Server)

	ui.Message(fmt.Sprintf(
		"Destroying server '%s' ('%s')...",
		server.Name,
		server.ID,
	))

	err := client.DeleteServer(server.ID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	err = client.WaitForDelete(compute.ResourceTypeServer, server.ID, 20*time.Minute)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf(
		"Destroyed server '%s' ('%s').",
		server.Name,
		server.ID,
	))

	state.Put("server", nil)

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step DestroyServer) Cleanup(state multistep.StateBag) {
	// TODO: Destroy server.
}

var _ multistep.Step = &DestroyServer{}
