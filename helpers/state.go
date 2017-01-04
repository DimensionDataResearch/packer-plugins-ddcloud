package helpers

import (
	"log"

	"runtime/debug"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/DimensionDataResearch/packer-plugins-ddcloud/artifacts"
	"github.com/mitchellh/multistep"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/packer"
)

// ForStateBag creates a new `State` helper for the specified multistep.StateBag.
func ForStateBag(stateBag multistep.StateBag) State {
	return State{
		Data: stateBag,
	}
}

// State is the helper for working with `multistep` state data.
type State struct {
	// The state data.
	Data multistep.StateBag
}

// Get retrieves the state data with the specified key.
func (state State) Get(key string) (value interface{}) {
	return state.Data.Get(key)
}

// GetOk retrieves the state data with the specified key, if it exists.
func (state State) GetOk(key string) (value interface{}, exists bool) {
	return state.Data.GetOk(key)
}

// Set updates the state data with the specified key and value.
func (state State) Set(key string, value interface{}) {
	state.Data.Put(key, value)
}

// GetLastError retrieves the last error (if any) from the state data.
func (state State) GetLastError() error {
	value, ok := state.Data.GetOk("error")
	if !ok || value == nil {
		return nil
	}

	return value.(error)
}

// SetLastError updates the last error (if any) in the state data.
func (state State) SetLastError(err error) {
	state.Data.Put("error", err)
}

// GetBuilderID gets the Id of the current builder plugin (if any) in the state data.
func (state State) GetBuilderID() string {
	value, ok := state.Data.GetOk("builder_id")
	if !ok || value == nil {
		return ""
	}

	return value.(string)
}

// SetBuilderID updates the Id of the current builder plugin (if any) in the state data.
func (state State) SetBuilderID(builderID string) {
	state.Data.Put("builder_id", builderID)
}

// GetUI gets a reference to the Packer UI from the state data.
func (state State) GetUI() packer.Ui {
	value, ok := state.Data.GetOk("ui")
	if !ok || value == nil {
		log.Printf("helpers.State.GetUI: Warning - UI not available.\n%s",
			debug.Stack(),
		)

		return nil
	}

	return value.(packer.Ui)
}

// SetUI updates the reference to the Packer UI in the state data.
func (state State) SetUI(ui packer.Ui) {
	state.Data.Put("ui", ui)
}

// GetHook gets a reference to the Packer extensibility hook from the state data.
func (state State) GetHook() packer.Hook {
	value, ok := state.Data.GetOk("hook")
	if !ok || value == nil {
		log.Printf("helpers.State.GetHook: Warning - Hook not available.\n%s",
			debug.Stack(),
		)

		return nil
	}

	return value.(packer.Hook)
}

// SetHook updates the reference to the Packer extensibility hook in the state data.
func (state State) SetHook(hook packer.Hook) {
	state.Data.Put("hook", hook)
}

// GetPackerConfig gets the Packer configuration from the state data.
func (state State) GetPackerConfig() *common.PackerConfig {
	value, ok := state.Data.GetOk("config")
	if !ok || value == nil {
		return nil
	}

	return value.(*common.PackerConfig)
}

// SetPackerConfig updates the Packer configuration in the state data.
func (state State) SetPackerConfig(config *common.PackerConfig) {
	state.Data.Put("config", config)
}

// GetSettings gets the plugin settings from the state data.
func (state State) GetSettings() PluginConfig {
	value, ok := state.Data.GetOk("settings")
	if !ok || value == nil {
		return nil
	}

	return value.(PluginConfig)
}

// SetSettings updates the plugin settings in the state data.
func (state State) SetSettings(config PluginConfig) {
	state.Data.Put("settings", config)
}

// GetClient gets the CloudControl API client from the state data.
func (state State) GetClient() *compute.Client {
	value, ok := state.Data.GetOk("client")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.Client)
}

// SetClient updates the CloudControl API client in the state data.
func (state State) SetClient(client *compute.Client) {
	state.Data.Put("client", client)
}

// GetTargetDatacenter gets the target datacenter from the state data.
func (state State) GetTargetDatacenter() *compute.Datacenter {
	value, ok := state.Data.GetOk("target_datacenter")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.Datacenter)
}

// SetTargetDatacenter updates the target datacenter in the state data.
func (state State) SetTargetDatacenter(datacenter *compute.Datacenter) {
	state.Data.Put("target_datacenter", datacenter)
}

// GetNetworkDomain gets the target network domain from the state data.
func (state State) GetNetworkDomain() *compute.NetworkDomain {
	value, ok := state.Data.GetOk("network_domain")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.NetworkDomain)
}

// SetNetworkDomain updates the target network domain in the state data.
func (state State) SetNetworkDomain(networkDomain *compute.NetworkDomain) {
	state.Data.Put("network_domain", networkDomain)
}

