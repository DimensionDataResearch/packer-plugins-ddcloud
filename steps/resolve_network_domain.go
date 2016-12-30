package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// ResolveNetworkDomain is the step that resolves the target network domain by name from CloudControl.
type ResolveNetworkDomain struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *ResolveNetworkDomain) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetSettings().(*config.Settings)
	client := state.GetClient()

	networkDomain, err := client.GetNetworkDomainByName(
		settings.NetworkDomainName,
		settings.DatacenterID,
	)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if networkDomain == nil {
		ui.Error(fmt.Sprintf(
			"Unable to find network domain '%s' in datacenter '%s'.",
			settings.NetworkDomainName,
			settings.DatacenterID,
		))

		return multistep.ActionHalt
	}

	state.SetNetworkDomain(networkDomain)

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *ResolveNetworkDomain) Cleanup(state multistep.StateBag) {
}

var _ multistep.Step = &ResolveNetworkDomain{}
