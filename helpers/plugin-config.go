package helpers

import (
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
)

// PluginConfig represents the basic configuration for a plugin.
type PluginConfig interface {
	// GetPackerConfig retrieves the common Packer configuration for the plugin.
	GetPackerConfig() *common.PackerConfig

	// GetCommunicatorConfig retrieves the Packer communicator configuration (if available) for the plugin.
	GetCommunicatorConfig() *communicator.Config

	// GetMCPUser retrieves the Cloud Control user name.
	GetMCPUser() string

	// GetMCPPassword retrieves the Cloud Control password.
	GetMCPPassword() string

	// Validate ensures that the configuration is valid.
	Validate() error
}
