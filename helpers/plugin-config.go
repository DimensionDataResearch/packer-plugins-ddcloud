package helpers

import (
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
)

// PluginConfig represents the basic configuration for a plugin.
type PluginConfig interface {
	// GetPackerConfig retrieves the common Packer configuration for the plugin.
	GetPackerConfig() *common.PackerConfig

	// GetCommunicatorConfig retrieves the Packer communicator configuration for the plugin.
	GetCommunicatorConfig() *communicator.Config

	// Validate ensures that the configuration is valid.
	Validate() error
}
