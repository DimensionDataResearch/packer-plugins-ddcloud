package ddcloud

import (
	"fmt"
	"os"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
	"github.com/mitchellh/mapstructure"
	"github.com/mitchellh/packer/packer"
)

// BuilderID is the unique Id for the ddcloud builder
const BuilderID = "dimension-data-research.ddcloud"

// Builder is the Builder plugin for Packer.
type Builder struct {
	config *Configuration
	client *compute.Client
}

// Configuration represents the configuration for the Builder.
type Configuration struct {
	McpRegion       string `mapstructure:"mcp_region"`
	McpUser         string `mapstructure:"mcp_user"`
	McpPassword     string `mapstructure:"mcp_password"`
	NetworkDomainID string `mapstructure:"networkdomain"`
	SourceImage     string `mapstructure:"source_image"`
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
	if config.SourceImage == "" {
		err = fmt.Errorf("'source_image' has not been specified in configuration")

		return
	}

	return
}

// Prepare the plugin to run.
func (builder *Builder) Prepare(configuration ...interface{}) (warnings []string, err error) {
	if len(configuration) == 0 {
		err = fmt.Errorf("No configuration")

		return
	}

	builder.config = &Configuration{}
	err = mapstructure.Decode(configuration[0], builder.config)
	if err != nil {
		return
	}

	err = builder.config.Validate()
	if err != nil {
		return
	}
	builder.client = compute.NewClient(
		builder.config.McpRegion,
		builder.config.McpUser,
		builder.config.McpPassword,
	)
	if os.Getenv("MCP_EXTENDED_LOGGING") != "" {
		builder.client.EnableExtendedLogging()
	}

	return
}

// Run the plugin.
func (builder *Builder) Run(ui packer.Ui, hook packer.Hook, cache packer.Cache) (packer.Artifact, error) {
	config := builder.config

	networkDomain, err := builder.client.GetNetworkDomain(config.NetworkDomainID)
	if err != nil {
		return nil, err
	}
	if networkDomain == nil {
		return nil, fmt.Errorf("Network domain '%s' not found.", config.NetworkDomainID)
	}

	ui.Message(fmt.Sprintf(
		"Deploy server from image '%s' in network domain '%s' (datacenter '%s').",
		builder.config.SourceImage,
		networkDomain.Name,
		networkDomain.DatacenterID,
	))

	return nil, nil
}

// Cancel plugin execution.
func (builder *Builder) Cancel() {
	if builder.client != nil {
		builder.client.Cancel()
	}
}
