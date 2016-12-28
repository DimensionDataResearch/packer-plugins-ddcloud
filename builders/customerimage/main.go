package main

import (
	"fmt"
	"os"
	"path"

	"github.com/DimensionDataResearch/packer-plugins-ddcloud"
	"github.com/mitchellh/packer/packer/plugin"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == "--version" {
		fmt.Printf("%s %s\n\n", path.Base(os.Args[0]), plugins.ProviderVersion)

		return
	}

	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}

	server.RegisterBuilder(new(Builder))
}
