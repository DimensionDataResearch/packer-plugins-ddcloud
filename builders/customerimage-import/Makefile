BUILD_ROOT = ../..

include ../../CommonVars.inc
include ../../CommonTargets.inc

PLUGIN_NAME = customerimage-import
PLUGIN_FOLDER = $(BUILDERS_ROOT)/$(PLUGIN_NAME)
EXECUTABLE_NAME = $(EXECUTABLE_PREFIX_BUILDER)-$(PLUGIN_NAME)

default: dev

# Perform a development (current-platform-only) build of the customer image builder plugin and publish it to ~/.packer.d/plugins.
dev: _dev

# Perform an all-platforms build of the customer image builder plugin.
build: _build

# Produce an archive containing the customer image builder plugin for a GitHub release.
dist: _dist
