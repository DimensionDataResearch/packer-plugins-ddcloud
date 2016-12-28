package ddcloud

import (
	"fmt"
	"os"

	"encoding/base64"
	"math/rand"

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
	VLANID          string `mapstructure:"vlan"`
	SourceImage     string `mapstructure:"source_image"`
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

	return
}

// Prepare the plugin to run.
func (builder *Builder) Prepare(configuration ...interface{}) (warnings []string, err error) {
	if len(configuration) == 0 {
		err = fmt.Errorf("No configuration")

		return
	}

	uniquenessKeyBytes := make([]byte, 8)
	_, err = rand.Read(uniquenessKeyBytes)
	if err != nil {
		return
	}
	uniquenessKey := base64.StdEncoding.EncodeToString(uniquenessKeyBytes)

	builder.config = &Configuration{
		uniquenessKey: uniquenessKey,
		serverName:    fmt.Sprintf("packer-build-%s", uniquenessKey),
	}
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
	client := builder.client
	config := builder.config

	networkDomain, err := client.GetNetworkDomain(config.NetworkDomainID)
	if err != nil {
		return nil, err
	}
	if networkDomain == nil {
		return nil, fmt.Errorf("Network domain '%s' not found.", config.NetworkDomainID)
	}

	ui.Message(fmt.Sprintf(
		"Deploy server '%s' from image '%s' in network domain '%s' (datacenter '%s').",
		config.serverName,
		builder.config.SourceImage,
		networkDomain.Name,
		networkDomain.DatacenterID,
	))

	image, err := resolveServerImage(config.SourceImage, networkDomain.DatacenterID, client)
	if err != nil {
		return nil, err
	}
	if image == nil {
		return nil, fmt.Errorf("Cannot find source image named '%s' in datacenter '%s'",
			config.SourceImage,
			networkDomain.DatacenterID,
		)
	}

	serverArtifact, err := deployServer(*config, image, *networkDomain, client)
	if err != nil {
		return nil, err
	}

	ui.Message(fmt.Sprintf(
		"Deployed server '%s' ('%s').",
		serverArtifact.Server.Name,
		serverArtifact.Server.ID,
	))

	// TODO: Create image from server.

	return serverArtifact, nil
}

// Cancel plugin execution.
func (builder *Builder) Cancel() {
	if builder.client != nil {
		builder.client.Cancel()
	}
}