// GetVLAN gets the target VLAN from the state data.
func (state State) GetVLAN() *compute.VLAN {
	value, ok := state.Data.GetOk("vlan")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.VLAN)
}

// SetVLAN updates the target VLAN in the state data.
func (state State) SetVLAN(vlan *compute.VLAN) {
	state.Data.Put("vlan", vlan)
}

// GetServer gets the target server from the state data.
func (state State) GetServer() *compute.Server {
	value, ok := state.Data.GetOk("server")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.Server)
}

// SetServer updates the target server in the state data.
func (state State) SetServer(server *compute.Server) {
	state.Data.Put("server", server)
}

// GetNATRule gets the NAT rule from the state data.
func (state State) GetNATRule() *compute.NATRule {
	value, ok := state.Data.GetOk("nat_rule")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.NATRule)
}

// SetNATRule updates the NAT rule in the state data.
func (state State) SetNATRule(natRule *compute.NATRule) {
	state.Data.Put("nat_rule", natRule)
}

// GetFirewallRule gets the firewall rule from the state data.
func (state State) GetFirewallRule() *compute.FirewallRule {
	value, ok := state.Data.GetOk("firewall_rule")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.FirewallRule)
}

// SetFirewallRule updates the firewall rule in the state data.
func (state State) SetFirewallRule(firewallRule *compute.FirewallRule) {
	state.Data.Put("firewall_rule", firewallRule)
}

// GetSourceImage gets the source image from the state data.
func (state State) GetSourceImage() compute.Image {
	value, ok := state.Data.GetOk("source_image")
	if !ok || value == nil {
		return nil
	}

	return value.(compute.Image)
}

// SetSourceImage updates the source image in the state data.
func (state State) SetSourceImage(image compute.Image) {
	state.Data.Put("source_image", image)
}

// GetSourceImageArtifact gets the source image artifact from the state data.
func (state State) GetSourceImageArtifact() *artifacts.Image {
	value, ok := state.Data.GetOk("source_image_artifact")
	if !ok || value == nil {
		return nil
	}

	return value.(*artifacts.Image)
}

// SetSourceImageArtifact updates the source image artifact in the state data.
func (state State) SetSourceImageArtifact(sourceArtifact *artifacts.Image) {
	state.Data.Put("source_image_artifact", sourceArtifact)
}

// GetTargetImage gets the target image from the state data.
func (state State) GetTargetImage() *compute.CustomerImage {
	value, ok := state.Data.GetOk("target_image")
	if !ok || value == nil {
		return nil
	}

	return value.(*compute.CustomerImage)
}

// SetTargetImage updates the target image in the state data.
func (state State) SetTargetImage(image *compute.CustomerImage) {
	state.Data.Put("target_image", image)
}

// GetTargetImageArtifact gets the target image artifact from the state data.
func (state State) GetTargetImageArtifact() *artifacts.Image {
	value, ok := state.Data.GetOk("target_image_artifact")
	if !ok || value == nil {
		return nil
	}

	return value.(*artifacts.Image)
}

// SetTargetImageArtifact updates the target image artifact in the state data.
func (state State) SetTargetImageArtifact(sourceArtifact *artifacts.Image) {
	state.Data.Put("target_image_artifact", sourceArtifact)
}

// GetRemoteOVFPackageArtifact gets the remote OVF package artifact from the state data.
func (state State) GetRemoteOVFPackageArtifact() *artifacts.RemoteOVFPackage {
	value, ok := state.Data.GetOk("remote_ovf_package_artifact")
	if !ok || value == nil {
		return nil
	}

	return value.(*artifacts.RemoteOVFPackage)
}

// SetRemoteOVFPackageArtifact updates the remote OVF package artifact in the state data.
func (state State) SetRemoteOVFPackageArtifact(packageArtifact *artifacts.RemoteOVFPackage) {
	state.Data.Put("remote_ovf_package_artifact", packageArtifact)
}

// GetSourceArtifact gets the source artifact from the state data.
func (state State) GetSourceArtifact() packer.Artifact {
	value, ok := state.Data.GetOk("source_artifact")
	if !ok || value == nil {
		return nil
	}

	return value.(*artifacts.Image)
}

// SetSourceArtifact updates the source artifact in the state data.
func (state State) SetSourceArtifact(sourceArtifact packer.Artifact) {
	state.Data.Put("source_artifact", sourceArtifact)
}

// GetTargetArtifact gets the target artifact from the state data.
func (state State) GetTargetArtifact() packer.Artifact {
	value, ok := state.Data.GetOk("target_artifact")
	if !ok || value == nil {
		return nil
	}

	return value.(*artifacts.Image)
}

// SetTargetArtifact updates the target artifact in the state data.
func (state State) SetTargetArtifact(targetArtifact packer.Artifact) {
	state.Data.Put("target_artifact", targetArtifact)
}
