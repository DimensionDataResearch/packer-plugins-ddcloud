package config

import (
	"fmt"
	"os"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
)

// Settings represents the settings for the Builder.
type Settings struct {
	PackerConfig common.PackerConfig `mapstructure:",squash"`

	McpRegion        string `mapstructure:"mcp_region"`
	McpUser          string `mapstructure:"mcp_user"`
	McpPassword      string `mapstructure:"mcp_password"`
	OVFPackagePrefix string `mapstructure:"ovf_package_prefix"`
}

var _ helpers.PluginConfig = &Settings{}

// GetPackerConfig retrieves the common Packer configuration for the plugin.
func (settings *Settings) GetPackerConfig() *common.PackerConfig {
	return &settings.PackerConfig
}

// GetCommunicatorConfig retrieves the Packer communicator configuration (if available) for the plugin.
func (settings *Settings) GetCommunicatorConfig() *communicator.Config {
	return nil
}

// Validate determines if the settings is valid.
func (settings *Settings) Validate() (err error) {
	if settings.McpRegion == "" {
		settings.McpRegion = os.Getenv("MCP_REGION")

		if settings.McpRegion == "" {
			err = packer.MultiErrorAppend(err,
				fmt.Errorf("'mcp_region' has not been specified in settings and the MCP_REGION environment variable has not been set"),
			)
		}
	}
	if settings.McpUser == "" {
		settings.McpUser = os.Getenv("MCP_USER")

		if settings.McpUser == "" {
			err = packer.MultiErrorAppend(err,
				fmt.Errorf("'mcp_user' has not been specified in settings and the MCP_USER environment variable has not been set"),
			)
		}
	}
	if settings.McpPassword == "" {
		settings.McpPassword = os.Getenv("MCP_PASSWORD")

		if settings.McpPassword == "" {
			err = packer.MultiErrorAppend(err,
				fmt.Errorf("'mcp_password' has not been specified in settings and the MCP_PASSWORD environment variable has not been set"),
			)
		}
	}
	if settings.OVFPackagePrefix == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'ovf_package_prefix' has not been specified in settings"),
		)
	}

	return
}
