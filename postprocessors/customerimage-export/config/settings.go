package config

import (
	"fmt"
	"os"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
)

// Settings represents the settings for the customer image export post-processor.
type Settings struct {
	PackerConfig common.PackerConfig `mapstructure:",squash"`

	McpRegion                string `mapstructure:"mcp_region"`
	McpUser                  string `mapstructure:"mcp_user"`
	McpPassword              string `mapstructure:"mcp_password"`
	DatacenterID             string `mapstructure:"datacenter"`
	TargetImageName          string `mapstructure:"target_image"`
	OVFPackagePrefix         string `mapstructure:"ovf_package_prefix"`
	DownloadToLocalDirectory string `mapstructure:"download_to_local_directory"`
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

// GetMCPUser retrieves the Cloud Control user name.
func (settings *Settings) GetMCPUser() string {
	return settings.McpUser
}

// GetMCPPassword retrieves the Cloud Control password.
func (settings *Settings) GetMCPPassword() string {
	return settings.McpPassword
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
	if settings.TargetImageName == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'target_image' has not been specified in settings"),
		)
	}
	if settings.DatacenterID == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'datacenter' has not been specified in settings"),
		)
	}
	if settings.OVFPackagePrefix == "" {
		settings.OVFPackagePrefix = settings.TargetImageName
	}
	if settings.DownloadToLocalDirectory != "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'download_to_local_directory' is not supported yet"),
		)
	}

	return
}
