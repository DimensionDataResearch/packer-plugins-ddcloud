package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// CreateNATRule is the step that exposes the target server using a NAT rule.
type CreateNATRule struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *CreateNATRule) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)

	ui := state.GetUI()

	settings := state.GetConfig().(*config.Settings)
	client := state.GetClient()
	networkDomain := state.GetNetworkDomain()
	server := state.GetServer()

	if settings.UsePrivateIPv4 {
		ui.Message(fmt.Sprintf(
			"Server '%s' will not be exposed because the configuration specifies 'use_private_ipv4'.",
			server.Name,
		))

		return multistep.ActionContinue
	}

	if settings.CommunicatorConfig.Type == "" {
		ui.Message(fmt.Sprintf(
			"Server '%s' will not be exposed because no communicator is configured.",
			server.Name,
		))

		return multistep.ActionContinue
	}

	privateIPv4Address := *server.Network.PrimaryAdapter.PrivateIPv4Address
	ui.Message(fmt.Sprintf(
		"Creating NAT rule for server '%s' ('%s') with private IPv4 address '%s'...",
		server.Name,
		server.ID,
		privateIPv4Address,
	))

	natRuleID, err := client.AddNATRule(
		networkDomain.ID,
		privateIPv4Address,
		nil, // Auto-select public IPv4 address
	)
	if err != nil {
		if compute.IsNoIPAddressAvailableError(err) {
			ui.Message(fmt.Sprintf(
				"Network domain '%s' ('%s') has no public IP addresses available; a new block will now be allocated...",
				networkDomain.Name,
				networkDomain.ID,
			))

			publicIPBlockID, err := client.AddPublicIPBlock(networkDomain.ID)
			if err != nil {
				ui.Error(err.Error())

				return multistep.ActionHalt
			}

			ui.Message(fmt.Sprintf(
				"Allocated new public IP block '%s' in network domain '%s' ('%s').",
				publicIPBlockID,
				networkDomain.Name,
				networkDomain.ID,
			))

			natRuleID, err = client.AddNATRule(
				networkDomain.ID,
				privateIPv4Address,
				nil, // Auto-select public IPv4 address
			)

			if err != nil {
				ui.Error(err.Error())

				return multistep.ActionHalt
			}
		} else {
			if err != nil {
				ui.Error(err.Error())

				return multistep.ActionHalt
			}
		}
	}

	natRule, err := client.GetNATRule(natRuleID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}
	if natRule == nil {
		ui.Error(fmt.Sprintf(
			"Cannot find newly-created NAT rule '%s'.",
			natRuleID,
		))

		return multistep.ActionHalt
	}

	ui.Message(fmt.Sprintf(
		"Created NAT rule '%s' for server '%s' ('%s') from private IPv4 address '%s' to public IPv4 address '%s'.",
		natRuleID,
		server.Name,
		server.ID,
		natRule.InternalIPAddress,
		natRule.ExternalIPAddress,
	))

	state.SetNATRule(natRule)

	// Override SSH / WinRM connection details, if required.
	communicatorConfig := &settings.CommunicatorConfig
	if communicatorConfig.SSHHost == natRule.InternalIPAddress {
		communicatorConfig.SSHHost = natRule.ExternalIPAddress
	}
	if communicatorConfig.WinRMHost == natRule.InternalIPAddress {
		communicatorConfig.WinRMHost = natRule.ExternalIPAddress
	}

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *CreateNATRule) Cleanup(stateBag multistep.StateBag) {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	client := state.GetClient()
	server := state.GetServer()

	natRule := state.GetNATRule()
	if natRule == nil {
		return // Nothing to do.
	}

	ui.Message(fmt.Sprintf(
		"Destroying NAT rule '%s' ('%s' -> '%s') for server '%s' ('%s')...",
		natRule.ID,
		natRule.ExternalIPAddress,
		natRule.InternalIPAddress,
		server.Name,
		server.ID,
	))

	err := client.DeleteNATRule(natRule.ID)
	if err != nil {
		ui.Error(err.Error())

		return
	}

	state.SetNATRule(nil)

	ui.Message(fmt.Sprintf(
		"Destroyed NAT rule '%s' ('%s' -> '%s') for server '%s' ('%s').",
		natRule.ID,
		natRule.ExternalIPAddress,
		natRule.InternalIPAddress,
		server.Name,
		server.ID,
	))
}

var _ multistep.Step = &CreateNATRule{}
