package steps

/*
 * Code imported from Packer due to vendoring weirdness.
 */

import (
	"log"
	"time"

	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// Provision runs the provisioners.
//
// Uses:
//   communicator packer.Communicator
//   hook         packer.Hook
//   ui           packer.Ui
//
// Produces:
//   <nothing>
type Provision struct {
	Comm packer.Communicator
}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step Provision) Run(state multistep.StateBag) multistep.StepAction {
	comm := step.Comm
	if comm == nil {
		raw, ok := state.Get("communicator").(packer.Communicator)
		if ok {
			comm = raw.(packer.Communicator)
		}
	}
	hook := state.Get("hook").(packer.Hook)
	ui := state.Get("ui").(packer.Ui)

	// Run the provisioner in a goroutine so we can continually check
	// for cancellations...
	log.Println("Running the provision hook")
	errCh := make(chan error, 1)
	go func() {
		errCh <- hook.Run(packer.HookProvision, ui, comm, nil)
	}()

	for {
		select {
		case err := <-errCh:
			if err != nil {
				state.Put("error", err)
				return multistep.ActionHalt
			}

			return multistep.ActionContinue
		case <-time.After(1 * time.Second):
			if _, ok := state.GetOk(multistep.StateCancelled); ok {
				log.Println("Cancelling provisioning due to interrupt...")
				hook.Cancel()
				return multistep.ActionHalt
			}
		}
	}
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step Provision) Cleanup(multistep.StateBag) {}

var _ multistep.Step = &Provision{}
