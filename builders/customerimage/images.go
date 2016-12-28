package main

import (
	"github.com/DimensionDataResearch/go-dd-cloud-compute/compute"
)

func resolveServerImage(imageName string, datacenterID string, client *compute.Client) (compute.Image, error) {
	osImage, err := client.FindOSImage(imageName, datacenterID)
	if err != nil {
		return nil, err
	}
	if osImage != nil {
		return osImage, nil
	}

	customerImage, err := client.FindCustomerImage(imageName, datacenterID)
	if err != nil {
		return nil, err
	}
	if customerImage != nil {
		return customerImage, nil
	}

	return nil, nil
}
