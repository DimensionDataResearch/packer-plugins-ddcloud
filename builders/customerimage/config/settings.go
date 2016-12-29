package config

import (
	"fmt"
	"os"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/communicator"
	"github.com/mitchellh/packer/packer"
)

// Settings represents the settings for the Builder.
type Settings struct {
	PackerConfig       common.PackerConfig `mapstructure:",squash"`
	CommunicatorConfig communicator.Config `mapstructure:",squash"`

	McpRegion         string `mapstructure:"mcp_region"`
	McpUser           string `mapstructure:"mcp_user"`
	McpPassword       string `mapstructure:"mcp_password"`
	DatacenterID      string `mapstructure:"datacenter"`
	NetworkDomainName string `mapstructure:"networkdomain"`
	VLANName          string `mapstructure:"vlan"`
	SourceImage       string `mapstructure:"source_image"`
	TargetImage       string `mapstructure:"target_image"`
	UniquenessKey     string
	ServerName        string
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
	if settings.CommunicatorConfig.Type == "" {
		settings.CommunicatorConfig.Type = "none"
	}
	if settings.CommunicatorConfig.SSHHost == "" {
		settings.CommunicatorConfig.SSHHost = settings.ServerName
	}
	if settings.CommunicatorConfig.SSHUsername == "" {
		settings.CommunicatorConfig.SSHUsername = "root"
	}
	if settings.CommunicatorConfig.SSHPassword == "" {
		settings.CommunicatorConfig.SSHPassword = settings.UniquenessKey
	}
	if settings.CommunicatorConfig.WinRMHost == "" {
		settings.CommunicatorConfig.WinRMHost = settings.ServerName
	}
	if settings.CommunicatorConfig.WinRMUser == "" {
		settings.CommunicatorConfig.WinRMUser = "Administrator"
	}
	if settings.CommunicatorConfig.WinRMPassword == "" {
		settings.CommunicatorConfig.WinRMHost = settings.UniquenessKey
	}

	return
}
