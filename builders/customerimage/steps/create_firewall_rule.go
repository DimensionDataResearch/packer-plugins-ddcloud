package steps

import (
	"fmt"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/builders/customerimage/config"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/packer"
)

// CreateFirewallRule is the step that exposes the target server using a firewall rule.
//
// The server's associated NAT rule must already have been created.
type CreateFirewallRule struct{}

// Run is called to perform the step's action.
//
// The return value determines whether multi-step sequences should continue or halt.
func (step *CreateFirewallRule) Run(state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packer.Ui)

	settings := state.Get("settings").(*config.Settings)
	client := state.Get("client").(*compute.Client)
	networkDomain := state.Get("network_domain").(*compute.NetworkDomain)
	server := state.Get("server").(*compute.Server)

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

	natRule := state.Get("nat_rule").(*compute.NATRule)

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

	state.Put("firewall_rule", firewallRule)

	return multistep.ActionContinue
}

// Cleanup is called in reverse order of the steps that have run
// and allow steps to clean up after themselves. Do not assume if this
// ran that the entire multi-step sequence completed successfully. This
// method can be ran in the face of errors and cancellations as well.
//
// The parameter is the same "state bag" as Run, and represents the
// state at the latest possible time prior to calling Cleanup.
func (step *CreateFirewallRule) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(packer.Ui)

	client := state.Get("client").(*compute.Client)
	server := state.Get("server").(*compute.Server)

	value, _ := state.GetOk("firewall_rule")
	if value == nil {
		return
	}
	firewallRule := value.(*compute.FirewallRule)

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

	state.Put("firewall_rule", nil)

	ui.Message(fmt.Sprintf(
		"Destroyed firewall rule '%s' ('%s') for server '%s' ('%s').",
		firewallRule.Name,
		firewallRule.ID,
		server.Name,
		server.ID,
	))
}

var _ multistep.Step = &CreateFirewallRule{}
