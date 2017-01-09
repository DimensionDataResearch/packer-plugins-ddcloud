# CloudControl plugins for Packer

Plugins for [Hashicorp Packer](https://packer.io/) that target Dimension Data CloudControl.

**Note**: this is a work-in-progress; it's not production-ready yet.

See the [plugin documentation](docs/plugins/README.md) for details.

## Installing

Download the appropriate package for the [latest release](https://github.com/DimensionDataResearch/packer-plugins-ddcloud/releases/latest).
Unzip the executable and place it in `~/.packer.d/plugins`.

Needs Packer and OSX or Linux. Both `curl` and VMWare's `ovftool` must be in a directory that's on your `$PATH`.

## Building

Needs Packer <= v0.8.6 (you can run the latest version of Packer, but we can only build against v0.8.6 or lower), Go >= 1.7, and GNU Make.

Sorry, the dependencies are a bit messy at the moment.

1. Fetch correct dependency versions by running `./init.sh` (one-time only).
2. Run `make dev`.

## Sample configurations

See the [plugin documentation](docs/plugins/README.md) for examples.
