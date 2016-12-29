package config

import (
	"fmt"

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
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'mcp_region' has not been specified in settings"),
		)
	}
	if settings.McpUser == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'mcp_user' has not been specified in settings"),
		)
	}
	if settings.McpPassword == "" {
		err = packer.MultiErrorAppend(err,
			fmt.Errorf("'mcp_password' has not been specified in settings"),
		)
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

	return
}
