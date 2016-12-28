package ddcloud

import (
	"time"

	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

func deployServer(config Configuration, image compute.Image, networkDomain compute.NetworkDomain, client *compute.Client) (*ServerArtifact, error) {
	deploymentConfiguration := compute.ServerDeploymentConfiguration{
		Name: config.serverName,
		Network: compute.VirtualMachineNetwork{
			NetworkDomainID: config.NetworkDomainID,
			PrimaryAdapter: compute.VirtualMachineNetworkAdapter{
				VLANID: &config.VLANID,
			},
		},
	}
	image.ApplyTo(&deploymentConfiguration)

	serverID, err := client.DeployServer(deploymentConfiguration)
	if err != nil {
		return nil, err
	}

	resource, err := client.WaitForDeploy(compute.ResourceTypeServer, serverID, 20*time.Minute)
	if err != nil {
		return nil, err
	}

	server := resource.(*compute.Server)

	artifact := &ServerArtifact{
		Server:        *server,
		NetworkDomain: networkDomain,
		deleteServer: func() error {
			deleteError := client.DeleteServer(server.ID)
			if deleteError != nil {
				return deleteError
			}

			return client.WaitForDelete(
				compute.ResourceTypeServer,
				server.ID,
				5*time.Minute,
			)
		},
	}

	return artifact, nil
}
