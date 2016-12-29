package config

import "fmt"

// Settings represents the settings for the Builder.
type Settings struct {
	McpRegion       string `mapstructure:"mcp_region"`
	McpUser         string `mapstructure:"mcp_user"`
	McpPassword     string `mapstructure:"mcp_password"`
	NetworkDomainID string `mapstructure:"networkdomain"`
	VLANID          string `mapstructure:"vlan"`
	SourceImage     string `mapstructure:"source_image"`
	TargetImage     string `mapstructure:"target_image"`
	UniquenessKey   string
	ServerName      string
}

// Validate determines if the settings is valid.
func (settings *Settings) Validate() (err error) {
	if settings.McpRegion == "" {
		err = fmt.Errorf("'mcp_region' has not been specified in settings")

		return
	}
	if settings.McpUser == "" {
		err = fmt.Errorf("'mcp_user' has not been specified in settings")

		return
	}
	if settings.McpPassword == "" {
		err = fmt.Errorf("'mcp_password' has not been specified in settings")

		return
	}
	if settings.NetworkDomainID == "" {
		err = fmt.Errorf("'networkdomain' has not been specified in settings")

		return
	}
	if settings.VLANID == "" {
		err = fmt.Errorf("'vlan' has not been specified in settings")

		return
	}
	if settings.SourceImage == "" {
		err = fmt.Errorf("'source_image' has not been specified in settings")

		return
	}
	if settings.TargetImage == "" {
		err = fmt.Errorf("'target_image' has not been specified in settings")

		return
	}

	return
}
