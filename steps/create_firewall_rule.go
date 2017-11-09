package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/multistep"
)

// CreateFirewallRule is the step that exposes the target server using a firewall rule.
//
// The server's associated NAT rule must already have been created.
type CreateFirewallRule struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *CreateFirewallRule) Run(stateBag multistep.StateBag) multistep.StepAction {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetSettings().(*config.Settings)
	client := state.GetClient()
	networkDomain := state.GetNetworkDomain()
	server := state.GetServer()
	natRule := state.GetNATRule()

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

	ui.Message(fmt.Sprintf(
		"Creating firewall rule to permit access for server '%s' ('%s') via public IPv4 address '%s'...",
		server.Name,
		server.ID,
		natRule.ExternalIPAddress,
	))

	firewallRuleConfiguration := &compute.FirewallRuleConfiguration{
		Name:            fmt.Sprintf("packer.%s.inbound", settings.UniquenessKey),
		NetworkDomainID: networkDomain.ID,
	}
	firewallRuleConfiguration.Accept()
	firewallRuleConfiguration.IPv4()
	firewallRuleConfiguration.IP()
	firewallRuleConfiguration.MatchSourceAddress(settings.ClientIP)
	firewallRuleConfiguration.MatchDestinationAddress(natRule.ExternalIPAddress)
	firewallRuleConfiguration.PlaceFirst()
	firewallRuleConfiguration.Enable()

	firewallRuleID, err := client.CreateFirewallRule(*firewallRuleConfiguration)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	firewallRule, err := client.GetFirewallRule(firewallRuleID)
	if err != nil {
		ui.Error(err.Error())

		return multistep.ActionHalt
	}

	if firewallRule == nil {
		ui.Error(fmt.Sprintf(
			"Cannot find newly-created firewall rule '%s'.",
			firewallRuleID,
		))

		return multistep.ActionHalt
	}

	state.SetFirewallRule(firewallRule)

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *CreateFirewallRule) Cleanup(stateBag multistep.StateBag) {
	state := helpers.ForStateBag(stateBag)
	ui := state.GetUI()

	settings := state.GetSettings().(*config.Settings)
	client := state.GetClient()
	server := state.GetServer()

	firewallRule := state.GetFirewallRule()
	if settings.UsePrivateIPv4 || firewallRule == nil {
		return // Nothing to do.
	}

	ui.Message(fmt.Sprintf(
		"Destroying firewall rule '%s' ('%s') for server '%s' ('%s')...",
		firewallRule.Name,
		firewallRule.ID,
		server.Name,
		server.ID,
	))

	err := client.DeleteFirewallRule(firewallRule.ID)
	if err != nil {
		ui.Error(err.Error())

		return
	}

	state.SetFirewallRule(nil)

	ui.Message(fmt.Sprintf(
		"Destroyed firewall rule '%s' ('%s') for server '%s' ('%s').",
		firewallRule.Name,
		firewallRule.ID,
		server.Name,
		server.ID,
	))
}

var _ multistep.Step = &CreateFirewallRule{}
