package main

import "fmt"

// Configuration represents the configuration for the Builder.
type Configuration struct {
	McpRegion       string `mapstructure:"mcp_region"`
	McpUser         string `mapstructure:"mcp_user"`
	McpPassword     string `mapstructure:"mcp_password"`
	NetworkDomainID string `mapstructure:"networkdomain"`
	VLANID          string `mapstructure:"vlan"`
	SourceImage     string `mapstructure:"source_image"`
	TargetImage     string `mapstructure:"target_image"`
	uniquenessKey   string
	serverName      string
}

// Validate determines if the configuration is valid.
func (config *Configuration) Validate() (err error) {
	if config.McpRegion == "" {
		err = fmt.Errorf("'mcp_region' has not been specified in configuration")

		return
	}
	if config.McpUser == "" {
		err = fmt.Errorf("'mcp_user' has not been specified in configuration")

		return
	}
	if config.McpPassword == "" {
		err = fmt.Errorf("'mcp_password' has not been specified in configuration")

		return
	}
	if config.NetworkDomainID == "" {
		err = fmt.Errorf("'networkdomain' has not been specified in configuration")

		return
	}
	if config.VLANID == "" {
		err = fmt.Errorf("'vlan' has not been specified in configuration")

		return
	}
	if config.SourceImage == "" {
		err = fmt.Errorf("'source_image' has not been specified in configuration")

		return
	}
	if config.TargetImage == "" {
		err = fmt.Errorf("'target_image' has not been specified in configuration")

		return
	}

	return
}
