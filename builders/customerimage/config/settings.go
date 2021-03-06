package config

import (
	"fmt"
	"os"

	"time"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud/helpers"
	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
)

// Settings represents the settings for the customer image builder.
type Settings struct {
	PackerConfig       common.PackerConfig `mapstructure:",squash"`
	CommunicatorConfig communicator.Config `mapstructure:",squash"`

	McpRegion            string `mapstructure:"mcp_region"`
	McpUser              string `mapstructure:"mcp_user"`
	McpPassword          string `mapstructure:"mcp_password"`
	DatacenterID         string `mapstructure:"datacenter"`
	NetworkDomainName    string `mapstructure:"networkdomain"`
	VLANName             string `mapstructure:"vlan"`
	SourceImage          string `mapstructure:"source_image"`
	TargetImage          string `mapstructure:"target_image"`
	InitialAdminPassword string `mapstructure:"initial_admin_password"`
	UsePrivateIPv4       bool   `mapstructure:"use_private_ipv4"`
	ClientIP             string `mapstructure:"client_ip"`
	UniquenessKey        string
	ServerName           string
}

var _ helpers.PluginConfig = &Settings{}

// GetPackerConfig retrieves the common Packer configuration for the plugin.
func (settings *Settings) GetPackerConfig() *common.PackerConfig {
	return &settings.PackerConfig
}

// GetCommunicatorConfig retrieves the Packer communicator configuration for the plugin.
func (settings *Settings) GetCommunicatorConfig() *communicator.Config {
	return &settings.CommunicatorConfig
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
	if settings.DatacenterID == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'datacenter' has not been specified in settings"),
		)
	}
	if settings.NetworkDomainName == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'networkdomain' has not been specified in settings"),
		)
	}
	if settings.VLANName == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'vlan' has not been specified in settings"),
		)
	}
	if settings.SourceImage == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'source_image' has not been specified in settings"),
		)
	}
	if settings.TargetImage == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'target_image' has not been specified in settings"),
		)
	}

	// Communicator defaults.
	settings.CommunicatorConfig.SSHTimeout = 2 * time.Minute
	settings.CommunicatorConfig.WinRMTimeout = 2 * time.Minute
	if settings.CommunicatorConfig.Type == "" {
		settings.CommunicatorConfig.Type = "none"
	} else if !settings.UsePrivateIPv4 {
		if settings.ClientIP == "" {
			err = packer.MultiErrorAppend(err,
				fmt.Errorf("'use_private_ipv4' has been specified in settings, but 'client_ip' has not"),
			)
		}
	}
	if settings.CommunicatorConfig.SSHHost == "" {
		settings.CommunicatorConfig.SSHHost = settings.ServerName
	}
	if settings.CommunicatorConfig.SSHPort == 0 {
		settings.CommunicatorConfig.SSHPort = 22
	}
	if settings.CommunicatorConfig.SSHUsername == "" {
		settings.CommunicatorConfig.SSHUsername = "root"
	}
	if settings.CommunicatorConfig.SSHPassword == "" {
		settings.CommunicatorConfig.SSHPassword = settings.InitialAdminPassword
	}
	if settings.CommunicatorConfig.WinRMHost == "" {
		settings.CommunicatorConfig.WinRMHost = settings.ServerName
	}
	if settings.CommunicatorConfig.WinRMPort == 0 {
		settings.CommunicatorConfig.WinRMPort = 5895
	}
	if settings.CommunicatorConfig.WinRMUser == "" {
		settings.CommunicatorConfig.WinRMUser = "Administrator"
	}
	if settings.CommunicatorConfig.WinRMPassword == "" {
		settings.CommunicatorConfig.WinRMPassword = settings.InitialAdminPassword
	}

	return
}
